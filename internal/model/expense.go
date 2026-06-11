package model

import "time"

type ExpenseCategory struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateExpenseCategoryRequest struct {
	Name string `json:"name"`
}

type UpdateExpenseCategoryRequest struct {
	Name string `json:"name"`
}

type Expense struct {
	ID                  int64     `json:"id"`
	Name                string    `json:"name"`
	Amount              float64   `json:"amount"`
	Date                time.Time `json:"date"`
	ExpenseCategoryID   int64     `json:"expense_category_id"`
	ExpenseCategoryName string    `json:"expense_category_name,omitempty"`
	AccountID           int64     `json:"account_id"`
	AccountName         string    `json:"account_name,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type CreateExpenseRequest struct {
	Name              string  `json:"name"`
	Amount            float64 `json:"amount"`
	Date              string  `json:"date"`
	ExpenseCategoryID int64   `json:"expense_category_id"`
	AccountID         int64   `json:"account_id"`
}

type UpdateExpenseRequest struct {
	Name              string  `json:"name"`
	Amount            float64 `json:"amount"`
	Date              string  `json:"date"`
	ExpenseCategoryID int64   `json:"expense_category_id"`
	AccountID         int64   `json:"account_id"`
}
