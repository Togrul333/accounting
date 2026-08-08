package service

import (
	"context"
	"errors"
	"strings"

	"accounting/internal/model"
	"accounting/internal/repository"
)

var ErrEmptyComment = errors.New("yorum boş olamaz")

type TaskCommentService struct {
	repo repository.TaskCommentRepository
}

func NewTaskCommentService(repo repository.TaskCommentRepository) *TaskCommentService {
	return &TaskCommentService{repo: repo}
}

func (s *TaskCommentService) GetByTaskID(ctx context.Context, taskID int64) ([]model.TaskComment, error) {
	return s.repo.GetByTaskID(ctx, taskID)
}

func (s *TaskCommentService) Create(ctx context.Context, taskID int64, userID *int64, body string) (*model.TaskComment, error) {
	body = strings.TrimSpace(body)
	if body == "" {
		return nil, ErrEmptyComment
	}
	return s.repo.Create(ctx, taskID, userID, body)
}

func (s *TaskCommentService) Update(ctx context.Context, id int64, body string) (*model.TaskComment, error) {
	body = strings.TrimSpace(body)
	if body == "" {
		return nil, ErrEmptyComment
	}
	return s.repo.Update(ctx, id, body)
}

func (s *TaskCommentService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
