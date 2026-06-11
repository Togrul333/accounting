package repository

import (
	"context"
	"database/sql"

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
	db *sql.DB
}

func NewTourCategoryRepository(db *sql.DB) TourCategoryRepository {
	return &tourCategoryRepo{db: db}
}

func (r *tourCategoryRepo) GetAll(ctx context.Context) ([]model.TourCategory, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, price, created_at, updated_at FROM tour_categories ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []model.TourCategory
	for rows.Next() {
		var c model.TourCategory
		if err := rows.Scan(&c.ID, &c.Name, &c.Price, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, rows.Err()
}

func (r *tourCategoryRepo) GetByID(ctx context.Context, id int64) (*model.TourCategory, error) {
	var c model.TourCategory
	err := r.db.QueryRowContext(ctx, `SELECT id, name, price, created_at, updated_at FROM tour_categories WHERE id=?`, id).
		Scan(&c.ID, &c.Name, &c.Price, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *tourCategoryRepo) Create(ctx context.Context, req model.CreateTourCategoryRequest) (*model.TourCategory, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO tour_categories (name, price) VALUES (?, ?)`, req.Name, req.Price)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *tourCategoryRepo) Update(ctx context.Context, id int64, req model.UpdateTourCategoryRequest) (*model.TourCategory, error) {
	res, err := r.db.ExecContext(ctx, `UPDATE tour_categories SET name=?, price=? WHERE id=?`, req.Name, req.Price, id)
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

func (r *tourCategoryRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM tour_categories WHERE id=?`, id)
	return err
}
