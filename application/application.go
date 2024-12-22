package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/zubrodin/calc_handler/pkg/calc"
)

type Config struct {
	Add string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Add = os.Getenv("PORT")
	if config.Add == "" {
		config.Add = "8080"
	}
	return config
}

type Application struct {
	Config *Config
}

func New() *Application {
	return &Application{
		Config: ConfigFromEnv(),
	}
}

type Request struct {
	Expression string `json:"expression"`
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	request := new(Request)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	result, err := calc.Calc(request.Expression)
	if err != nil {
		if errors.Is(err, calc.ErrEmptyExpression) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "error: %s", err.Error())
			return

		} else if errors.Is(err, calc.ErrInvalidExpression) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "error: %s", err.Error())
			return

		} else if errors.Is(err, calc.ErrDivisionByZero) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "error: %s", err.Error())
			return

		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error: %s", "unknown error")
			return
		}

	} else {
		fmt.Fprintf(w, "result: %f", result)
	}
}

func (a *Application) RunServer() error {
	http.HandleFunc("/api/v1/calculate", CalcHandler)
	return http.ListenAndServe(":"+a.Config.Add, nil)
}
