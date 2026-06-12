package model

import "time"

type User struct {
	ID        int64     `gorm:"primaryKey"`
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateProfileRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdatePasswordRequest struct {
	Current string `json:"current"`
	New     string `json:"new"`
}
