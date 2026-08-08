package model

import "time"

// MetaAdAccount — настройка одного рекламного кабинета Meta.
type MetaAdAccount struct {
	ID                int64      `json:"id" gorm:"primaryKey"`
	AdAccountID       string     `json:"ad_account_id"`
	Name              string     `json:"name"`
	Currency          string     `json:"currency"`
	AccessToken       string     `json:"-"` // наружу не отдаём, только флаг has_token
	ExpenseCategoryID *int64     `json:"expense_category_id"`
	AccountID         *int64     `json:"account_id"`
	AutoExpense       bool       `json:"auto_expense"`
	LastSyncedAt      *time.Time `json:"last_synced_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	// Поля из JOIN-ов и вычисляемые для UI (в таблице не хранятся).
	ExpenseCategoryName string `json:"expense_category_name,omitempty" gorm:"<-:false"`
	AccountName         string `json:"account_name,omitempty" gorm:"<-:false"`
	HasToken            bool   `json:"has_token" gorm:"-"`
}

type CreateMetaAdAccountRequest struct {
	AdAccountID       string `json:"ad_account_id"`
	Name              string `json:"name"`
	AccessToken       string `json:"access_token"`
	ExpenseCategoryID *int64 `json:"expense_category_id"`
	AccountID         *int64 `json:"account_id"`
	AutoExpense       bool   `json:"auto_expense"`
}

type UpdateMetaAdAccountRequest struct {
	Name              string `json:"name"`
	AccessToken       string `json:"access_token"` // пусто = не менять существующий токен
	ExpenseCategoryID *int64 `json:"expense_category_id"`
	AccountID         *int64 `json:"account_id"`
	AutoExpense       bool   `json:"auto_expense"`
}

// MetaAdSpend — расход по одной кампании за один день.
type MetaAdSpend struct {
	ID           int64     `json:"id" gorm:"primaryKey"`
	AdAccountID  string    `json:"ad_account_id"`
	CampaignID   string    `json:"campaign_id"`
	CampaignName string    `json:"campaign_name"`
	Date         time.Time `json:"date"`
	Spend        float64   `json:"spend"`
	Impressions  int64     `json:"impressions"`
	Clicks       int64     `json:"clicks"`
	Reach        int64     `json:"reach"`
	CPC          float64   `json:"cpc" gorm:"column:cpc"`
	CTR          float64   `json:"ctr" gorm:"column:ctr"`
	Currency     string    `json:"currency"`
	TourID       *int64    `json:"tour_id"`
	ExpenseID    *int64    `json:"expense_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	TourCode string `json:"tour_code,omitempty" gorm:"-"`
}

// TableName — таблица в единственном числе, GORM сам бы сделал meta_ad_spends.
func (MetaAdSpend) TableName() string { return "meta_ad_spend" }

// MetaCampaignSummary — агрегат по кампании за период (для таблицы в UI).
type MetaCampaignSummary struct {
	CampaignID   string  `json:"campaign_id"`
	CampaignName string  `json:"campaign_name"`
	Spend        float64 `json:"spend"`
	Impressions  int64   `json:"impressions"`
	Clicks       int64   `json:"clicks"`
	Currency     string  `json:"currency"`
	Days         int64   `json:"days"`
	CPC          float64 `json:"cpc" gorm:"-"`
	CTR          float64 `json:"ctr" gorm:"-"`
}

// MetaSyncRequest — параметры ручной синхронизации.
type MetaSyncRequest struct {
	Since string `json:"since"` // YYYY-MM-DD
	Until string `json:"until"` // YYYY-MM-DD
}

// MetaSyncResult — итог синхронизации, показывается пользователю.
type MetaSyncResult struct {
	AdAccountID     string  `json:"ad_account_id"`
	Since           string  `json:"since"`
	Until           string  `json:"until"`
	RowsFetched     int     `json:"rows_fetched"`
	RowsSaved       int     `json:"rows_saved"`
	ExpensesCreated int     `json:"expenses_created"`
	TotalSpend      float64 `json:"total_spend"`
	Currency        string  `json:"currency"`
}
