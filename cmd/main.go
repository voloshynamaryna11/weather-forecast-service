package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
	"weather-forecast-service/internal/config"
	customHttp "weather-forecast-service/internal/http"
	"weather-forecast-service/internal/http/handlers"
	"weather-forecast-service/internal/persistence/repo"
	"weather-forecast-service/internal/persistence/sqlite"
	"weather-forecast-service/internal/service"
	"weather-forecast-service/internal/thirdpaty"
	"weather-forecast-service/internal/thirdpaty/weather"
)

func main() {

	// Logger
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	// Load config
	cnf := config.MustLoad()

	// SQLite storage
	storage, err := sqlite.New(cnf.StoragePath, log)
	if err != nil {
		log.Error("failed to init storage", "err", err)
		os.Exit(1)
	}

	// Repositories
	wRepo := repo.NewWeatherRepo(storage.DB())
	sRepo := repo.NewSubscriptionRepo(storage.DB())
	uRepo := repo.NewUserRepo(storage.DB())

	// Weather fetcher
	ctx := context.Background()
	provider := weather.NewStubProvider()
	cities := []string{"Kyiv", "Lviv", "Odesa"}

	fetcher := thirdpaty.NewFetcher(provider, cities, wRepo)
	fetcher.Start(ctx)

	// Services
	wSvc := service.NewWeatherService(wRepo)
	subSvc := service.NewSubscriptionService(sRepo, uRepo)

	// Handlers
	weatherH := handlers.NewWeatherHandler(wSvc)
	subH := handlers.NewSubscribeHandler(subSvc)
	confirmH := handlers.NewConfirmHandler(subSvc)
	unsubH := handlers.NewUnsubscribeHandler(subSvc)

	// Router
	router := customHttp.NewRouter(
		weatherH,
		subH,
		confirmH,
		unsubH,
	)

	// Use config values for HTTP server
	srv := &http.Server{
		Addr:        cnf.HTTPServer.Address,
		Handler:     router,
		ReadTimeout: cnf.HTTPServer.Timeout,
	}

	// Start server
	go func() {
		log.Info("HTTP server listening", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("shutdown error", "err", err)
	}
	log.Info("server stopped")
}
