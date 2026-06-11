package model

import "time"

type Account struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	AccountNumber string    `json:"account_number"`
	Currency      string    `json:"currency"`
	Balance       float64   `json:"balance"`
	Description   string    `json:"description,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateAccountRequest struct {
	Name          string `json:"name"`
	AccountNumber string `json:"account_number"`
	Currency      string `json:"currency"`
	Description   string `json:"description"`
}

type UpdateAccountRequest struct {
	Name          string `json:"name"`
	AccountNumber string `json:"account_number"`
	Currency      string `json:"currency"`
	Description   string `json:"description"`
}
