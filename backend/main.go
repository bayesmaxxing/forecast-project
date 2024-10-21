package main

import (
	"backend/internal/cache"
	"backend/internal/database"
	"backend/internal/handlers"
	"backend/internal/repository"
	"backend/internal/services"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func setupRoutes(mux *http.ServeMux, db *database.DB) {
	cacheInstance := cache.NewCache()

	forecastRepo := repository.NewForecastRepository(db)
	forecastPointRepo := repository.NewForecastPointRepository(db)
	forecastService := services.NewForecastService(forecastRepo, forecastPointRepo)
	forecastHandler := handlers.NewForecastHandler(forecastService, cacheInstance)

	mux.HandleFunc("GET /forecasts", forecastHandler.ListForecasts)
	mux.HandleFunc("GET /forecasts/{id}", forecastHandler.GetForecast)
	mux.HandleFunc("POST /forecasts", forecastHandler.CreateForecast)
	mux.HandleFunc("DELETE /forecasts/{id}", forecastHandler.DeleteForecast)
	mux.HandleFunc("PUT /resolve/{id}", forecastHandler.ResolveForecast)
	mux.HandleFunc("GET /scores", forecastHandler.GetAggregatedScores)

	forecastPointService := services.NewForecastPointService(forecastPointRepo)
	forecastPointHandler := handlers.NewForecastPointHandler(forecastPointService, cacheInstance)

	mux.HandleFunc("GET /forecast-points/{id}", forecastPointHandler.ListForecastPointsbyID)
	mux.HandleFunc("GET /forecast-points", forecastPointHandler.ListAllForecastPoints)
	mux.HandleFunc("POST /forecast-points", forecastPointHandler.CreateForecastPoint)
	mux.HandleFunc("GET /forecast-points/latest", forecastPointHandler.ListLatestForecastPoints)

	blogpostRepo := repository.NewBlogpostRepository(db)
	blogpostService := services.NewBlogpostService(blogpostRepo)
	blogpostHandler := handlers.NewBlogpostHandler(blogpostService)

	mux.HandleFunc("GET /blogposts", blogpostHandler.ListBlogposts)
	mux.HandleFunc("GET /blogposts/{slug}", blogpostHandler.GetBlogpostBySlug)
	mux.HandleFunc("POST /blogposts", blogpostHandler.CreateBlogpost)
}

func getDBConnectionString() string {
	dbName := os.Getenv("DB_CONNECTION_STRING")
	return dbName
}

// CORS middleware
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "https://www.samuelsforecasts.com")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	db_connection := getDBConnectionString()
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
	handler := CORSMiddleware(mux)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}
