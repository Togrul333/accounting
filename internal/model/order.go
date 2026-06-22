package model

import "time"

type Order struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	ClientID    int64     `json:"client_id"`
	ClientName  string    `json:"client_name,omitempty" gorm:"<-:false"`
	TourID      int64     `json:"tour_id"`
	TourCode    string    `json:"tour_code,omitempty" gorm:"<-:false"`
	IncomeCount    int        `json:"income_count" gorm:"<-:false"`
	IncomeTotal    float64    `json:"income_total" gorm:"<-:false"`
	Incomes        []Income   `json:"incomes,omitempty" gorm:"-"`
	DiscountCount  int        `json:"discount_count" gorm:"<-:false"`
	DiscountTotal  float64    `json:"discount_total" gorm:"<-:false"`
	Discounts      []Discount `json:"discounts,omitempty" gorm:"-"`
	NetTotal       float64    `json:"net_total" gorm:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (o *Order) ComputeNet() {
	o.NetTotal = o.IncomeTotal - o.DiscountTotal
}

type CreateOrderRequest struct {
	ClientID int64                 `json:"client_id"`
	TourID   int64                 `json:"tour_id"`
	Incomes  []CreateIncomeRequest `json:"incomes"`
}

type UpdateOrderRequest struct {
	ClientID int64 `json:"client_id"`
	TourID   int64 `json:"tour_id"`
}
