package service

import (
	"context"

	"accounting/internal/model"
	"accounting/internal/repository"
)

type FlightService struct {
	repo repository.FlightRepository
}

func NewFlightService(repo repository.FlightRepository) *FlightService {
	return &FlightService{repo: repo}
}

func (s *FlightService) GetAll(ctx context.Context) ([]model.Flight, error) {
	return s.repo.GetAll(ctx)
}

func (s *FlightService) GetByID(ctx context.Context, id int64) (*model.Flight, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *FlightService) Create(ctx context.Context, req model.CreateFlightRequest) (*model.Flight, error) {
	return s.repo.Create(ctx, req)
}

func (s *FlightService) Update(ctx context.Context, id int64, req model.UpdateFlightRequest) (*model.Flight, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *FlightService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
