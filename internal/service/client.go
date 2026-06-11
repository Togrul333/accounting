package service

import (
	"context"

	"accounting/internal/model"
	"accounting/internal/repository"
)

type ClientService struct {
	repo repository.ClientRepository
}

func NewClientService(repo repository.ClientRepository) *ClientService {
	return &ClientService{repo: repo}
}

func (s *ClientService) GetAll(ctx context.Context) ([]model.Client, error) {
	return s.repo.GetAll(ctx)
}

func (s *ClientService) GetByID(ctx context.Context, id int64) (*model.Client, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ClientService) Create(ctx context.Context, req model.CreateClientRequest) (*model.Client, error) {
	return s.repo.Create(ctx, req)
}

func (s *ClientService) Update(ctx context.Context, id int64, req model.UpdateClientRequest) (*model.Client, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *ClientService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
