package model

import "time"

type Task struct {
	ID           int64      `json:"id" gorm:"primaryKey"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Status       string     `json:"status"`
	DueDate      *time.Time `json:"due_date"`
	HocaUserID   *int64     `json:"hoca_user_id"`
	HocaUserName string     `json:"hoca_user_name,omitempty" gorm:"<-:false"`
	TourID       *int64     `json:"tour_id"`
	TourCode     string     `json:"tour_code,omitempty" gorm:"<-:false"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	DueDate     string `json:"due_date"`
	HocaUserID  *int64 `json:"hoca_user_id"`
	TourID      *int64 `json:"tour_id"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	DueDate     string `json:"due_date"`
	HocaUserID  *int64 `json:"hoca_user_id"`
	TourID      *int64 `json:"tour_id"`
}

type UpdateTaskStatusRequest struct {
	Status string `json:"status"`
}
