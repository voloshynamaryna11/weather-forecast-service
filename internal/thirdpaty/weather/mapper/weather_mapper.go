package mapper

import (
	"time"
	"weather-forecast-service/internal/thirdpaty/weather"

	"weather-forecast-service/internal/domain/entity"
)

func ToWeather(data weather.WeatherData) *entity.Weather {
	return &entity.Weather{
		City:        data.City,
		Description: data.Description,
		Temperature: data.Temperature,
		Humidity:    data.Humidity,
		Date:        time.Now(),
	}
}
