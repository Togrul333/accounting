package repository

import (
	"context"

	"gorm.io/gorm"

	"accounting/internal/model"
)

type IncomeCategoryRepository interface {
	GetAll(ctx context.Context) ([]model.IncomeCategory, error)
	GetByID(ctx context.Context, id int64) (*model.IncomeCategory, error)
	Create(ctx context.Context, req model.CreateIncomeCategoryRequest) (*model.IncomeCategory, error)
	Update(ctx context.Context, id int64, req model.UpdateIncomeCategoryRequest) (*model.IncomeCategory, error)
	Delete(ctx context.Context, id int64) error
}

type incomeCategoryRepo struct {
	db *gorm.DB
}

func NewIncomeCategoryRepository(db *gorm.DB) IncomeCategoryRepository {
	return &incomeCategoryRepo{db: db}
}

func (r *incomeCategoryRepo) GetAll(ctx context.Context) ([]model.IncomeCategory, error) {
	var cats []model.IncomeCategory
	err := r.db.WithContext(ctx).Order("id").Find(&cats).Error
	return cats, err
}

func (r *incomeCategoryRepo) GetByID(ctx context.Context, id int64) (*model.IncomeCategory, error) {
	var c model.IncomeCategory
	err := r.db.WithContext(ctx).First(&c, id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *incomeCategoryRepo) Create(ctx context.Context, req model.CreateIncomeCategoryRequest) (*model.IncomeCategory, error) {
	c := model.IncomeCategory{Name: req.Name}
	if err := r.db.WithContext(ctx).Create(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *incomeCategoryRepo) Update(ctx context.Context, id int64, req model.UpdateIncomeCategoryRequest) (*model.IncomeCategory, error) {
	result := r.db.WithContext(ctx).Model(&model.IncomeCategory{}).Where("id = ?", id).Update("name", req.Name)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return r.GetByID(ctx, id)
}

func (r *incomeCategoryRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.IncomeCategory{}, id).Error
}
