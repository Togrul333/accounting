package repository

import (
	"context"
	"database/sql"

	"accounting/internal/model"
)

type TourRepository interface {
	GetAll(ctx context.Context) ([]model.Tour, error)
	GetByID(ctx context.Context, id int64) (*model.Tour, error)
	Create(ctx context.Context, req model.CreateTourRequest) (*model.Tour, error)
	Update(ctx context.Context, id int64, req model.UpdateTourRequest) (*model.Tour, error)
	Delete(ctx context.Context, id int64) error
}

type tourRepo struct {
	db *sql.DB
}

func NewTourRepository(db *sql.DB) TourRepository {
	return &tourRepo{db: db}
}

const tourSelectQuery = `
SELECT t.id, t.code, t.start_date, t.end_date,
       t.tour_category_id, tc.name, tc.price,
       t.room_id, r.price, r.beds_count,
       t.created_at, t.updated_at
FROM tours t
JOIN tour_categories tc ON tc.id = t.tour_category_id
JOIN rooms r ON r.id = t.room_id`

func scanTour(s interface{ Scan(...any) error }) (model.Tour, error) {
	var t model.Tour
	err := s.Scan(
		&t.ID, &t.Code, &t.StartDate, &t.EndDate,
		&t.TourCategoryID, &t.TourCategoryName, &t.TourCategoryPrice,
		&t.RoomID, &t.RoomPrice, &t.RoomBedsCount,
		&t.CreatedAt, &t.UpdatedAt,
	)
	return t, err
}

func (r *tourRepo) GetAll(ctx context.Context) ([]model.Tour, error) {
	rows, err := r.db.QueryContext(ctx, tourSelectQuery+` ORDER BY t.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tours []model.Tour
	for rows.Next() {
		t, err := scanTour(rows)
		if err != nil {
			return nil, err
		}
		tours = append(tours, t)
	}
	return tours, rows.Err()
}

func (r *tourRepo) GetByID(ctx context.Context, id int64) (*model.Tour, error) {
	t, err := scanTour(r.db.QueryRowContext(ctx, tourSelectQuery+` WHERE t.id=?`, id))
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *tourRepo) Create(ctx context.Context, req model.CreateTourRequest) (*model.Tour, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO tours (code, start_date, end_date, tour_category_id, room_id) VALUES (?, ?, ?, ?, ?)`,
		req.Code, req.StartDate, req.EndDate, req.TourCategoryID, req.RoomID,
	)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *tourRepo) Update(ctx context.Context, id int64, req model.UpdateTourRequest) (*model.Tour, error) {
	res, err := r.db.ExecContext(ctx,
		`UPDATE tours SET code=?, start_date=?, end_date=?, tour_category_id=?, room_id=? WHERE id=?`,
		req.Code, req.StartDate, req.EndDate, req.TourCategoryID, req.RoomID, id,
	)
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

func (r *tourRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM tours WHERE id=?`, id)
	return err
}
