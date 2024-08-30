package main

import (
	"go_api/internal/database"
	"go_api/internal/handlers"
	"go_api/internal/repository"
	"go_api/internal/services"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5"
)

func setupRoutes(mux *http.ServeMux, db *database.DB) {
	forecastRepo := repository.NewForecastRepository(db)
	forecastPointRepo := repository.NewForecastPointRepository(db)
	forecastService := services.NewForecastService(forecastRepo, forecastPointRepo)
	forecastHandler := handlers.NewForecastHandler(forecastService)

	mux.HandleFunc("GET /forecasts", forecastHandler.ListForecasts)
	mux.HandleFunc("GET /forecasts/{id}", forecastHandler.GetForecast)
	mux.HandleFunc("POST /forecasts", forecastHandler.CreateForecast)
	mux.HandleFunc("DELETE /forecasts/{id}", forecastHandler.DeleteForecast)
	mux.HandleFunc("GET /scores", forecastHandler.GetAggregatedScores)

	forecastPointService := services.NewForecastPointService(forecastPointRepo)
	forecastPointHandler := handlers.NewForecastPointHandler(forecastPointService)

	mux.HandleFunc("GET /forecast-points/{id}", forecastPointHandler.ListForecastPointsbyID)
	mux.HandleFunc("GET /forecast-points", forecastPointHandler.ListAllForecastPoints)
	mux.HandleFunc("POST /forecast-points", forecastPointHandler.CreateForecastPoint)
}

func main() {
	db, err := database.NewDB("connection_string")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	mux := http.NewServeMux()
	setupRoutes(mux, db)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}
