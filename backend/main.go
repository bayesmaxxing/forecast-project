package main

import (
	"backend/internal/cache"
	"backend/internal/database"
	"backend/internal/handlers"
	"backend/internal/middleware"
	"backend/internal/repository"
	"backend/internal/routes"
	"backend/internal/services"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func getDBConnectionString() string {
	dbName := os.Getenv("DB_CONNECTION_STRING")
	return dbName
}

// CORS middleware
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*") // https://www.samuelsforecasts.com
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

	repositories := &routes.Repositories{
		Forecast:      repository.NewForecastRepository(db),
		ForecastPoint: repository.NewForecastPointRepository(db),
		User:          repository.NewUserRepository(db),
		Score:         repository.NewScoreRepository(db),
	}

	cache := cache.NewCache()

	services := &routes.Services{
		Forecast:      services.NewForecastService(repositories.Forecast, repositories.ForecastPoint, repositories.Score, cache),
		ForecastPoint: services.NewForecastPointService(repositories.ForecastPoint, repositories.Forecast, cache),
		User:          services.NewUserService(repositories.User, cache),
		Score:         services.NewScoreService(repositories.Score, cache),
		News:          services.NewNewsService(),
	}

	handlers := &routes.Handlers{
		Forecast:      handlers.NewForecastHandler(services.Forecast),
		ForecastPoint: handlers.NewForecastPointHandler(services.ForecastPoint),
		User:          handlers.NewUserHandler(services.User),
		Score:         handlers.NewScoreHandler(services.Score),
		News:          handlers.NewNewsHandler(services.News),
	}

	mux := http.NewServeMux()
	routes.Setup(mux, handlers)
	handler := CORSMiddleware(mux)
	handler = middleware.RequestLogger(handler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}
