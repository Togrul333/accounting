package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"accounting/internal/model"
)

type IncomeRepository interface {
	GetAll(ctx context.Context) ([]model.Income, error)
	GetByID(ctx context.Context, id int64) (*model.Income, error)
	GetByAccountID(ctx context.Context, accountID int64) ([]model.Income, error)
	GetByOrderID(ctx context.Context, orderID int64) ([]model.Income, error)
	Create(ctx context.Context, req model.CreateIncomeRequest) (*model.Income, error)
	BulkCreate(ctx context.Context, reqs []model.CreateIncomeRequest) ([]model.Income, error)
	Update(ctx context.Context, id int64, req model.UpdateIncomeRequest) (*model.Income, error)
	Delete(ctx context.Context, id int64) error
}

type incomeRepo struct {
	db *gorm.DB
}

func NewIncomeRepository(db *gorm.DB) IncomeRepository {
	return &incomeRepo{db: db}
}

const incomeSelectQuery = `
	SELECT i.id, i.name, i.amount, i.date,
	       i.income_category_id, c.name AS income_category_name,
	       i.account_id, a.name AS account_name,
	       i.tour_id, t.code AS tour_code,
	       i.order_id,
	       i.created_at, i.updated_at
	FROM incomes i
	JOIN income_categories c ON c.id = i.income_category_id
	JOIN accounts a ON a.id = i.account_id
	LEFT JOIN tours t ON t.id = i.tour_id`

func (r *incomeRepo) GetAll(ctx context.Context) ([]model.Income, error) {
	var incomes []model.Income
	err := r.db.WithContext(ctx).Raw(incomeSelectQuery + ` ORDER BY i.date DESC, i.id DESC`).Scan(&incomes).Error
	if incomes == nil {
		incomes = []model.Income{}
	}
	return incomes, err
}

func (r *incomeRepo) GetByID(ctx context.Context, id int64) (*model.Income, error) {
	var inc model.Income
	result := r.db.WithContext(ctx).Raw(incomeSelectQuery+` WHERE i.id = ?`, id).Scan(&inc)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &inc, nil
}

func (r *incomeRepo) GetByAccountID(ctx context.Context, accountID int64) ([]model.Income, error) {
	var incomes []model.Income
	err := r.db.WithContext(ctx).Raw(incomeSelectQuery+` WHERE i.account_id = ? ORDER BY i.date DESC, i.id DESC`, accountID).Scan(&incomes).Error
	if incomes == nil {
		incomes = []model.Income{}
	}
	return incomes, err
}

func (r *incomeRepo) GetByOrderID(ctx context.Context, orderID int64) ([]model.Income, error) {
	var incomes []model.Income
	err := r.db.WithContext(ctx).Raw(incomeSelectQuery+` WHERE i.order_id = ? ORDER BY i.date DESC, i.id DESC`, orderID).Scan(&incomes).Error
	if incomes == nil {
		incomes = []model.Income{}
	}
	return incomes, err
}

func (r *incomeRepo) Create(ctx context.Context, req model.CreateIncomeRequest) (*model.Income, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, err
	}
	inc := model.Income{
		Name:             req.Name,
		Amount:           req.Amount,
		Date:             date,
		IncomeCategoryID: req.IncomeCategoryID,
		AccountID:        req.AccountID,
		TourID:           req.TourID,
		OrderID:          req.OrderID,
	}
	if err := r.db.WithContext(ctx).Create(&inc).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, inc.ID)
}

func (r *incomeRepo) BulkCreate(ctx context.Context, reqs []model.CreateIncomeRequest) ([]model.Income, error) {
	var ids []int64
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, req := range reqs {
			date, err := time.Parse("2006-01-02", req.Date)
			if err != nil {
				return err
			}
			inc := model.Income{
				Name:             req.Name,
				Amount:           req.Amount,
				Date:             date,
				IncomeCategoryID: req.IncomeCategoryID,
				AccountID:        req.AccountID,
				TourID:           req.TourID,
				OrderID:          req.OrderID,
			}
			if err := tx.Create(&inc).Error; err != nil {
				return err
			}
			ids = append(ids, inc.ID)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	result := make([]model.Income, 0, len(ids))
	for _, id := range ids {
		inc, err := r.GetByID(ctx, id)
		if err != nil {
			continue
		}
		result = append(result, *inc)
	}
	return result, nil
}

func (r *incomeRepo) Update(ctx context.Context, id int64, req model.UpdateIncomeRequest) (*model.Income, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, err
	}
	result := r.db.WithContext(ctx).Model(&model.Income{}).Where("id = ?", id).Updates(map[string]any{
		"name":               req.Name,
		"amount":             req.Amount,
		"date":               date,
		"income_category_id": req.IncomeCategoryID,
		"account_id":         req.AccountID,
		"tour_id":            req.TourID,
		"order_id":           req.OrderID,
	})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return r.GetByID(ctx, id)
}

func (r *incomeRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Income{}, id).Error
}
