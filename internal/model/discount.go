package model

import "time"

type DiscountCategory struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateDiscountCategoryRequest struct {
	Name string `json:"name"`
}

type UpdateDiscountCategoryRequest struct {
	Name string `json:"name"`
}

type Discount struct {
	ID                   int64     `json:"id" gorm:"primaryKey"`
	Amount               float64   `json:"amount"`
	DiscountCategoryID   int64     `json:"discount_category_id"`
	DiscountCategoryName string    `json:"discount_category_name,omitempty" gorm:"<-:false"`
	OrderID              *int64    `json:"order_id"`
	OrderClientName      string    `json:"order_client_name,omitempty" gorm:"<-:false"`
	OrderTourCode        string    `json:"order_tour_code,omitempty" gorm:"<-:false"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type CreateDiscountRequest struct {
	Amount             float64 `json:"amount"`
	DiscountCategoryID int64   `json:"discount_category_id"`
	OrderID            *int64  `json:"order_id"`
}

type UpdateDiscountRequest struct {
	Amount             float64 `json:"amount"`
	DiscountCategoryID int64   `json:"discount_category_id"`
	OrderID            *int64  `json:"order_id"`
}
