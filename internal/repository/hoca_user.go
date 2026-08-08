package repository

import (
	"context"

	"gorm.io/gorm"

	"accounting/internal/model"
)

type HocaUserRepository interface {
	GetAll(ctx context.Context) ([]model.HocaUser, error)
	GetByID(ctx context.Context, id int64) (*model.HocaUser, error)
	Create(ctx context.Context, req model.CreateHocaUserRequest) (*model.HocaUser, error)
	Update(ctx context.Context, id int64, req model.UpdateHocaUserRequest) (*model.HocaUser, error)
	Delete(ctx context.Context, id int64) error
}

type hocaUserRepo struct {
	db *gorm.DB
}

func NewHocaUserRepository(db *gorm.DB) HocaUserRepository {
	return &hocaUserRepo{db: db}
}

func (r *hocaUserRepo) GetAll(ctx context.Context) ([]model.HocaUser, error) {
	var users []model.HocaUser
	err := r.db.WithContext(ctx).Order("id").Find(&users).Error
	return users, err
}

func (r *hocaUserRepo) GetByID(ctx context.Context, id int64) (*model.HocaUser, error) {
	var u model.HocaUser
	if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *hocaUserRepo) Create(ctx context.Context, req model.CreateHocaUserRequest) (*model.HocaUser, error) {
	u := model.HocaUser{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
	}
	if err := r.db.WithContext(ctx).Create(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *hocaUserRepo) Update(ctx context.Context, id int64, req model.UpdateHocaUserRequest) (*model.HocaUser, error) {
	result := r.db.WithContext(ctx).Model(&model.HocaUser{}).Where("id = ?", id).Updates(map[string]any{
		"first_name": req.FirstName,
		"last_name":  req.LastName,
		"phone":      req.Phone,
	})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		// Данные могли не измениться — проверяем, существует ли запись.
		if _, err := r.GetByID(ctx, id); err != nil {
			return nil, err
		}
	}
	return r.GetByID(ctx, id)
}

func (r *hocaUserRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.HocaUser{}, id).Error
}
