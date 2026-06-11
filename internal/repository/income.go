package repository

import (
	"context"
	"database/sql"

	"accounting/internal/model"
)

type IncomeRepository interface {
	GetAll(ctx context.Context) ([]model.Income, error)
	GetByID(ctx context.Context, id int64) (*model.Income, error)
	GetByAccountID(ctx context.Context, accountID int64) ([]model.Income, error)
	Create(ctx context.Context, req model.CreateIncomeRequest) (*model.Income, error)
	BulkCreate(ctx context.Context, reqs []model.CreateIncomeRequest) ([]model.Income, error)
	Update(ctx context.Context, id int64, req model.UpdateIncomeRequest) (*model.Income, error)
	Delete(ctx context.Context, id int64) error
}

type incomeRepo struct {
	db *sql.DB
}

func NewIncomeRepository(db *sql.DB) IncomeRepository {
	return &incomeRepo{db: db}
}

func (r *incomeRepo) GetAll(ctx context.Context) ([]model.Income, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT i.id, i.name, i.amount, i.date, i.income_category_id, c.name, i.account_id, a.name, i.created_at, i.updated_at
		FROM incomes i
		JOIN income_categories c ON c.id = i.income_category_id
		JOIN accounts a ON a.id = i.account_id
		ORDER BY i.date DESC, i.id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incomes []model.Income
	for rows.Next() {
		var inc model.Income
		if err := rows.Scan(&inc.ID, &inc.Name, &inc.Amount, &inc.Date, &inc.IncomeCategoryID, &inc.IncomeCategoryName, &inc.AccountID, &inc.AccountName, &inc.CreatedAt, &inc.UpdatedAt); err != nil {
			return nil, err
		}
		incomes = append(incomes, inc)
	}
	return incomes, rows.Err()
}

func (r *incomeRepo) GetByID(ctx context.Context, id int64) (*model.Income, error) {
	var inc model.Income
	err := r.db.QueryRowContext(ctx, `
		SELECT i.id, i.name, i.amount, i.date, i.income_category_id, c.name, i.account_id, a.name, i.created_at, i.updated_at
		FROM incomes i
		JOIN income_categories c ON c.id = i.income_category_id
		JOIN accounts a ON a.id = i.account_id
		WHERE i.id=?
	`, id).Scan(&inc.ID, &inc.Name, &inc.Amount, &inc.Date, &inc.IncomeCategoryID, &inc.IncomeCategoryName, &inc.AccountID, &inc.AccountName, &inc.CreatedAt, &inc.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &inc, nil
}

func (r *incomeRepo) GetByAccountID(ctx context.Context, accountID int64) ([]model.Income, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT i.id, i.name, i.amount, i.date, i.income_category_id, c.name, i.account_id, a.name, i.created_at, i.updated_at
		FROM incomes i
		JOIN income_categories c ON c.id = i.income_category_id
		JOIN accounts a ON a.id = i.account_id
		WHERE i.account_id = ?
		ORDER BY i.date DESC, i.id DESC
	`, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incomes []model.Income
	for rows.Next() {
		var inc model.Income
		if err := rows.Scan(&inc.ID, &inc.Name, &inc.Amount, &inc.Date, &inc.IncomeCategoryID, &inc.IncomeCategoryName, &inc.AccountID, &inc.AccountName, &inc.CreatedAt, &inc.UpdatedAt); err != nil {
			return nil, err
		}
		incomes = append(incomes, inc)
	}
	return incomes, rows.Err()
}

func (r *incomeRepo) Create(ctx context.Context, req model.CreateIncomeRequest) (*model.Income, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO incomes (name, amount, date, income_category_id, account_id) VALUES (?, ?, ?, ?, ?)
	`, req.Name, req.Amount, req.Date, req.IncomeCategoryID, req.AccountID)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *incomeRepo) BulkCreate(ctx context.Context, reqs []model.CreateIncomeRequest) ([]model.Income, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	ids := make([]int64, 0, len(reqs))
	for _, req := range reqs {
		res, err := tx.ExecContext(ctx,
			`INSERT INTO incomes (name, amount, date, income_category_id, account_id) VALUES (?, ?, ?, ?, ?)`,
			req.Name, req.Amount, req.Date, req.IncomeCategoryID, req.AccountID)
		if err != nil {
			return nil, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := tx.Commit(); err != nil {
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
	res, err := r.db.ExecContext(ctx, `
		UPDATE incomes SET name=?, amount=?, date=?, income_category_id=?, account_id=? WHERE id=?
	`, req.Name, req.Amount, req.Date, req.IncomeCategoryID, req.AccountID, id)
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

func (r *incomeRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM incomes WHERE id=?`, id)
	return err
}
