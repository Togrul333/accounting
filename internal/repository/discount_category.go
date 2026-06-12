package repository

import (
	"context"

	"gorm.io/gorm"

	"accounting/internal/model"
)

type DiscountCategoryRepository interface {
	GetAll(ctx context.Context) ([]model.DiscountCategory, error)
	GetByID(ctx context.Context, id int64) (*model.DiscountCategory, error)
	Create(ctx context.Context, req model.CreateDiscountCategoryRequest) (*model.DiscountCategory, error)
	Update(ctx context.Context, id int64, req model.UpdateDiscountCategoryRequest) (*model.DiscountCategory, error)
	Delete(ctx context.Context, id int64) error
}

type discountCategoryRepo struct {
	db *gorm.DB
}

func NewDiscountCategoryRepository(db *gorm.DB) DiscountCategoryRepository {
	return &discountCategoryRepo{db: db}
}

func (r *discountCategoryRepo) GetAll(ctx context.Context) ([]model.DiscountCategory, error) {
	var cats []model.DiscountCategory
	err := r.db.WithContext(ctx).Order("id").Find(&cats).Error
	return cats, err
}

func (r *discountCategoryRepo) GetByID(ctx context.Context, id int64) (*model.DiscountCategory, error) {
	var c model.DiscountCategory
	err := r.db.WithContext(ctx).First(&c, id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *discountCategoryRepo) Create(ctx context.Context, req model.CreateDiscountCategoryRequest) (*model.DiscountCategory, error) {
	c := model.DiscountCategory{Name: req.Name}
	if err := r.db.WithContext(ctx).Create(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *discountCategoryRepo) Update(ctx context.Context, id int64, req model.UpdateDiscountCategoryRequest) (*model.DiscountCategory, error) {
	result := r.db.WithContext(ctx).Model(&model.DiscountCategory{}).Where("id = ?", id).Update("name", req.Name)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return r.GetByID(ctx, id)
}

func (r *discountCategoryRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.DiscountCategory{}, id).Error
}
