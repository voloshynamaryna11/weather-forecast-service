package repo

import (
	"context"
	"time"
	"weather-forecast-service/internal/persistence/sqlite"

	"gorm.io/gorm"

	"weather-forecast-service/internal/domain/entity"
	"weather-forecast-service/internal/domain/repository"
)

type WeatherRepo struct{ db *gorm.DB }

func NewWeatherRepo(db *gorm.DB) repository.WeatherRepository {
	return &WeatherRepo{db: db}
}

func (r *WeatherRepo) Save(ctx context.Context, w *entity.Weather) error {
	return r.db.WithContext(ctx).Save(sqlite.FromEntityWeather(w)).Error
}

func (r *WeatherRepo) FindInRange(ctx context.Context, city string, from, to time.Time) (*entity.Weather, error) {
	var w entity.Weather
	err := r.db.WithContext(ctx).
		Where("city = ? AND date BETWEEN ? AND ?", city, from, to).
		Order("date DESC").
		First(&w).Error

	if err != nil {
		return nil, err
	}

	return &w, nil
}

var _ repository.WeatherRepository = (*WeatherRepo)(nil)
