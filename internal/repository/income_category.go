package repository

import (
	"context"
	"database/sql"

	"accounting/internal/model"
)

type IncomeCategoryRepository interface {
	GetAll(ctx context.Context) ([]model.IncomeCategory, error)
	GetByID(ctx context.Context, id int64) (*model.IncomeCategory, error)
	Create(ctx context.Context, req model.CreateIncomeCategoryRequest) (*model.IncomeCategory, error)
	Update(ctx context.Context, id int64, req model.UpdateIncomeCategoryRequest) (*model.IncomeCategory, error)
	Delete(ctx context.Context, id int64) error
}

type incomeCategoryRepo struct {
	db *sql.DB
}

func NewIncomeCategoryRepository(db *sql.DB) IncomeCategoryRepository {
	return &incomeCategoryRepo{db: db}
}

func (r *incomeCategoryRepo) GetAll(ctx context.Context) ([]model.IncomeCategory, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, created_at, updated_at FROM income_categories ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []model.IncomeCategory
	for rows.Next() {
		var c model.IncomeCategory
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, rows.Err()
}

func (r *incomeCategoryRepo) GetByID(ctx context.Context, id int64) (*model.IncomeCategory, error) {
	var c model.IncomeCategory
	err := r.db.QueryRowContext(ctx, `SELECT id, name, created_at, updated_at FROM income_categories WHERE id=?`, id).
		Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *incomeCategoryRepo) Create(ctx context.Context, req model.CreateIncomeCategoryRequest) (*model.IncomeCategory, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO income_categories (name) VALUES (?)`, req.Name)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *incomeCategoryRepo) Update(ctx context.Context, id int64, req model.UpdateIncomeCategoryRequest) (*model.IncomeCategory, error) {
	res, err := r.db.ExecContext(ctx, `UPDATE income_categories SET name=? WHERE id=?`, req.Name, id)
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

func (r *incomeCategoryRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM income_categories WHERE id=?`, id)
	return err
}
