package repository

import (
	"context"
	"time"
	"weather-forecast-service/internal/domain/entity"
)

type WeatherRepository interface {
	Save(ctx context.Context, w *entity.Weather) error
	Find(ctx context.Context, city string, date time.Time) (*entity.Weather, error)
}

type SubscriptionRepository interface {
	Create(ctx context.Context, s *entity.Subscription) error
	Confirm(ctx context.Context, userID int64) error
	DeleteByUserId(ctx context.Context, id int64) error
}

type UserRepository interface {
	Get(ctx context.Context, id int64) (*entity.User, error)
	Save(ctx context.Context, u *entity.User) error
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByToken(ctx context.Context, token string) (*entity.User, error)
}
