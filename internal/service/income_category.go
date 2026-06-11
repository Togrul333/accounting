package service

import (
	"context"

	"accounting/internal/model"
	"accounting/internal/repository"
)

type IncomeCategoryService struct {
	repo repository.IncomeCategoryRepository
}

func NewIncomeCategoryService(repo repository.IncomeCategoryRepository) *IncomeCategoryService {
	return &IncomeCategoryService{repo: repo}
}

func (s *IncomeCategoryService) GetAll(ctx context.Context) ([]model.IncomeCategory, error) {
	return s.repo.GetAll(ctx)
}

func (s *IncomeCategoryService) GetByID(ctx context.Context, id int64) (*model.IncomeCategory, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *IncomeCategoryService) Create(ctx context.Context, req model.CreateIncomeCategoryRequest) (*model.IncomeCategory, error) {
	return s.repo.Create(ctx, req)
}

func (s *IncomeCategoryService) Update(ctx context.Context, id int64, req model.UpdateIncomeCategoryRequest) (*model.IncomeCategory, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *IncomeCategoryService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
