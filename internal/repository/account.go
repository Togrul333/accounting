package repository

import (
	"context"

	"gorm.io/gorm"

	"accounting/internal/model"
)

type AccountRepository interface {
	GetAll(ctx context.Context) ([]model.Account, error)
	GetByID(ctx context.Context, id int64) (*model.Account, error)
	Create(ctx context.Context, req model.CreateAccountRequest) (*model.Account, error)
	Update(ctx context.Context, id int64, req model.UpdateAccountRequest) (*model.Account, error)
	Delete(ctx context.Context, id int64) error
}

type accountRepo struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepo{db: db}
}

const accountBalanceQuery = `
	SELECT a.id, a.name, a.account_number, a.currency,
	  COALESCE((SELECT SUM(amount) FROM incomes WHERE account_id = a.id), 0) -
	  COALESCE((SELECT SUM(amount) FROM expenses WHERE account_id = a.id), 0) AS balance,
	  a.description, a.created_at, a.updated_at
	FROM accounts a`

func (r *accountRepo) GetAll(ctx context.Context) ([]model.Account, error) {
	var accounts []model.Account
	err := r.db.WithContext(ctx).Raw(accountBalanceQuery + ` ORDER BY a.id`).Scan(&accounts).Error
	if accounts == nil {
		accounts = []model.Account{}
	}
	return accounts, err
}

func (r *accountRepo) GetByID(ctx context.Context, id int64) (*model.Account, error) {
	var a model.Account
	result := r.db.WithContext(ctx).Raw(accountBalanceQuery+` WHERE a.id = ?`, id).Scan(&a)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &a, nil
}

func (r *accountRepo) Create(ctx context.Context, req model.CreateAccountRequest) (*model.Account, error) {
	a := model.Account{
		Name:          req.Name,
		AccountNumber: req.AccountNumber,
		Currency:      req.Currency,
		Description:   req.Description,
	}
	if err := r.db.WithContext(ctx).Create(&a).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, a.ID)
}

func (r *accountRepo) Update(ctx context.Context, id int64, req model.UpdateAccountRequest) (*model.Account, error) {
	result := r.db.WithContext(ctx).Model(&model.Account{}).Where("id = ?", id).Updates(map[string]any{
		"name":           req.Name,
		"account_number": req.AccountNumber,
		"currency":       req.Currency,
		"description":    req.Description,
	})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return r.GetByID(ctx, id)
}

func (r *accountRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Account{}, id).Error
}
