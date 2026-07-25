package model

import "time"

type Account struct {
	ID            int64     `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name"`
	AccountNumber string    `json:"account_number"`
	Currency      string    `json:"currency"`
	Balance       float64   `json:"balance" gorm:"<-:false"`
	Description   string    `json:"description,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateAccountRequest struct {
	Name          string `json:"name"`
	AccountNumber string `json:"account_number"`
	Currency      string `json:"currency"`
	Description   string `json:"description"`
}

type UpdateAccountRequest struct {
	Name          string `json:"name"`
	AccountNumber string `json:"account_number"`
	Currency      string `json:"currency"`
	Description   string `json:"description"`
}

type StatementRow struct {
	Date            string  `json:"date"`
	Ref             string  `json:"ref"`
	CP              string  `json:"cp"`
	Debit           float64 `json:"debit"`
	Credit          float64 `json:"credit"`
	Desc            string  `json:"desc"`
	Tax             string  `json:"tax"`
	AlreadyImported bool    `json:"already_imported"`
}

// ImportGelirRow banka ekstresindeki tek bir gelir satırı ile birlikte,
// kullanıcının o satır için seçtiği kategori/tur/müşteriyi taşır.
type ImportGelirRow struct {
	StatementRow
	IncomeCategoryID int64  `json:"income_category_id"`
	TourID           *int64 `json:"tour_id"`
	ClientID         *int64 `json:"client_id"`
}

type ImportGelirlerRequest struct {
	Rows []ImportGelirRow `json:"rows"`
}

// ImportGiderRow banka ekstresindeki tek bir gider satırıdır — kategori tüm
// satırlar için tek seferde (ImportGiderlerRequest.ExpenseCategoryID) seçilir.
type ImportGiderRow struct {
	StatementRow
}

type ImportGiderlerRequest struct {
	ExpenseCategoryID int64            `json:"expense_category_id"`
	Rows              []ImportGiderRow `json:"rows"`
}

type StatementPreview struct {
	IBAN        string         `json:"iban"`
	Gelirler    []StatementRow `json:"gelirler"`
	Giderler    []StatementRow `json:"giderler"`
	TotalCredit float64        `json:"total_credit"`
	TotalDebit  float64        `json:"total_debit"`
}
