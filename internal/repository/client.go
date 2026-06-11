package repository

import (
	"context"

	"gorm.io/gorm"

	"accounting/internal/model"
)

type ClientRepository interface {
	GetAll(ctx context.Context) ([]model.Client, error)
	GetByID(ctx context.Context, id int64) (*model.Client, error)
	Create(ctx context.Context, req model.CreateClientRequest) (*model.Client, error)
	Update(ctx context.Context, id int64, req model.UpdateClientRequest) (*model.Client, error)
	Delete(ctx context.Context, id int64) error
}

type clientRepo struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) ClientRepository {
	return &clientRepo{db: db}
}

func (r *clientRepo) GetAll(ctx context.Context) ([]model.Client, error) {
	var clients []model.Client
	err := r.db.WithContext(ctx).Order("last_name, first_name").Find(&clients).Error
	return clients, err
}

func (r *clientRepo) GetByID(ctx context.Context, id int64) (*model.Client, error) {
	var c model.Client
	err := r.db.WithContext(ctx).First(&c, id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *clientRepo) Create(ctx context.Context, req model.CreateClientRequest) (*model.Client, error) {
	c := model.Client{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		BirthYear: req.BirthYear,
	}
	if err := r.db.WithContext(ctx).Create(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *clientRepo) Update(ctx context.Context, id int64, req model.UpdateClientRequest) (*model.Client, error) {
	result := r.db.WithContext(ctx).Model(&model.Client{}).Where("id = ?", id).Updates(map[string]any{
		"first_name": req.FirstName,
		"last_name":  req.LastName,
		"email":      req.Email,
		"phone":      req.Phone,
		"birth_year": req.BirthYear,
	})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return r.GetByID(ctx, id)
}

func (r *clientRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Client{}, id).Error
}
