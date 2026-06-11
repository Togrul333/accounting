package service

import (
	"context"

	"accounting/internal/model"
	"accounting/internal/repository"
)

type ExpenseCategoryService struct {
	repo repository.ExpenseCategoryRepository
}

func NewExpenseCategoryService(repo repository.ExpenseCategoryRepository) *ExpenseCategoryService {
	return &ExpenseCategoryService{repo: repo}
}

func (s *ExpenseCategoryService) GetAll(ctx context.Context) ([]model.ExpenseCategory, error) {
	return s.repo.GetAll(ctx)
}

func (s *ExpenseCategoryService) GetByID(ctx context.Context, id int64) (*model.ExpenseCategory, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ExpenseCategoryService) Create(ctx context.Context, req model.CreateExpenseCategoryRequest) (*model.ExpenseCategory, error) {
	return s.repo.Create(ctx, req)
}

func (s *ExpenseCategoryService) Update(ctx context.Context, id int64, req model.UpdateExpenseCategoryRequest) (*model.ExpenseCategory, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *ExpenseCategoryService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
