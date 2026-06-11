package repository

import (
	"context"

	"gorm.io/gorm"

	"accounting/internal/model"
)

type ExpenseCategoryRepository interface {
	GetAll(ctx context.Context) ([]model.ExpenseCategory, error)
	GetByID(ctx context.Context, id int64) (*model.ExpenseCategory, error)
	Create(ctx context.Context, req model.CreateExpenseCategoryRequest) (*model.ExpenseCategory, error)
	Update(ctx context.Context, id int64, req model.UpdateExpenseCategoryRequest) (*model.ExpenseCategory, error)
	Delete(ctx context.Context, id int64) error
}

type expenseCategoryRepo struct {
	db *gorm.DB
}

func NewExpenseCategoryRepository(db *gorm.DB) ExpenseCategoryRepository {
	return &expenseCategoryRepo{db: db}
}

func (r *expenseCategoryRepo) GetAll(ctx context.Context) ([]model.ExpenseCategory, error) {
	var cats []model.ExpenseCategory
	err := r.db.WithContext(ctx).Order("id").Find(&cats).Error
	return cats, err
}

func (r *expenseCategoryRepo) GetByID(ctx context.Context, id int64) (*model.ExpenseCategory, error) {
	var c model.ExpenseCategory
	err := r.db.WithContext(ctx).First(&c, id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *expenseCategoryRepo) Create(ctx context.Context, req model.CreateExpenseCategoryRequest) (*model.ExpenseCategory, error) {
	c := model.ExpenseCategory{Name: req.Name}
	if err := r.db.WithContext(ctx).Create(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *expenseCategoryRepo) Update(ctx context.Context, id int64, req model.UpdateExpenseCategoryRequest) (*model.ExpenseCategory, error) {
	result := r.db.WithContext(ctx).Model(&model.ExpenseCategory{}).Where("id = ?", id).Update("name", req.Name)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return r.GetByID(ctx, id)
}

func (r *expenseCategoryRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.ExpenseCategory{}, id).Error
}
