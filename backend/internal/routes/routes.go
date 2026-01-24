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
	Calibration   *handlers.CalibrationHandler
}

type Services struct {
	Forecast      *services.ForecastService
	ForecastPoint *services.ForecastPointService
	User          *services.UserService
	Score         *services.ScoreService
	Calibration   *services.CalibrationService
}

type Repositories struct {
	Forecast      repository.ForecastRepository
	ForecastPoint repository.ForecastPointRepository
	User          repository.UserRepository
	Score         repository.ScoreRepository
	Calibration   repository.CalibrationRepository
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
	mux.HandleFunc("GET /forecasts", handlers.Forecast.ListForecasts)
	mux.HandleFunc("GET /forecasts/{id}", handlers.Forecast.GetForecast)
	mux.HandleFunc("GET /forecasts/llm/{user_id}", handlers.Forecast.GetStaleAndNewForecasts)

	// forecast points
	mux.HandleFunc("GET /forecast-points", handlers.ForecastPoint.ListForecastPoints)

	// scores (single-score)
	mux.HandleFunc("GET /scores", handlers.Score.GetScores)
	mux.HandleFunc("GET /scores/average", handlers.Score.GetAverageScores)

	// scores (aggregate)
	mux.HandleFunc("GET /scores/aggregate", handlers.Score.GetAggregateScores)
	mux.HandleFunc("GET /scores/aggregate/users", handlers.Score.GetAggregateScoresGroupedByUsers)

	// users
	mux.HandleFunc("GET /users", handlers.User.ListUsers)
	mux.HandleFunc("POST /users", handlers.User.CreateUser)
	mux.HandleFunc("POST /users/login", handlers.User.Login)
	// mux.HandleFunc("POST /users/reset-password", handlers.User.AdminResetPassword)

	// calibration
	mux.HandleFunc("GET /calibration", handlers.Calibration.GetCalibration)
	mux.HandleFunc("GET /calibration/users", handlers.Calibration.GetCalibrationByUsers)
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
