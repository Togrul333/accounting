package repository

import (
	"context"

	"gorm.io/gorm"

	"accounting/internal/model"
)

type OrderRepository interface {
	GetAll(ctx context.Context) ([]model.Order, error)
	GetByID(ctx context.Context, id int64) (*model.Order, error)
	Create(ctx context.Context, clientID, tourID int64) (*model.Order, error)
	Update(ctx context.Context, id int64, req model.UpdateOrderRequest) (*model.Order, error)
	Delete(ctx context.Context, id int64) error
}

type orderRepo struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepo{db: db}
}

const orderListQuery = `
	SELECT o.id, o.client_id,
	       CONCAT(c.first_name, ' ', c.last_name) AS client_name,
	       o.tour_id, t.code AS tour_code,
	       COUNT(i.id)              AS income_count,
	       COALESCE(SUM(i.amount), 0) AS income_total,
	       o.created_at, o.updated_at
	FROM orders o
	JOIN clients c ON c.id = o.client_id
	JOIN tours t   ON t.id = o.tour_id
	LEFT JOIN incomes i ON i.order_id = o.id
	GROUP BY o.id, o.client_id, client_name, o.tour_id, t.code, o.created_at, o.updated_at`

const orderBaseQuery = `
	SELECT o.id, o.client_id,
	       CONCAT(c.first_name, ' ', c.last_name) AS client_name,
	       o.tour_id, t.code AS tour_code,
	       o.created_at, o.updated_at
	FROM orders o
	JOIN clients c ON c.id = o.client_id
	JOIN tours t   ON t.id = o.tour_id`

func (r *orderRepo) GetAll(ctx context.Context) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.WithContext(ctx).Raw(orderListQuery + ` ORDER BY o.id DESC`).Scan(&orders).Error
	if orders == nil {
		orders = []model.Order{}
	}
	return orders, err
}

func (r *orderRepo) GetByID(ctx context.Context, id int64) (*model.Order, error) {
	var order model.Order
	result := r.db.WithContext(ctx).Raw(orderBaseQuery+` WHERE o.id = ?`, id).Scan(&order)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &order, nil
}

func (r *orderRepo) Create(ctx context.Context, clientID, tourID int64) (*model.Order, error) {
	order := model.Order{
		ClientID: clientID,
		TourID:   tourID,
	}
	if err := r.db.WithContext(ctx).Create(&order).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, order.ID)
}

func (r *orderRepo) Update(ctx context.Context, id int64, req model.UpdateOrderRequest) (*model.Order, error) {
	result := r.db.WithContext(ctx).Model(&model.Order{}).Where("id = ?", id).Updates(map[string]any{
		"client_id": req.ClientID,
		"tour_id":   req.TourID,
	})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return r.GetByID(ctx, id)
}

func (r *orderRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Order{}, id).Error
}
