package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
	"weather-forecast-service/internal/config"
	customHttp "weather-forecast-service/internal/http"
	"weather-forecast-service/internal/persistence/repo"
	"weather-forecast-service/internal/persistence/sqlite"
)

func main() {

	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	cfg := config.MustLoad()
	fmt.Printf("Loaded config: %+v\n", cfg)

	storage, err := sqlite.New(cfg.StoragePath, log)

	if err != nil {
		log.Error("failed to init storage", "err", err)
	}

	_ = storage

	wRepo := repo.NewWeatherRepo(storage.DB())

	/* ---------- HTTP router ------------ */
	router := customHttp.NewRouter(wRepo)

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
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	log.Info("shutting down ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Info("bye")
}
