package model

import "time"

type SheetLink struct {
	ID            int64     `json:"id" gorm:"primaryKey"`
	URL           string    `json:"url"`
	SpreadsheetID string    `json:"spreadsheet_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
