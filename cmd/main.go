package main

import (
	"github.com/zubrodin/calc_handler/application"
)

func main() {
	app := application.New()
	app.RunServer()
}
