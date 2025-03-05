package main

import (
	"log"
	"net/http"

	"github.com/zubrodin/calc_handler/api"
)

func main() {

	api.SetupRoutes()

	port := "8080"
	log.Printf("Сервер запущен на порту %s", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %s", err)
	}
}
