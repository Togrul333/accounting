package model

import "time"

type TourCategory struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTourCategoryRequest struct {
	Name string `json:"name"`
}

type UpdateTourCategoryRequest struct {
	Name string `json:"name"`
}

// Room — справочник типов номеров. Цена здесь не хранится: она задаётся
// на конкретном туре, в пивот-таблице tour_rooms.
type Room struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Code      string    `json:"code,omitempty"`
	BedsCount int       `json:"beds_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateRoomRequest struct {
	Code      string `json:"code"`
	BedsCount int    `json:"beds_count"`
}

type UpdateRoomRequest struct {
	Code      string `json:"code"`
	BedsCount int    `json:"beds_count"`
}

// TourRoom — комната, привязанная к туру, вместе с ценой из tour_rooms.
type TourRoom struct {
	RoomID    int64   `json:"room_id"`
	Code      string  `json:"code,omitempty"`
	BedsCount int     `json:"beds_count"`
	Price     float64 `json:"price"`
}

// TourRoomInput — одна строка привязки комнаты в запросах создания/обновления тура.
type TourRoomInput struct {
	RoomID int64   `json:"room_id"`
	Price  float64 `json:"price"`
}

type Tour struct {
	ID                int64      `json:"id" gorm:"primaryKey"`
	Code              string     `json:"code"`
	StartDate         time.Time  `json:"start_date"`
	EndDate           time.Time  `json:"end_date"`
	TourCategoryID    int64      `json:"tour_category_id"`
	TourCategoryName  string     `json:"tour_category_name,omitempty" gorm:"<-:false"`
	Rooms             []TourRoom `json:"rooms" gorm:"-"`
	Flights           []Flight   `json:"flights" gorm:"-"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type CreateTourRequest struct {
	Code           string          `json:"code"`
	StartDate      string          `json:"start_date"`
	EndDate        string          `json:"end_date"`
	TourCategoryID int64           `json:"tour_category_id"`
	Rooms          []TourRoomInput `json:"rooms"`
	FlightIDs      []int64         `json:"flight_ids"`
}

type UpdateTourRequest struct {
	Code           string          `json:"code"`
	StartDate      string          `json:"start_date"`
	EndDate        string          `json:"end_date"`
	TourCategoryID int64           `json:"tour_category_id"`
	Rooms          []TourRoomInput `json:"rooms"`
	FlightIDs      []int64         `json:"flight_ids"`
}
