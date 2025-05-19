package service

import (
	"context"
	"errors"
	"time"

	"weather-forecast-service/internal/domain/repository"
)

type WeatherResponse struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Description string  `json:"description"`
}

var ErrNotFound = errors.New("weather not found")

type WeatherService interface {
	Get(ctx context.Context, city string) (WeatherResponse, error)
}

type weatherService struct {
	repo repository.WeatherRepository
}

func NewWeatherService(repo repository.WeatherRepository) WeatherService {
	return &weatherService{repo: repo}
}

func (s *weatherService) Get(ctx context.Context, city string) (WeatherResponse, error) {
	day := time.Now().Truncate(24 * time.Hour)

	w, err := s.repo.Find(ctx, city, day)
	if err != nil {
		return WeatherResponse{}, ErrNotFound
	}
	return WeatherResponse{
		Temperature: w.Temperature,
		Humidity:    w.Humidity,
		Description: w.Description,
	}, nil
}
