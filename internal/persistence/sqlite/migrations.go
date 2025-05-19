package sqlite

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func initMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "0001_initial",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&WeatherModel{}, &UserModel{}, &SubscriptionModel{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&SubscriptionModel{}, &UserModel{}, &WeatherModel{})
			},
		},
	}
}
