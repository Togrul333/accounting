package repository

import (
	"context"
	"database/sql"

	"accounting/internal/model"
)

type ClientRepository interface {
	GetAll(ctx context.Context) ([]model.Client, error)
	GetByID(ctx context.Context, id int64) (*model.Client, error)
	Create(ctx context.Context, req model.CreateClientRequest) (*model.Client, error)
	Update(ctx context.Context, id int64, req model.UpdateClientRequest) (*model.Client, error)
	Delete(ctx context.Context, id int64) error
}

type clientRepo struct {
	db *sql.DB
}

func NewClientRepository(db *sql.DB) ClientRepository {
	return &clientRepo{db: db}
}

const clientCols = `id, first_name, last_name, email, phone, birth_year, created_at, updated_at`

func scanClient(s interface{ Scan(...any) error }) (model.Client, error) {
	var c model.Client
	err := s.Scan(&c.ID, &c.FirstName, &c.LastName, &c.Email, &c.Phone, &c.BirthYear, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *clientRepo) GetAll(ctx context.Context) ([]model.Client, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+clientCols+` FROM clients ORDER BY last_name, first_name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []model.Client
	for rows.Next() {
		c, err := scanClient(rows)
		if err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}
	return clients, rows.Err()
}

func (r *clientRepo) GetByID(ctx context.Context, id int64) (*model.Client, error) {
	c, err := scanClient(r.db.QueryRowContext(ctx, `SELECT `+clientCols+` FROM clients WHERE id=?`, id))
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *clientRepo) Create(ctx context.Context, req model.CreateClientRequest) (*model.Client, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO clients (first_name, last_name, email, phone, birth_year) VALUES (?, ?, ?, ?, ?)`,
		req.FirstName, req.LastName, req.Email, req.Phone, req.BirthYear,
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

func (r *clientRepo) Update(ctx context.Context, id int64, req model.UpdateClientRequest) (*model.Client, error) {
	res, err := r.db.ExecContext(ctx,
		`UPDATE clients SET first_name=?, last_name=?, email=?, phone=?, birth_year=? WHERE id=?`,
		req.FirstName, req.LastName, req.Email, req.Phone, req.BirthYear, id,
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

func (r *clientRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM clients WHERE id=?`, id)
	return err
}
