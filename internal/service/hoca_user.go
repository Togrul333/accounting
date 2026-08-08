package service

import (
	"context"

	"accounting/internal/model"
	"accounting/internal/repository"
)

type HocaUserService struct {
	repo repository.HocaUserRepository
}

func NewHocaUserService(repo repository.HocaUserRepository) *HocaUserService {
	return &HocaUserService{repo: repo}
}

func (s *HocaUserService) GetAll(ctx context.Context) ([]model.HocaUser, error) {
	return s.repo.GetAll(ctx)
}

func (s *HocaUserService) GetByID(ctx context.Context, id int64) (*model.HocaUser, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *HocaUserService) Create(ctx context.Context, req model.CreateHocaUserRequest) (*model.HocaUser, error) {
	return s.repo.Create(ctx, req)
}

func (s *HocaUserService) Update(ctx context.Context, id int64, req model.UpdateHocaUserRequest) (*model.HocaUser, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *HocaUserService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
