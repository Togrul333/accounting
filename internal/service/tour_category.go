package service

import (
	"context"

	"accounting/internal/model"
	"accounting/internal/repository"
)

type TourCategoryService struct {
	repo repository.TourCategoryRepository
}

func NewTourCategoryService(repo repository.TourCategoryRepository) *TourCategoryService {
	return &TourCategoryService{repo: repo}
}

func (s *TourCategoryService) GetAll(ctx context.Context) ([]model.TourCategory, error) {
	return s.repo.GetAll(ctx)
}

func (s *TourCategoryService) GetByID(ctx context.Context, id int64) (*model.TourCategory, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TourCategoryService) Create(ctx context.Context, req model.CreateTourCategoryRequest) (*model.TourCategory, error) {
	return s.repo.Create(ctx, req)
}

func (s *TourCategoryService) Update(ctx context.Context, id int64, req model.UpdateTourCategoryRequest) (*model.TourCategory, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *TourCategoryService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
