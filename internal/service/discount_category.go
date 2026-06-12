package service

import (
	"context"

	"accounting/internal/model"
	"accounting/internal/repository"
)

type DiscountCategoryService struct {
	repo repository.DiscountCategoryRepository
}

func NewDiscountCategoryService(repo repository.DiscountCategoryRepository) *DiscountCategoryService {
	return &DiscountCategoryService{repo: repo}
}

func (s *DiscountCategoryService) GetAll(ctx context.Context) ([]model.DiscountCategory, error) {
	return s.repo.GetAll(ctx)
}

func (s *DiscountCategoryService) GetByID(ctx context.Context, id int64) (*model.DiscountCategory, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *DiscountCategoryService) Create(ctx context.Context, req model.CreateDiscountCategoryRequest) (*model.DiscountCategory, error) {
	return s.repo.Create(ctx, req)
}

func (s *DiscountCategoryService) Update(ctx context.Context, id int64, req model.UpdateDiscountCategoryRequest) (*model.DiscountCategory, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *DiscountCategoryService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
