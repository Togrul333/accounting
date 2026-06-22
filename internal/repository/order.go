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
	       tc.name AS tour_category_name,
	       tc.price AS tour_category_price,
	       r.price AS room_price,
	       COALESCE(inc.income_count, 0)   AS income_count,
	       COALESCE(inc.income_total, 0)   AS income_total,
	       COALESCE(disc.discount_count, 0) AS discount_count,
	       COALESCE(disc.discount_total, 0) AS discount_total,
	       o.created_at, o.updated_at
	FROM orders o
	JOIN clients c ON c.id = o.client_id
	JOIN tours t   ON t.id = o.tour_id
	JOIN tour_categories tc ON tc.id = t.tour_category_id
	JOIN rooms r ON r.id = t.room_id
	LEFT JOIN (
	    SELECT order_id, COUNT(*) AS income_count, SUM(amount) AS income_total
	    FROM incomes GROUP BY order_id
	) inc ON inc.order_id = o.id
	LEFT JOIN (
	    SELECT order_id, COUNT(*) AS discount_count, SUM(amount) AS discount_total
	    FROM discounts GROUP BY order_id
	) disc ON disc.order_id = o.id`

const orderBaseQuery = `
	SELECT o.id, o.client_id,
	       CONCAT(c.first_name, ' ', c.last_name) AS client_name,
	       o.tour_id, t.code AS tour_code,
	       tc.name AS tour_category_name,
	       tc.price AS tour_category_price,
	       r.price AS room_price,
	       o.created_at, o.updated_at
	FROM orders o
	JOIN clients c ON c.id = o.client_id
	JOIN tours t   ON t.id = o.tour_id
	JOIN tour_categories tc ON tc.id = t.tour_category_id
	JOIN rooms r ON r.id = t.room_id`

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
