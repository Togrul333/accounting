package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"accounting/internal/model"
)

type TaskRepository interface {
	GetAll(ctx context.Context) ([]model.Task, error)
	GetByID(ctx context.Context, id int64) (*model.Task, error)
	Create(ctx context.Context, req model.CreateTaskRequest) (*model.Task, error)
	Update(ctx context.Context, id int64, req model.UpdateTaskRequest) (*model.Task, error)
	UpdateStatus(ctx context.Context, id int64, status string) (*model.Task, error)
	Delete(ctx context.Context, id int64) error
}

type taskRepo struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepo{db: db}
}

func parseTaskDueDate(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	d, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *taskRepo) GetAll(ctx context.Context) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.WithContext(ctx).Order("id DESC").Find(&tasks).Error
	if tasks == nil {
		tasks = []model.Task{}
	}
	return tasks, err
}

func (r *taskRepo) GetByID(ctx context.Context, id int64) (*model.Task, error) {
	var t model.Task
	if err := r.db.WithContext(ctx).First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *taskRepo) Create(ctx context.Context, req model.CreateTaskRequest) (*model.Task, error) {
	dueDate, err := parseTaskDueDate(req.DueDate)
	if err != nil {
		return nil, err
	}
	status := req.Status
	if status == "" {
		status = "todo"
	}
	t := model.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		DueDate:     dueDate,
	}
	if err := r.db.WithContext(ctx).Create(&t).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, t.ID)
}

func (r *taskRepo) Update(ctx context.Context, id int64, req model.UpdateTaskRequest) (*model.Task, error) {
	dueDate, err := parseTaskDueDate(req.DueDate)
	if err != nil {
		return nil, err
	}
	result := r.db.WithContext(ctx).Model(&model.Task{}).Where("id = ?", id).Updates(map[string]any{
		"title":       req.Title,
		"description": req.Description,
		"status":      req.Status,
		"due_date":    dueDate,
	})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return r.GetByID(ctx, id)
}

func (r *taskRepo) UpdateStatus(ctx context.Context, id int64, status string) (*model.Task, error) {
	result := r.db.WithContext(ctx).Model(&model.Task{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return r.GetByID(ctx, id)
}

func (r *taskRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Task{}, id).Error
}
