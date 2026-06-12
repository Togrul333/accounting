package repository

import (
	"context"

	"gorm.io/gorm"

	"accounting/internal/model"
)

type DiscountRepository interface {
	GetAll(ctx context.Context) ([]model.Discount, error)
	GetByID(ctx context.Context, id int64) (*model.Discount, error)
	Create(ctx context.Context, req model.CreateDiscountRequest) (*model.Discount, error)
	Update(ctx context.Context, id int64, req model.UpdateDiscountRequest) (*model.Discount, error)
	Delete(ctx context.Context, id int64) error
}

type discountRepo struct {
	db *gorm.DB
}

func NewDiscountRepository(db *gorm.DB) DiscountRepository {
	return &discountRepo{db: db}
}

const discountSelectQuery = `
	SELECT d.id, d.amount, d.discount_category_id, c.name AS discount_category_name,
	       d.created_at, d.updated_at
	FROM discounts d
	JOIN discount_categories c ON c.id = d.discount_category_id`

func (r *discountRepo) GetAll(ctx context.Context) ([]model.Discount, error) {
	var discounts []model.Discount
	err := r.db.WithContext(ctx).Raw(discountSelectQuery + ` ORDER BY d.id DESC`).Scan(&discounts).Error
	if discounts == nil {
		discounts = []model.Discount{}
	}
	return discounts, err
}

func (r *discountRepo) GetByID(ctx context.Context, id int64) (*model.Discount, error) {
	var d model.Discount
	result := r.db.WithContext(ctx).Raw(discountSelectQuery+` WHERE d.id = ?`, id).Scan(&d)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &d, nil
}

func (r *discountRepo) Create(ctx context.Context, req model.CreateDiscountRequest) (*model.Discount, error) {
	d := model.Discount{
		Amount:             req.Amount,
		DiscountCategoryID: req.DiscountCategoryID,
	}
	if err := r.db.WithContext(ctx).Create(&d).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, d.ID)
}

func (r *discountRepo) Update(ctx context.Context, id int64, req model.UpdateDiscountRequest) (*model.Discount, error) {
	result := r.db.WithContext(ctx).Model(&model.Discount{}).Where("id = ?", id).Updates(map[string]any{
		"amount":               req.Amount,
		"discount_category_id": req.DiscountCategoryID,
	})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return r.GetByID(ctx, id)
}

func (r *discountRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Discount{}, id).Error
}
