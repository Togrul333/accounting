package service

import (
	"context"

	"accounting/internal/model"
	"accounting/internal/repository"
)

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) GetAll(ctx context.Context) ([]model.Task, error) {
	return s.repo.GetAll(ctx)
}

func (s *TaskService) GetByID(ctx context.Context, id int64) (*model.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TaskService) Create(ctx context.Context, req model.CreateTaskRequest) (*model.Task, error) {
	return s.repo.Create(ctx, req)
}

func (s *TaskService) Update(ctx context.Context, id int64, req model.UpdateTaskRequest) (*model.Task, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *TaskService) UpdateStatus(ctx context.Context, id int64, status string) (*model.Task, error) {
	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *TaskService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
