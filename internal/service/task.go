package service

import (
	"context"

	"accounting/internal/model"
	"accounting/internal/repository"
)

type TaskService struct {
	repo     repository.TaskRepository
	telegram *TelegramService
}

func NewTaskService(repo repository.TaskRepository, telegram *TelegramService) *TaskService {
	return &TaskService{repo: repo, telegram: telegram}
}

func (s *TaskService) GetAll(ctx context.Context) ([]model.Task, error) {
	return s.repo.GetAll(ctx)
}

func (s *TaskService) GetByID(ctx context.Context, id int64) (*model.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TaskService) Create(ctx context.Context, req model.CreateTaskRequest) (*model.Task, error) {
	task, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	s.notifyAssignee(task)
	return task, nil
}

// Update — при смене исполнителя новый получает уведомление в Telegram.
func (s *TaskService) Update(ctx context.Context, id int64, req model.UpdateTaskRequest) (*model.Task, error) {
	before, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	task, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}
	if assigneeChanged(before, task) {
		s.notifyAssignee(task)
	}
	return task, nil
}

func (s *TaskService) UpdateStatus(ctx context.Context, id int64, status string) (*model.Task, error) {
	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *TaskService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// notifyAssignee — уведомление уходит в фоне: медленный или недоступный
// Telegram не должен задерживать ответ и не отменяет сохранение задачи.
func (s *TaskService) notifyAssignee(task *model.Task) {
	if s.telegram == nil || task == nil || task.HocaUserID == nil {
		return
	}
	go s.telegram.NotifyTaskAssigned(*task)
}

func assigneeChanged(before, after *model.Task) bool {
	if after.HocaUserID == nil {
		return false
	}
	if before.HocaUserID == nil {
		return true
	}
	return *before.HocaUserID != *after.HocaUserID
}
