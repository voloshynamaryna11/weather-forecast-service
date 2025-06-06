package sqlite

import (
	"time"
	"weather-forecast-service/internal/domain/entity"
	"weather-forecast-service/internal/domain/entity/enum"
)

type WeatherModel struct {
	ID          int64  `gorm:"primaryKey"`
	City        string `gorm:"index:idx_city_date,unique"`
	Description string
	Temperature float64
	Humidity    float64
	Date        time.Time `gorm:"index:idx_city_date,unique"`
}

type SubscriptionModel struct {
	ID        int64          `gorm:"primaryKey"`
	UserID    int64          `gorm:"not null;uniqueIndex:idx_user_city"`
	User      entity.User    `gorm:"constraint:OnDelete:CASCADE"`
	Frequency enum.Frequency `gorm:"type:text"`
	City      string         `gorm:"not null;uniqueIndex:idx_user_city"`
	Confirmed bool
	CreatedAt time.Time
}

type UserModel struct {
	ID    int64  `gorm:"primaryKey"`
	Email string `gorm:"uniqueIndex"`
	Token string `gorm:"uniqueIndex"`
}
