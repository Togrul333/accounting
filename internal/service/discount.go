package service

import (
	"context"

	"accounting/internal/model"
	"accounting/internal/repository"
)

type DiscountService struct {
	repo repository.DiscountRepository
}

func NewDiscountService(repo repository.DiscountRepository) *DiscountService {
	return &DiscountService{repo: repo}
}

func (s *DiscountService) GetAll(ctx context.Context) ([]model.Discount, error) {
	return s.repo.GetAll(ctx)
}

func (s *DiscountService) GetByID(ctx context.Context, id int64) (*model.Discount, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *DiscountService) Create(ctx context.Context, req model.CreateDiscountRequest) (*model.Discount, error) {
	return s.repo.Create(ctx, req)
}

func (s *DiscountService) Update(ctx context.Context, id int64, req model.UpdateDiscountRequest) (*model.Discount, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *DiscountService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
