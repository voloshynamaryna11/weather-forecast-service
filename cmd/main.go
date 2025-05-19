package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
	customHttp "weather-forecast-service/internal/http"
	"weather-forecast-service/internal/http/handlers"
	"weather-forecast-service/internal/persistence/repo"
	"weather-forecast-service/internal/persistence/sqlite"
	"weather-forecast-service/internal/service"
	"weather-forecast-service/internal/thirdpaty"
	"weather-forecast-service/internal/thirdpaty/weather"
)

func main() {

	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	ctx := context.Background()
	provider := weather.NewStubProvider()
	cities := []string{"Kyiv", "Lviv", "Odesa"}

	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "file:weather.db?_fk=1"
	}
	log.Info("using SQLite DSN", "dsn", storagePath)

	storage, err := sqlite.New(storagePath, log)
	if err != nil {
		log.Error("failed to init storage", "err", err)
		os.Exit(1)
	}

	wRepo := repo.NewWeatherRepo(storage.DB())
	sRepo := repo.NewSubscriptionRepo(storage.DB())
	uRepo := repo.NewUserRepo(storage.DB())

	fetcher := thirdpaty.NewFetcher(provider, cities, wRepo)
	fetcher.Start(ctx)

	wSvc := service.NewWeatherService(wRepo)
	subSvc := service.NewSubscriptionService(sRepo, uRepo)

	weatherH := handlers.NewWeatherHandler(wSvc)
	subH := handlers.NewSubscribeHandler(subSvc)
	confirmH := handlers.NewConfirmHandler(subSvc)
	unsubH := handlers.NewUnsubscribeHandler(subSvc)

	router := customHttp.NewRouter(
		weatherH,
		subH,
		confirmH,
		unsubH,
	)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Info("HTTP server listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

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
