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

func (r *SubscriptionRepo) Confirm(ctx context.Context, userID int64, email string) error {
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

	// doSomeApiCallToSendEmail(email) + rethrow error if failed
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

func (r *SubscriptionRepo) FindByUserAndCity(ctx context.Context, userID int64, city string) (*entity.Subscription, error) {
	var m sqlite.SubscriptionModel
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND city = ?", userID, city).
		First(&m).Error
	if err != nil {
		return nil, err
	}
	return sqlite.ToEntitySubscription(&m), nil
}

func (r *SubscriptionRepo) Update(ctx context.Context, s *entity.Subscription) error {
	return r.db.WithContext(ctx).
		Save(sqlite.FromEntitySubscription(s)).Error
}

func (r *SubscriptionRepo) FindByUser(ctx context.Context, userID int64, isConfirmed bool) ([]*entity.Subscription, error) {
	var models []sqlite.SubscriptionModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND confirmed = ?", userID, isConfirmed).
		Find(&models).Error; err != nil {
		return nil, err
	}
	subs := make([]*entity.Subscription, len(models))
	for i := range models {
		subs[i] = sqlite.ToEntitySubscription(&models[i])
	}
	return subs, nil
}

var _ repository.SubscriptionRepository = (*SubscriptionRepo)(nil)
