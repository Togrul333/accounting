package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"accounting/internal/model"
)

type ExpenseRepository interface {
	GetAll(ctx context.Context) ([]model.Expense, error)
	GetByID(ctx context.Context, id int64) (*model.Expense, error)
	GetByAccountID(ctx context.Context, accountID int64) ([]model.Expense, error)
	Create(ctx context.Context, req model.CreateExpenseRequest) (*model.Expense, error)
	BulkCreate(ctx context.Context, reqs []model.CreateExpenseRequest) ([]model.Expense, error)
	Update(ctx context.Context, id int64, req model.UpdateExpenseRequest) (*model.Expense, error)
	Delete(ctx context.Context, id int64) error
}

type expenseRepo struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) ExpenseRepository {
	return &expenseRepo{db: db}
}

const expenseSelectQuery = `
	SELECT e.id, e.name, e.amount, e.date,
	       e.expense_category_id, c.name AS expense_category_name,
	       e.account_id, a.name AS account_name,
	       e.created_at, e.updated_at
	FROM expenses e
	JOIN expense_categories c ON c.id = e.expense_category_id
	JOIN accounts a ON a.id = e.account_id`

func (r *expenseRepo) GetAll(ctx context.Context) ([]model.Expense, error) {
	var expenses []model.Expense
	err := r.db.WithContext(ctx).Raw(expenseSelectQuery + ` ORDER BY e.date DESC, e.id DESC`).Scan(&expenses).Error
	if expenses == nil {
		expenses = []model.Expense{}
	}
	return expenses, err
}

func (r *expenseRepo) GetByID(ctx context.Context, id int64) (*model.Expense, error) {
	var exp model.Expense
	result := r.db.WithContext(ctx).Raw(expenseSelectQuery+` WHERE e.id = ?`, id).Scan(&exp)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &exp, nil
}

func (r *expenseRepo) GetByAccountID(ctx context.Context, accountID int64) ([]model.Expense, error) {
	var expenses []model.Expense
	err := r.db.WithContext(ctx).Raw(expenseSelectQuery+` WHERE e.account_id = ? ORDER BY e.date DESC, e.id DESC`, accountID).Scan(&expenses).Error
	if expenses == nil {
		expenses = []model.Expense{}
	}
	return expenses, err
}

func (r *expenseRepo) Create(ctx context.Context, req model.CreateExpenseRequest) (*model.Expense, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, err
	}
	exp := model.Expense{
		Name:              req.Name,
		Amount:            req.Amount,
		Date:              date,
		ExpenseCategoryID: req.ExpenseCategoryID,
		AccountID:         req.AccountID,
	}
	if err := r.db.WithContext(ctx).Create(&exp).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, exp.ID)
}

func (r *expenseRepo) BulkCreate(ctx context.Context, reqs []model.CreateExpenseRequest) ([]model.Expense, error) {
	var ids []int64
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, req := range reqs {
			date, err := time.Parse("2006-01-02", req.Date)
			if err != nil {
				return err
			}
			exp := model.Expense{
				Name:              req.Name,
				Amount:            req.Amount,
				Date:              date,
				ExpenseCategoryID: req.ExpenseCategoryID,
				AccountID:         req.AccountID,
			}
			if err := tx.Create(&exp).Error; err != nil {
				return err
			}
			ids = append(ids, exp.ID)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	result := make([]model.Expense, 0, len(ids))
	for _, id := range ids {
		exp, err := r.GetByID(ctx, id)
		if err != nil {
			continue
		}
		result = append(result, *exp)
	}
	return result, nil
}

func (r *expenseRepo) Update(ctx context.Context, id int64, req model.UpdateExpenseRequest) (*model.Expense, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, err
	}
	result := r.db.WithContext(ctx).Model(&model.Expense{}).Where("id = ?", id).Updates(map[string]any{
		"name":                req.Name,
		"amount":              req.Amount,
		"date":                date,
		"expense_category_id": req.ExpenseCategoryID,
		"account_id":          req.AccountID,
	})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return r.GetByID(ctx, id)
}

func (r *expenseRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Expense{}, id).Error
}
