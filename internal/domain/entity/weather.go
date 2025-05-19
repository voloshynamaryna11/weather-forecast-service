package entity

import "time"

type Weather struct {
	ID          int64
	City        string
	Description string
	Temperature float64
	Humidity    float64
	Date        time.Time
}

func (Weather) TableName() string {
	return "weather_models"
}
