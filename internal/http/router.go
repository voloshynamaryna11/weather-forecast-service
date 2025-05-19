package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"weather-forecast-service/internal/http/handlers"
)

func NewRouter(
	weatherHandler *handlers.WeatherHandler,
	subscribeHandler *handlers.SubscriptionHandler,
	confirmHandler *handlers.ConfirmHandler,
	unsubscribeHandler *handlers.UnsubscribeHandler,
) http.Handler {
	r := chi.NewRouter()

	// common middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Get("/weather", weatherHandler.Get)

		r.Post("/subscribe", subscribeHandler.Post)

		r.Get("/confirm/{token}", confirmHandler.Get)

		r.Get("/unsubscribe/{token}", unsubscribeHandler.Get)
	})

	return r
}
