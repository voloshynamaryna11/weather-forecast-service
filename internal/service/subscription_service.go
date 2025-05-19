package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"weather-forecast-service/internal/domain/entity"
	"weather-forecast-service/internal/domain/entity/enum"
	"weather-forecast-service/internal/domain/repository"
)

var (
	ErrAlreadyExists        = errors.New("email already subscribed")
	ErrConfirmationNeeded   = errors.New("confirm your first subscription")
	ErrInvalidToken         = errors.New("invalid token")
	ErrTokenNotFound        = errors.New("token not found")
	ErrSubscriptionNotFound = errors.New("subscription not found")
)

type SubscriptionService interface {
	Subscribe(ctx context.Context, email, city string, freq enum.Frequency) error
	Confirm(ctx context.Context, token string) error
	Unsubscribe(ctx context.Context, token string) error
}

type subscriptionService struct {
	subRepo  repository.SubscriptionRepository
	userRepo repository.UserRepository
}

func NewSubscriptionService(subRepo repository.SubscriptionRepository, userRepo repository.UserRepository) SubscriptionService {
	return &subscriptionService{subRepo: subRepo, userRepo: userRepo}
}

func (s *subscriptionService) Subscribe(ctx context.Context, email, city string, freq enum.Frequency) error {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		var entityToSave = &entity.User{
			Email: email,
			Token: uuid.NewString(),
		}
		user, err = s.userRepo.Save(ctx, entityToSave)

		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	existing, err := s.subRepo.FindByUserAndCity(ctx, user.ID, city)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if err == nil && existing != nil {
		if existing.Confirmed {
			return ErrAlreadyExists
		}

		existing.Frequency = freq
		existing.CreatedAt = time.Now()
		return s.subRepo.Update(ctx, existing)
	}

	subs1, err := s.subRepo.FindByUser(ctx, user.ID, false)
	subs2, err := s.subRepo.FindByUser(ctx, user.ID, true)
	if err != nil {
		return err
	}

	if len(subs1) >= 1 && len(subs2) == 0 {
		return ErrConfirmationNeeded
	}

	confirmed := false
	if len(subs2) > 0 {
		confirmed = true
	}

	sub := &entity.Subscription{
		UserID:    user.ID,
		City:      city,
		Frequency: freq,
		Confirmed: confirmed,
		CreatedAt: time.Now(),
	}

	return s.subRepo.Create(ctx, sub)
}

func (s *subscriptionService) Confirm(ctx context.Context, token string) error {
	u, err := s.userRepo.GetByToken(ctx, token)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTokenNotFound
		}
		return err
	}

	if err := s.subRepo.Confirm(ctx, u.ID, u.Email); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSubscriptionNotFound
		}
		return err
	}

	return nil
}

func (s *subscriptionService) Unsubscribe(ctx context.Context, token string) error {
	u, err := s.userRepo.GetByToken(ctx, token)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidToken
		}
		return err
	}

	return s.subRepo.DeleteByUserId(ctx, u.ID)
}
