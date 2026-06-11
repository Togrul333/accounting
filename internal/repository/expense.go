package repository

import (
	"context"
	"database/sql"

	"accounting/internal/model"
)

type ExpenseRepository interface {
	GetAll(ctx context.Context) ([]model.Expense, error)
	GetByID(ctx context.Context, id int64) (*model.Expense, error)
	GetByAccountID(ctx context.Context, accountID int64) ([]model.Expense, error)
	Create(ctx context.Context, req model.CreateExpenseRequest) (*model.Expense, error)
	Update(ctx context.Context, id int64, req model.UpdateExpenseRequest) (*model.Expense, error)
	Delete(ctx context.Context, id int64) error
}

type expenseRepo struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) ExpenseRepository {
	return &expenseRepo{db: db}
}

func (r *expenseRepo) GetAll(ctx context.Context) ([]model.Expense, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT e.id, e.name, e.amount, e.date, e.expense_category_id, c.name, e.account_id, a.name, e.created_at, e.updated_at
		FROM expenses e
		JOIN expense_categories c ON c.id = e.expense_category_id
		JOIN accounts a ON a.id = e.account_id
		ORDER BY e.date DESC, e.id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []model.Expense
	for rows.Next() {
		var exp model.Expense
		if err := rows.Scan(&exp.ID, &exp.Name, &exp.Amount, &exp.Date, &exp.ExpenseCategoryID, &exp.ExpenseCategoryName, &exp.AccountID, &exp.AccountName, &exp.CreatedAt, &exp.UpdatedAt); err != nil {
			return nil, err
		}
		expenses = append(expenses, exp)
	}
	return expenses, rows.Err()
}

func (r *expenseRepo) GetByID(ctx context.Context, id int64) (*model.Expense, error) {
	var exp model.Expense
	err := r.db.QueryRowContext(ctx, `
		SELECT e.id, e.name, e.amount, e.date, e.expense_category_id, c.name, e.account_id, a.name, e.created_at, e.updated_at
		FROM expenses e
		JOIN expense_categories c ON c.id = e.expense_category_id
		JOIN accounts a ON a.id = e.account_id
		WHERE e.id=?
	`, id).Scan(&exp.ID, &exp.Name, &exp.Amount, &exp.Date, &exp.ExpenseCategoryID, &exp.ExpenseCategoryName, &exp.AccountID, &exp.AccountName, &exp.CreatedAt, &exp.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &exp, nil
}

func (r *expenseRepo) GetByAccountID(ctx context.Context, accountID int64) ([]model.Expense, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT e.id, e.name, e.amount, e.date, e.expense_category_id, c.name, e.account_id, a.name, e.created_at, e.updated_at
		FROM expenses e
		JOIN expense_categories c ON c.id = e.expense_category_id
		JOIN accounts a ON a.id = e.account_id
		WHERE e.account_id = ?
		ORDER BY e.date DESC, e.id DESC
	`, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []model.Expense
	for rows.Next() {
		var exp model.Expense
		if err := rows.Scan(&exp.ID, &exp.Name, &exp.Amount, &exp.Date, &exp.ExpenseCategoryID, &exp.ExpenseCategoryName, &exp.AccountID, &exp.AccountName, &exp.CreatedAt, &exp.UpdatedAt); err != nil {
			return nil, err
		}
		expenses = append(expenses, exp)
	}
	return expenses, rows.Err()
}

func (r *expenseRepo) Create(ctx context.Context, req model.CreateExpenseRequest) (*model.Expense, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO expenses (name, amount, date, expense_category_id, account_id) VALUES (?, ?, ?, ?, ?)
	`, req.Name, req.Amount, req.Date, req.ExpenseCategoryID, req.AccountID)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *expenseRepo) Update(ctx context.Context, id int64, req model.UpdateExpenseRequest) (*model.Expense, error) {
	res, err := r.db.ExecContext(ctx, `
		UPDATE expenses SET name=?, amount=?, date=?, expense_category_id=?, account_id=? WHERE id=?
	`, req.Name, req.Amount, req.Date, req.ExpenseCategoryID, req.AccountID, id)
	if err != nil {
		return nil, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, sql.ErrNoRows
	}
	return r.GetByID(ctx, id)
}

func (r *expenseRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM expenses WHERE id=?`, id)
	return err
}
