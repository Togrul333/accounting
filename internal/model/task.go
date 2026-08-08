package model

import "time"

type Task struct {
	ID           int64      `json:"id" gorm:"primaryKey"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Status       string     `json:"status"`
	Priority     string     `json:"priority"`
	DueDate      *time.Time `json:"due_date"`
	HocaUserID   *int64     `json:"hoca_user_id"`
	HocaUserName string     `json:"hoca_user_name,omitempty" gorm:"<-:false"`
	TourID       *int64     `json:"tour_id"`
	TourCode     string     `json:"tour_code,omitempty" gorm:"<-:false"`
	OrderID      *int64     `json:"order_id"`
	OrderLabel   string     `json:"order_label,omitempty" gorm:"<-:false"`
	ClientID     *int64     `json:"client_id"`
	ClientName   string     `json:"client_name,omitempty" gorm:"<-:false"`
	CommentCount int        `json:"comment_count" gorm:"<-:false"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	DueDate     string `json:"due_date"`
	HocaUserID  *int64 `json:"hoca_user_id"`
	TourID      *int64 `json:"tour_id"`
	OrderID     *int64 `json:"order_id"`
	ClientID    *int64 `json:"client_id"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	DueDate     string `json:"due_date"`
	HocaUserID  *int64 `json:"hoca_user_id"`
	TourID      *int64 `json:"tour_id"`
	OrderID     *int64 `json:"order_id"`
	ClientID    *int64 `json:"client_id"`
}

type UpdateTaskStatusRequest struct {
	Status string `json:"status"`
}

// TaskComment — комментарий в ленте обсуждения задачи.
type TaskComment struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	TaskID    int64     `json:"task_id"`
	UserID    *int64    `json:"user_id"`
	UserName  string    `json:"user_name,omitempty" gorm:"<-:false"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTaskCommentRequest struct {
	Body string `json:"body"`
}

type UpdateTaskCommentRequest struct {
	Body string `json:"body"`
}
