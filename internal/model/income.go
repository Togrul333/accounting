package model

import "time"

type IncomeCategory struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateIncomeCategoryRequest struct {
	Name string `json:"name"`
}

type UpdateIncomeCategoryRequest struct {
	Name string `json:"name"`
}

type Income struct {
	ID                 int64     `json:"id" gorm:"primaryKey"`
	Name               string    `json:"name"`
	Amount             float64   `json:"amount"`
	Date               time.Time `json:"date"`
	IncomeCategoryID   int64     `json:"income_category_id"`
	IncomeCategoryName string    `json:"income_category_name,omitempty" gorm:"<-:false"`
	AccountID          int64     `json:"account_id"`
	AccountName        string    `json:"account_name,omitempty" gorm:"<-:false"`
	TourID             *int64    `json:"tour_id"`
	TourCode           string    `json:"tour_code,omitempty" gorm:"<-:false"`
	OrderID            *int64    `json:"order_id"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type CreateIncomeRequest struct {
	Name             string  `json:"name"`
	Amount           float64 `json:"amount"`
	Date             string  `json:"date"`
	IncomeCategoryID int64   `json:"income_category_id"`
	AccountID        int64   `json:"account_id"`
	TourID           *int64  `json:"tour_id"`
	OrderID          *int64  `json:"order_id"`
}

type UpdateIncomeRequest struct {
	Name             string  `json:"name"`
	Amount           float64 `json:"amount"`
	Date             string  `json:"date"`
	IncomeCategoryID int64   `json:"income_category_id"`
	AccountID        int64   `json:"account_id"`
	TourID           *int64  `json:"tour_id"`
	OrderID          *int64  `json:"order_id"`
}
