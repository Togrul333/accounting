package service

import (
	"context"

	"accounting/internal/model"
	"accounting/internal/repository"
)

type TourService struct {
	repo repository.TourRepository
}

func NewTourService(repo repository.TourRepository) *TourService {
	return &TourService{repo: repo}
}

func (s *TourService) GetAll(ctx context.Context) ([]model.Tour, error) {
	return s.repo.GetAll(ctx)
}

func (s *TourService) GetByID(ctx context.Context, id int64) (*model.Tour, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TourService) Create(ctx context.Context, req model.CreateTourRequest) (*model.Tour, error) {
	return s.repo.Create(ctx, req)
}

func (s *TourService) Update(ctx context.Context, id int64, req model.UpdateTourRequest) (*model.Tour, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *TourService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
