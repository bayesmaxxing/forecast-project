package routes

import (
	"backend/internal/auth"
	"backend/internal/handlers"
	"backend/internal/repository"
	"backend/internal/services"
	"net/http"
)

type Handlers struct {
	Forecast      *handlers.ForecastHandler
	ForecastPoint *handlers.ForecastPointHandler
	User          *handlers.UserHandler
	Score         *handlers.ScoreHandler
}

type Services struct {
	Forecast      *services.ForecastService
	ForecastPoint *services.ForecastPointService
	User          *services.UserService
	Score         *services.ScoreService
}

type Repositories struct {
	Forecast      repository.ForecastRepository
	ForecastPoint repository.ForecastPointRepository
	User          repository.UserRepository
	Score         repository.ScoreRepository
}

func Setup(mux *http.ServeMux, handlers *Handlers) {
	//public routes
	setupPublicRoutes(mux, handlers)

	// protected routes
	protected := http.NewServeMux()
	setupProtectedRoutes(protected, handlers)

	apiHandler := http.StripPrefix("/api", protected)
	mux.Handle("/api/", auth.AuthMiddleware(apiHandler))
}

func setupPublicRoutes(mux *http.ServeMux, handlers *Handlers) {
	// forecasts
	mux.HandleFunc("POST /forecasts", handlers.Forecast.ListForecasts)
	mux.HandleFunc("GET /forecasts/{id}", handlers.Forecast.GetForecast)
	mux.HandleFunc("GET /forecasts/stale-and-new/{user_id}", handlers.Forecast.GetStaleAndNewForecasts)

	// forecast points
	mux.HandleFunc("GET /forecast-points/{id}", handlers.ForecastPoint.ListForecastPointsbyID)
	mux.HandleFunc("POST /forecast-points/user", handlers.ForecastPoint.ListForecastPointsbyIDAndUser)
	mux.HandleFunc("GET /forecast-points", handlers.ForecastPoint.ListAllForecastPoints)
	mux.HandleFunc("GET /forecast-points/latest", handlers.ForecastPoint.ListLatestForecastPoints)
	mux.HandleFunc("GET /forecast-points/latest_by_user", handlers.ForecastPoint.ListLatestForecastPointsByUser)
	mux.HandleFunc("GET /forecast-points/ordered/{id}", handlers.ForecastPoint.ListOrderedForecastPoints)

	// scores (single-score)
	mux.HandleFunc("POST /scores", handlers.Score.GetScores)
	mux.HandleFunc("GET /scores/all", handlers.Score.GetAllScores)
	mux.HandleFunc("GET /scores/average", handlers.Score.GetAverageScores)
	mux.HandleFunc("GET /scores/average/{id}", handlers.Score.GetAverageScoreByForecastID)

	// scores (aggregate)
	mux.HandleFunc("POST /scores/aggregate", handlers.Score.GetAggregateScores)

	// users
	mux.HandleFunc("GET /users", handlers.User.ListUsers)
	mux.HandleFunc("POST /users", handlers.User.CreateUser)
	mux.HandleFunc("POST /users/login", handlers.User.Login)
	// mux.HandleFunc("POST /users/reset-password", handlers.User.AdminResetPassword)
}

func setupProtectedRoutes(mux *http.ServeMux, handlers *Handlers) {
	// forecasts
	mux.HandleFunc("POST /forecasts/create", handlers.Forecast.CreateForecast)
	mux.HandleFunc("DELETE /forecasts", handlers.Forecast.DeleteForecast)
	mux.HandleFunc("PUT /resolve", handlers.Forecast.ResolveForecast)

	// forecast points
	mux.HandleFunc("POST /forecast-points", handlers.ForecastPoint.CreateForecastPoint)

	// scores (single-score)
	mux.HandleFunc("POST /scores", handlers.Score.CreateScore)
	mux.HandleFunc("DELETE /scores", handlers.Score.DeleteScore)

	// users
	mux.HandleFunc("DELETE /users", handlers.User.DeleteUser)
	mux.HandleFunc("PUT /users/password", handlers.User.ChangePassword)

}
