package service

import (
	"context"

	"accounting/internal/model"
	"accounting/internal/repository"
)

type IncomeService struct {
	repo repository.IncomeRepository
}

func NewIncomeService(repo repository.IncomeRepository) *IncomeService {
	return &IncomeService{repo: repo}
}

func (s *IncomeService) GetAll(ctx context.Context) ([]model.Income, error) {
	return s.repo.GetAll(ctx)
}

func (s *IncomeService) GetByID(ctx context.Context, id int64) (*model.Income, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *IncomeService) Create(ctx context.Context, req model.CreateIncomeRequest) (*model.Income, error) {
	return s.repo.Create(ctx, req)
}

func (s *IncomeService) Update(ctx context.Context, id int64, req model.UpdateIncomeRequest) (*model.Income, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *IncomeService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
