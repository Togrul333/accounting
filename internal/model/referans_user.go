package model

import "time"

// ReferansUser — реферер: человек, который приводит клиентов.
type ReferansUser struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ReferansUser) TableName() string {
	return "referans_users"
}

type CreateReferansUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

// ReferansOrder — строка заказа в списках кандидатов и подтверждённых рефералов.
type ReferansOrder struct {
	ID            int64     `json:"id"`
	ClientID      int64     `json:"client_id"`
	ClientName    string    `json:"client_name"`
	ReferenceName string    `json:"reference_name"`
	TourCode      string    `json:"tour_code"`
	IncomeTotal   float64   `json:"income_total"`
	Confirmed     bool      `json:"confirmed"`
	CreatedAt     time.Time `json:"created_at"`
}

type AddReferansOrderRequest struct {
	OrderID int64 `json:"order_id"`
}

type UpdateReferansUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}
