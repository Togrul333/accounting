package model

import "time"

// HocaUser — сотрудник турагентства; на него можно назначать задачи.
type HocaUser struct {
	ID        int64  `json:"id" gorm:"primaryKey"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	// TelegramChatID заполняется после того, как сотрудник нажал Start в боте.
	TelegramChatID   string    `json:"telegram_chat_id"`
	TelegramUsername string    `json:"telegram_username"`
	TelegramLinkCode string    `json:"telegram_link_code"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// FullName — имя для подписи в интерфейсе и уведомлениях.
func (h HocaUser) FullName() string {
	if h.LastName == "" {
		return h.FirstName
	}
	return h.FirstName + " " + h.LastName
}

func (HocaUser) TableName() string {
	return "hoca_users"
}

type CreateHocaUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

type UpdateHocaUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}
