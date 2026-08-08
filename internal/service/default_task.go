package service

import (
	"context"

	"accounting/internal/model"
	"accounting/internal/repository"
)

type DefaultTaskService struct {
	repo repository.DefaultTaskRepository
}

func NewDefaultTaskService(repo repository.DefaultTaskRepository) *DefaultTaskService {
	return &DefaultTaskService{repo: repo}
}

func (s *DefaultTaskService) GetAll(ctx context.Context) ([]model.DefaultTask, error) {
	return s.repo.GetAll(ctx)
}

func (s *DefaultTaskService) GetByID(ctx context.Context, id int64) (*model.DefaultTask, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *DefaultTaskService) Create(ctx context.Context, req model.CreateDefaultTaskRequest) (*model.DefaultTask, error) {
	return s.repo.Create(ctx, req)
}

func (s *DefaultTaskService) Update(ctx context.Context, id int64, req model.UpdateDefaultTaskRequest) (*model.DefaultTask, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *DefaultTaskService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
