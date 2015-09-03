package main

import (
	"finance/yahoo"
	"app"
	"main/controller"
)

func main() {
	var financeService *yahoo.YahooFinanceService
	financeService, err := yahoo.New(yahoo.Config{Url:"https://query.yahooapis.com/v1/public/yql"})
	if err != nil {
		panic(err)
		return
	}

	mainController := controller.New(financeService)

	app := app.New()
	app.Get("/", mainController)
	app.Listen(8080)
}