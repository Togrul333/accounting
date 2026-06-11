package model

import "time"

type TourCategory struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTourCategoryRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type UpdateTourCategoryRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type Room struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Price     float64   `json:"price"`
	BedsCount int       `json:"beds_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateRoomRequest struct {
	Price     float64 `json:"price"`
	BedsCount int     `json:"beds_count"`
}

type UpdateRoomRequest struct {
	Price     float64 `json:"price"`
	BedsCount int     `json:"beds_count"`
}

type Tour struct {
	ID                int64     `json:"id" gorm:"primaryKey"`
	Code              string    `json:"code"`
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`
	TourCategoryID    int64     `json:"tour_category_id"`
	TourCategoryName  string    `json:"tour_category_name,omitempty" gorm:"<-:false"`
	TourCategoryPrice float64   `json:"tour_category_price,omitempty" gorm:"<-:false"`
	RoomID            int64     `json:"room_id"`
	RoomPrice         float64   `json:"room_price,omitempty" gorm:"<-:false"`
	RoomBedsCount     int       `json:"room_beds_count,omitempty" gorm:"<-:false"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CreateTourRequest struct {
	Code            string `json:"code"`
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
	TourCategoryID  int64  `json:"tour_category_id"`
	RoomID          int64  `json:"room_id"`
}

type UpdateTourRequest struct {
	Code            string `json:"code"`
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
	TourCategoryID  int64  `json:"tour_category_id"`
	RoomID          int64  `json:"room_id"`
}
