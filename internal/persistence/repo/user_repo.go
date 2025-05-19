package repo

import (
	"context"
	"weather-forecast-service/internal/persistence/sqlite"

	"gorm.io/gorm"

	"weather-forecast-service/internal/domain/entity"
	"weather-forecast-service/internal/domain/repository"
)

type UserRepo struct{ db *gorm.DB }

func NewUserRepo(db *gorm.DB) repository.UserRepository {
	return &UserRepo{db: db}
}

func (r *UserRepo) Get(ctx context.Context, id int64) (*entity.User, error) {
	var m sqlite.UserModel
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return sqlite.ToEntityUser(&m), nil
}

func (r *UserRepo) Save(ctx context.Context, u *entity.User) error {
	return r.db.WithContext(ctx).Save(sqlite.FromEntityUser(u)).Error
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var m sqlite.UserModel
	if err := r.db.WithContext(ctx).
		Where("email = ?", email).First(&m).Error; err != nil {
		return nil, err
	}
	return sqlite.ToEntityUser(&m), nil
}

func (r *UserRepo) GetByToken(ctx context.Context, token string) (*entity.User, error) {
	var m sqlite.UserModel
	if err := r.db.WithContext(ctx).
		Where("token = ?", token).First(&m).Error; err != nil {
		return nil, err
	}
	return sqlite.ToEntityUser(&m), nil
}

var _ repository.UserRepository = (*UserRepo)(nil)
