package http

import (
	"net/http"
	"weather-forecast-service/internal/domain/repository"
	"weather-forecast-service/internal/http/handlers"
	"weather-forecast-service/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(wRepo repository.WeatherRepository) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.Logger, middleware.Recoverer)

	weatherSvc := service.NewWeatherService(wRepo)
	weatherH := handlers.New(weatherSvc)

	r.Route("/api", func(api chi.Router) {
		api.Get("/weather", weatherH.Get)
	})
	return r
}
