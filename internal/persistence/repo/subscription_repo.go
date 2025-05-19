package repo

import (
	"context"
	"weather-forecast-service/internal/persistence/sqlite"

	"gorm.io/gorm"

	"weather-forecast-service/internal/domain/entity"
	"weather-forecast-service/internal/domain/repository"
)

type SubscriptionRepo struct{ db *gorm.DB }

func NewSubscriptionRepo(db *gorm.DB) repository.SubscriptionRepository {
	return &SubscriptionRepo{db: db}
}

func (r *SubscriptionRepo) Create(ctx context.Context, s *entity.Subscription) error {
	return r.db.WithContext(ctx).Create(sqlite.FromEntitySubscription(s)).Error
}

func (r *SubscriptionRepo) Confirm(ctx context.Context, userID int64) error {
	res := r.db.WithContext(ctx).
		Model(&sqlite.SubscriptionModel{}).
		Where("user_id = ?", userID).
		Update("confirmed", true)

	if err := res.Error; err != nil {
		return err
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *SubscriptionRepo) DeleteByUserId(ctx context.Context, userID int64) error {
	res := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&sqlite.SubscriptionModel{})

	if err := res.Error; err != nil {
		return err
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

var _ repository.SubscriptionRepository = (*SubscriptionRepo)(nil)
