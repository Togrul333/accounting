package model

import "time"

type Flight struct {
	ID            int64     `json:"id" gorm:"primaryKey"`
	PNR           string    `json:"pnr"`
	DepartureTime time.Time `json:"departure_time"`
	Price         float64   `json:"price"`
	Deposit       float64   `json:"deposit"`
	PaxCount      int       `json:"pax_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateFlightRequest struct {
	PNR           string  `json:"pnr"`
	DepartureTime string  `json:"departure_time"`
	Price         float64 `json:"price"`
	Deposit       float64 `json:"deposit"`
	PaxCount      int     `json:"pax_count"`
}

type UpdateFlightRequest struct {
	PNR           string  `json:"pnr"`
	DepartureTime string  `json:"departure_time"`
	Price         float64 `json:"price"`
	Deposit       float64 `json:"deposit"`
	PaxCount      int     `json:"pax_count"`
}
