package weather

import (
	"context"
)

// Structure to hold retrieved data from thirdparty provider
type WeatherData struct {
	City        string
	Temperature float64
	Humidity    float64
	Description string
}

type Provider interface {
	FetchWeather(ctx context.Context, city string) (WeatherData, error)
}
