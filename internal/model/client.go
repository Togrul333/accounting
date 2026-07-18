package model

import "time"

type Client struct {
	ID            int64     `json:"id" gorm:"primaryKey"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	BirthDate     string    `json:"birth_date"`
	FinCode       string    `json:"fin_code"`
	IDCardNumber  string    `json:"id_card_number"`
	DocumentPhoto string    `json:"document_photo"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateClientRequest struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	BirthDate    string `json:"birth_date"`
	FinCode      string `json:"fin_code"`
	IDCardNumber string `json:"id_card_number"`
}

type UpdateClientRequest struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	BirthDate    string `json:"birth_date"`
	FinCode      string `json:"fin_code"`
	IDCardNumber string `json:"id_card_number"`
}
