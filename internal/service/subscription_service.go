package service

import (
	"context"
	"errors"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"weather-forecast-service/internal/domain/entity"
	"weather-forecast-service/internal/domain/entity/enum"
	"weather-forecast-service/internal/domain/repository"
)

var (
	ErrInvalidFrequency     = errors.New("invalid frequency")
	ErrAlreadyExists        = errors.New("subscription already exists")
	ErrInvalidToken         = errors.New("invalid token")
	ErrInvalidEmail         = errors.New("invalid email")
	ErrInvalidCity          = errors.New("city is required")
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

	//Params validation
	if email == "" {
		return ErrInvalidEmail
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return ErrInvalidEmail
	}

	if city == "" {
		return ErrInvalidCity
	}

	if !freq.IsValid() {
		return ErrInvalidFrequency
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user = &entity.User{
			Email: email,
			Token: uuid.NewString(),
		}
		if err := s.userRepo.Save(ctx, user); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	sub := &entity.Subscription{
		UserID:    user.ID,
		City:      city,
		Frequency: freq,
		Confirmed: false,
		CreatedAt: time.Now(),
	}
	if err := s.subRepo.Create(ctx, sub); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrAlreadyExists
		}
		return err
	}

	return nil
}

func (s *subscriptionService) Confirm(ctx context.Context, token string) error {
	u, err := s.userRepo.GetByToken(ctx, token)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidToken
		}
		return err
	}

	if err := s.subRepo.Confirm(ctx, u.ID); err != nil {
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
