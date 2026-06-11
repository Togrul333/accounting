package repository

import (
	"context"
	"database/sql"

	"accounting/internal/model"
)

type RoomRepository interface {
	GetAll(ctx context.Context) ([]model.Room, error)
	GetByID(ctx context.Context, id int64) (*model.Room, error)
	Create(ctx context.Context, req model.CreateRoomRequest) (*model.Room, error)
	Update(ctx context.Context, id int64, req model.UpdateRoomRequest) (*model.Room, error)
	Delete(ctx context.Context, id int64) error
}

type roomRepo struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) RoomRepository {
	return &roomRepo{db: db}
}

func (r *roomRepo) GetAll(ctx context.Context) ([]model.Room, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, price, beds_count, created_at, updated_at FROM rooms ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []model.Room
	for rows.Next() {
		var rm model.Room
		if err := rows.Scan(&rm.ID, &rm.Price, &rm.BedsCount, &rm.CreatedAt, &rm.UpdatedAt); err != nil {
			return nil, err
		}
		rooms = append(rooms, rm)
	}
	return rooms, rows.Err()
}

func (r *roomRepo) GetByID(ctx context.Context, id int64) (*model.Room, error) {
	var rm model.Room
	err := r.db.QueryRowContext(ctx, `SELECT id, price, beds_count, created_at, updated_at FROM rooms WHERE id=?`, id).
		Scan(&rm.ID, &rm.Price, &rm.BedsCount, &rm.CreatedAt, &rm.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &rm, nil
}

func (r *roomRepo) Create(ctx context.Context, req model.CreateRoomRequest) (*model.Room, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO rooms (price, beds_count) VALUES (?, ?)`, req.Price, req.BedsCount)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *roomRepo) Update(ctx context.Context, id int64, req model.UpdateRoomRequest) (*model.Room, error) {
	res, err := r.db.ExecContext(ctx, `UPDATE rooms SET price=?, beds_count=? WHERE id=?`, req.Price, req.BedsCount, id)
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

func (r *roomRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM rooms WHERE id=?`, id)
	return err
}
