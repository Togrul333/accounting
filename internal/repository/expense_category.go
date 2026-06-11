package repository

import (
	"context"
	"database/sql"

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
	db *sql.DB
}

func NewExpenseCategoryRepository(db *sql.DB) ExpenseCategoryRepository {
	return &expenseCategoryRepo{db: db}
}

func (r *expenseCategoryRepo) GetAll(ctx context.Context) ([]model.ExpenseCategory, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, created_at, updated_at FROM expense_categories ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []model.ExpenseCategory
	for rows.Next() {
		var c model.ExpenseCategory
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, rows.Err()
}

func (r *expenseCategoryRepo) GetByID(ctx context.Context, id int64) (*model.ExpenseCategory, error) {
	var c model.ExpenseCategory
	err := r.db.QueryRowContext(ctx, `SELECT id, name, created_at, updated_at FROM expense_categories WHERE id=?`, id).
		Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *expenseCategoryRepo) Create(ctx context.Context, req model.CreateExpenseCategoryRequest) (*model.ExpenseCategory, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO expense_categories (name) VALUES (?)`, req.Name)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *expenseCategoryRepo) Update(ctx context.Context, id int64, req model.UpdateExpenseCategoryRequest) (*model.ExpenseCategory, error) {
	res, err := r.db.ExecContext(ctx, `UPDATE expense_categories SET name=? WHERE id=?`, req.Name, id)
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

func (r *expenseCategoryRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM expense_categories WHERE id=?`, id)
	return err
}
