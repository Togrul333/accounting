package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"accounting/internal/model"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
	UpdateProfile(ctx context.Context, id int64, req model.UpdateProfileRequest) error
	UpdatePassword(ctx context.Context, id int64, hash string) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
	var u model.User
	if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("kullanıcı bulunamadı")
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) UpdateProfile(ctx context.Context, id int64, req model.UpdateProfileRequest) error {
	return r.db.WithContext(ctx).
		Model(&model.User{}).Where("id = ?", id).
		Updates(map[string]any{"name": req.Name, "email": req.Email}).Error
}

func (r *userRepo) UpdatePassword(ctx context.Context, id int64, hash string) error {
	return r.db.WithContext(ctx).
		Model(&model.User{}).Where("id = ?", id).
		Update("password", hash).Error
}
