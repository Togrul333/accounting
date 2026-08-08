package model

import "time"

// DefaultTask — шаблон задачи: при создании тура из каждого шаблона
// создаётся задача, привязанная к этому туру.
type DefaultTask struct {
	ID              int64     `json:"id" gorm:"primaryKey"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	DaysBeforeStart int       `json:"days_before_start"`
	HocaUserID      *int64    `json:"hoca_user_id"`
	HocaUserName    string    `json:"hoca_user_name,omitempty" gorm:"<-:false"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (DefaultTask) TableName() string {
	return "default_tasks"
}

type CreateDefaultTaskRequest struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	DaysBeforeStart int    `json:"days_before_start"`
	HocaUserID      *int64 `json:"hoca_user_id"`
}

type UpdateDefaultTaskRequest struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	DaysBeforeStart int    `json:"days_before_start"`
	HocaUserID      *int64 `json:"hoca_user_id"`
}
