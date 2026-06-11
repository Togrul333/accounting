package repository

import (
	"context"

	"gorm.io/gorm"

	"accounting/internal/model"
)

type TourCategoryRepository interface {
	GetAll(ctx context.Context) ([]model.TourCategory, error)
	GetByID(ctx context.Context, id int64) (*model.TourCategory, error)
	Create(ctx context.Context, req model.CreateTourCategoryRequest) (*model.TourCategory, error)
	Update(ctx context.Context, id int64, req model.UpdateTourCategoryRequest) (*model.TourCategory, error)
	Delete(ctx context.Context, id int64) error
}

type tourCategoryRepo struct {
	db *gorm.DB
}

func NewTourCategoryRepository(db *gorm.DB) TourCategoryRepository {
	return &tourCategoryRepo{db: db}
}

func (r *tourCategoryRepo) GetAll(ctx context.Context) ([]model.TourCategory, error) {
	var cats []model.TourCategory
	err := r.db.WithContext(ctx).Order("id").Find(&cats).Error
	return cats, err
}

func (r *tourCategoryRepo) GetByID(ctx context.Context, id int64) (*model.TourCategory, error) {
	var c model.TourCategory
	err := r.db.WithContext(ctx).First(&c, id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *tourCategoryRepo) Create(ctx context.Context, req model.CreateTourCategoryRequest) (*model.TourCategory, error) {
	c := model.TourCategory{Name: req.Name, Price: req.Price}
	if err := r.db.WithContext(ctx).Create(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *tourCategoryRepo) Update(ctx context.Context, id int64, req model.UpdateTourCategoryRequest) (*model.TourCategory, error) {
	result := r.db.WithContext(ctx).Model(&model.TourCategory{}).Where("id = ?", id).Updates(map[string]any{
		"name":  req.Name,
		"price": req.Price,
	})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return r.GetByID(ctx, id)
}

func (r *tourCategoryRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.TourCategory{}, id).Error
}
