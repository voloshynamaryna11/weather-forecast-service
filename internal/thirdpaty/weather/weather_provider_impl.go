package weather

import (
	"context"
	"math/rand"
)

type WeatherProvider struct{}

func NewStubProvider() *WeatherProvider {
	return &WeatherProvider{}
}

func (s *WeatherProvider) FetchWeather(ctx context.Context, city string) (WeatherData, error) {
	// Here we are going to implement call to thirdparty to retrieve needed data for our weather app
	// Also some error processing could be present here
	return WeatherData{
		City:        city,
		Temperature: rand.Float64()*30 + 5,
		Humidity:    rand.Float64()*50 + 30,
		Description: "Sunny with clouds",
	}, nil
}
