package sqlite

import (
	"weather-forecast-service/internal/domain/entity"
)

func ToEntityWeather(m *WeatherModel) *entity.Weather {
	return &entity.Weather{
		ID:          m.ID,
		City:        m.City,
		Description: m.Description,
		Temperature: m.Temperature,
		Humidity:    m.Humidity,
		Date:        m.Date,
	}
}

func FromEntityWeather(e *entity.Weather) *WeatherModel {
	return &WeatherModel{
		ID:          e.ID,
		City:        e.City,
		Description: e.Description,
		Temperature: e.Temperature,
		Humidity:    e.Humidity,
		Date:        e.Date,
	}
}

func ToEntitySubscription(m *SubscriptionModel) *entity.Subscription {
	return &entity.Subscription{
		ID:        m.ID,
		UserID:    m.UserID,
		Frequency: m.Frequency,
		City:      m.City,
		Confirmed: m.Confirmed,
		CreatedAt: m.CreatedAt,
	}
}

func FromEntitySubscription(e *entity.Subscription) *SubscriptionModel {
	return &SubscriptionModel{
		ID:        e.ID,
		UserID:    e.UserID,
		Frequency: e.Frequency,
		City:      e.City,
		Confirmed: e.Confirmed,
		CreatedAt: e.CreatedAt,
	}
}

func ToEntityUser(m *UserModel) *entity.User {
	return &entity.User{
		ID:    m.ID,
		Email: m.Email,
		Token: m.Token,
	}
}

func FromEntityUser(e *entity.User) *UserModel {
	return &UserModel{
		ID:    e.ID,
		Email: e.Email,
		Token: e.Token,
	}
}
