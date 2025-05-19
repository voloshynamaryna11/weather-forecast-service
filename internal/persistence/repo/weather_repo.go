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

func (r *WeatherRepo) Find(ctx context.Context, city string, date time.Time) (*entity.Weather, error) {
	var m sqlite.WeatherModel
	err := r.db.WithContext(ctx).
		Where("city = ? AND date = ?", city, date).
		First(&m).Error
	if err != nil {
		return nil, err
	}
	return sqlite.ToEntityWeather(&m), nil
}

var _ repository.WeatherRepository = (*WeatherRepo)(nil)
