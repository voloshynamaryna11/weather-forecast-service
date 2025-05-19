package thirdpaty

import (
	"context"
	"log"
	"time"
	"weather-forecast-service/internal/domain/repository"
	"weather-forecast-service/internal/thirdpaty/weather"
	"weather-forecast-service/internal/thirdpaty/weather/mapper"
)

type Fetcher struct {
	Provider weather.Provider
	Cities   []string
	Repo     repository.WeatherRepository
}

func NewFetcher(provider weather.Provider, cities []string, repo repository.WeatherRepository) *Fetcher {
	return &Fetcher{
		Provider: provider,
		Cities:   cities,
		Repo:     repo,
	}
}

func (f *Fetcher) Start(ctx context.Context) {
	f.fetchAll(ctx)

	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				f.fetchAll(ctx)
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func (f *Fetcher) fetchAll(ctx context.Context) {
	for _, city := range f.Cities {
		data, err := f.Provider.FetchWeather(ctx, city)
		if err != nil {
			log.Printf("Failed to fetch weather for %s: %v", city, err)
			continue
		}

		weatherEntity := mapper.ToWeather(data)

		log.Printf("Fetched and mapped weather for %s: %+v", city, weatherEntity)

		if err := f.Repo.Save(ctx, weatherEntity); err != nil {
			log.Printf("Failed to save weather for %s: %v", city, err)
			continue
		}

		log.Printf("Weather saved for %s", city)
	}
}
