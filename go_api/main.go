package main

import (
	"go_api/internal/database"
	"go_api/internal/handlers"
	"go_api/internal/repository"
	"go_api/internal/services"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
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

	blogpostRepo := repository.NewBlogpostRepository(db)
	blogpostService := services.NewBlogpostService(blogpostRepo)
	blogpostHandler := handlers.NewBlogpostHandler(blogpostService)

	mux.HandleFunc("GET /blogposts", blogpostHandler.ListBlogposts)
	mux.HandleFunc("GET /blogposts/{slug}", blogpostHandler.GetBlogpostBySlug)
	mux.HandleFunc("POST /blogposts", blogpostHandler.CreateBlogpost)
}

func getDBConnectionString() string {
	dbName := os.Getenv("DB_CONNECTION_STRING")

	// Debug logging
	log.Printf("DB_NAME: %s", dbName)
	// Don't log the password

	return dbName
}

func main() {
	db_connection := getDBConnectionString()
	log.Printf("Attempting to connect using: %s", db_connection)

	db, err := database.NewDB(db_connection)
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
