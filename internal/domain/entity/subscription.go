package entity

import (
	"time"
	"weather-forecast-service/internal/domain/entity/enum"
)

type Subscription struct {
	ID        int64
	UserID    int64
	Frequency enum.Frequency
	City      string
	Confirmed bool
	CreatedAt time.Time
}
