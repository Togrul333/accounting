package service

import (
	"context"

	"accounting/internal/model"
	"accounting/internal/repository"
)

type ExpenseService struct {
	repo repository.ExpenseRepository
}

func NewExpenseService(repo repository.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: repo}
}

func (s *ExpenseService) GetAll(ctx context.Context) ([]model.Expense, error) {
	return s.repo.GetAll(ctx)
}

func (s *ExpenseService) GetByID(ctx context.Context, id int64) (*model.Expense, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ExpenseService) GetByAccountID(ctx context.Context, accountID int64) ([]model.Expense, error) {
	return s.repo.GetByAccountID(ctx, accountID)
}

func (s *ExpenseService) Create(ctx context.Context, req model.CreateExpenseRequest) (*model.Expense, error) {
	return s.repo.Create(ctx, req)
}

func (s *ExpenseService) BulkCreate(ctx context.Context, reqs []model.CreateExpenseRequest) ([]model.Expense, error) {
	return s.repo.BulkCreate(ctx, reqs)
}

func (s *ExpenseService) Update(ctx context.Context, id int64, req model.UpdateExpenseRequest) (*model.Expense, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *ExpenseService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
