package repository

import (
	"context"
	"database/sql"

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
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &accountRepo{db: db}
}

func (r *accountRepo) GetAll(ctx context.Context) ([]model.Account, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT a.id, a.name, a.account_number, a.currency,
		  COALESCE((SELECT SUM(amount) FROM incomes WHERE account_id = a.id), 0) -
		  COALESCE((SELECT SUM(amount) FROM expenses WHERE account_id = a.id), 0) AS balance,
		  a.description, a.created_at, a.updated_at
		FROM accounts a ORDER BY a.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []model.Account
	for rows.Next() {
		var a model.Account
		if err := rows.Scan(&a.ID, &a.Name, &a.AccountNumber, &a.Currency, &a.Balance, &a.Description, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	return accounts, rows.Err()
}

func (r *accountRepo) GetByID(ctx context.Context, id int64) (*model.Account, error) {
	var a model.Account
	err := r.db.QueryRowContext(ctx, `
		SELECT a.id, a.name, a.account_number, a.currency,
		  COALESCE((SELECT SUM(amount) FROM incomes WHERE account_id = a.id), 0) -
		  COALESCE((SELECT SUM(amount) FROM expenses WHERE account_id = a.id), 0) AS balance,
		  a.description, a.created_at, a.updated_at
		FROM accounts a WHERE a.id = ?
	`, id).Scan(&a.ID, &a.Name, &a.AccountNumber, &a.Currency, &a.Balance, &a.Description, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *accountRepo) Create(ctx context.Context, req model.CreateAccountRequest) (*model.Account, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO accounts (name, account_number, currency, description)
		VALUES (?, ?, ?, ?)
	`, req.Name, req.AccountNumber, req.Currency, req.Description)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *accountRepo) Update(ctx context.Context, id int64, req model.UpdateAccountRequest) (*model.Account, error) {
	res, err := r.db.ExecContext(ctx, `
		UPDATE accounts
		SET name=?, account_number=?, currency=?, description=?
		WHERE id=?
	`, req.Name, req.AccountNumber, req.Currency, req.Description, id)
	if err != nil {
		return nil, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, sql.ErrNoRows
	}
	return r.GetByID(ctx, id)
}

func (r *accountRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM accounts WHERE id=?`, id)
	return err
}
