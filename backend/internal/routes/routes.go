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
	Blogpost      *handlers.BlogpostHandler
	User          *handlers.UserHandler
	Score         *handlers.ScoreHandler
}

type Services struct {
	Forecast      *services.ForecastService
	ForecastPoint *services.ForecastPointService
	Blogpost      *services.BlogpostService
	User          *services.UserService
	Score         *services.ScoreService
}

type Repositories struct {
	Forecast      repository.ForecastRepository
	ForecastPoint repository.ForecastPointRepository
	Blogpost      repository.BlogpostRepository
	User          repository.UserRepository
	Score         repository.ScoreRepository
}

func Setup(mux *http.ServeMux, handlers *Handlers) {
	//public routes
	setupPublicRoutes(mux, handlers)

	// protected routes
	protected := http.NewServeMux()
	setupProtectedRoutes(protected, handlers)
	mux.Handle("/api/", auth.AuthMiddleware(protected))
}

func setupPublicRoutes(mux *http.ServeMux, handlers *Handlers) {
	// forecasts
	mux.HandleFunc("POST /forecasts", handlers.Forecast.ListForecasts)
	mux.HandleFunc("GET /forecasts/{id}", handlers.Forecast.GetForecast)

	// forecast points
	mux.HandleFunc("GET /forecast-points/{id}", handlers.ForecastPoint.ListForecastPointsbyID)
	mux.HandleFunc("GET /forecast-points", handlers.ForecastPoint.ListAllForecastPoints)
	mux.HandleFunc("GET /forecast-points/latest", handlers.ForecastPoint.ListLatestForecastPoints)
	mux.HandleFunc("GET /forecast-points/latest/{user_id}", handlers.ForecastPoint.ListLatestForecastPointsByUser)
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

	// blogposts
	mux.HandleFunc("GET /blogposts", handlers.Blogpost.ListBlogposts)
	mux.HandleFunc("GET /blogposts/{slug}", handlers.Blogpost.GetBlogpostBySlug)
	mux.HandleFunc("POST /blogposts", handlers.Blogpost.CreateBlogpost)
}

func setupProtectedRoutes(mux *http.ServeMux, handlers *Handlers) {
	// forecasts
	mux.HandleFunc("POST /forecasts", handlers.Forecast.CreateForecast)
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
