package service

import (
	"context"
	"errors"

	"accounting/internal/model"
	"accounting/internal/repository"
)

var ErrFlightRequired = errors.New("en az bir uçuş seçimi zorunludur")

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
	if len(req.FlightIDs) == 0 {
		return nil, ErrFlightRequired
	}
	return s.repo.Create(ctx, req)
}

func (s *TourService) Update(ctx context.Context, id int64, req model.UpdateTourRequest) (*model.Tour, error) {
	if len(req.FlightIDs) == 0 {
		return nil, ErrFlightRequired
	}
	return s.repo.Update(ctx, id, req)
}

func (s *TourService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
