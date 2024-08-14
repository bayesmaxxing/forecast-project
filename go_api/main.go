package main

import (
	"net/http"
)

type App struct {
	router *http.ServeMux
	// more here
}

func NewApp() *App {
	app := &App{
		router: http.NewServeMux(),
	}
	app.routes()
	return app
}

func (a *App) routes() {
	// Forecast methods, all forecasts
	a.router.HandleFunc("GET /forecasts", a.getForecasts)

	// Specific forecast methods
	a.router.HandleFunc("POST /forecasts", a.createForecast)
	a.router.HandleFunc("GET /forecasts/{id}", a.getForecast)
	a.router.HandleFunc("PUT /forecasts", a.resolveForecast)

	//Aggregate methods
	a.router.HandleFunc("GET /forecasts/scores")

	// Forecast point methods
	a.router.HandleFunc("GET /forecast_points", a.getForecastPoints)
	a.router.HandleFunc("POST /forecast_points", a.createForecastPoint)
	a.router.HandleFunc("GET /forecast_points", a.handleForecastPoints)

}

func (a *App) getForecasts(w http.ResponseWriter, r *http.Request) {
	// implement get forecasts logic
}
