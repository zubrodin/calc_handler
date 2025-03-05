package api

import (
	"encoding/json"
	"net/http"

	"github.com/zubrodin/calc_handler/internal/models"
	"github.com/zubrodin/calc_handler/internal/orchestrator"
)

func SetupRoutes() {
	http.HandleFunc("/api/v1/calculate", orchestrator.AddExpression)
	http.HandleFunc("/api/v1/expressions", orchestrator.GetExpressions)
	http.HandleFunc("/api/v1/expressions/", orchestrator.GetExpressionByID)
	http.HandleFunc("/internal/task", orchestrator.GetTask)
	http.HandleFunc("/internal/task/result", orchestrator.HandleResult)
}

func HandleResults(w http.ResponseWriter, r *http.Request) {
	var result models.Result
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, "Invalid data", http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusOK)
}
