package model

import "time"

type IncomeCategory struct {
	ID        int64     `json:"id"`
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
	ID                 int64     `json:"id"`
	Name               string    `json:"name"`
	Amount             float64   `json:"amount"`
	Date               time.Time `json:"date"`
	IncomeCategoryID   int64     `json:"income_category_id"`
	IncomeCategoryName string    `json:"income_category_name,omitempty"`
	AccountID          int64     `json:"account_id"`
	AccountName        string    `json:"account_name,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type CreateIncomeRequest struct {
	Name             string  `json:"name"`
	Amount           float64 `json:"amount"`
	Date             string  `json:"date"`
	IncomeCategoryID int64   `json:"income_category_id"`
	AccountID        int64   `json:"account_id"`
}

type UpdateIncomeRequest struct {
	Name             string  `json:"name"`
	Amount           float64 `json:"amount"`
	Date             string  `json:"date"`
	IncomeCategoryID int64   `json:"income_category_id"`
	AccountID        int64   `json:"account_id"`
}
