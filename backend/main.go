package main

import (
	"backend/internal/auth"
	"backend/internal/cache"
	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handlers"
	"backend/internal/middleware"
	"backend/internal/repository"
	"backend/internal/routes"
	"backend/internal/services"
	"context"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// CORSMiddleware creates a CORS middleware with the specified allowed origin.
func CORSMiddleware(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
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
}

func main() {
	ctx := context.Background()

	// Load configuration from Secret Manager and environment
	cfg, err := config.Load(ctx)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize auth with JWT secret
	if err := auth.Init(cfg.JWTSecret); err != nil {
		log.Fatalf("Error initializing auth: %v", err)
	}

	db, err := database.NewDB(cfg.DBConnString)
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
	}

	handlers := &routes.Handlers{
		Forecast:      handlers.NewForecastHandler(services.Forecast),
		ForecastPoint: handlers.NewForecastPointHandler(services.ForecastPoint),
		User:          handlers.NewUserHandler(services.User),
		Score:         handlers.NewScoreHandler(services.Score),
	}

	mux := http.NewServeMux()
	routes.Setup(mux, handlers)
	handler := CORSMiddleware(cfg.AllowedOrigin)(mux)
	handler = middleware.RequestLogger(handler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}
