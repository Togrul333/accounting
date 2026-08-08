package repository

import (
	"context"

	"gorm.io/gorm"

	"accounting/internal/model"
)

type DefaultTaskRepository interface {
	GetAll(ctx context.Context) ([]model.DefaultTask, error)
	GetByID(ctx context.Context, id int64) (*model.DefaultTask, error)
	Create(ctx context.Context, req model.CreateDefaultTaskRequest) (*model.DefaultTask, error)
	Update(ctx context.Context, id int64, req model.UpdateDefaultTaskRequest) (*model.DefaultTask, error)
	Delete(ctx context.Context, id int64) error
}

type defaultTaskRepo struct {
	db *gorm.DB
}

func NewDefaultTaskRepository(db *gorm.DB) DefaultTaskRepository {
	return &defaultTaskRepo{db: db}
}

// defaultTaskBaseQuery — шаблоны вместе с именем исполнителя по умолчанию.
const defaultTaskBaseQuery = `
	SELECT d.id, d.title, d.description, d.days_before_start, d.hoca_user_id,
	       TRIM(CONCAT(h.first_name, ' ', h.last_name)) AS hoca_user_name,
	       d.created_at, d.updated_at
	FROM default_tasks d
	LEFT JOIN hoca_users h ON h.id = d.hoca_user_id`

func (r *defaultTaskRepo) GetAll(ctx context.Context) ([]model.DefaultTask, error) {
	var tasks []model.DefaultTask
	err := r.db.WithContext(ctx).Raw(defaultTaskBaseQuery + ` ORDER BY d.days_before_start DESC, d.id`).Scan(&tasks).Error
	if tasks == nil {
		tasks = []model.DefaultTask{}
	}
	return tasks, err
}

func (r *defaultTaskRepo) GetByID(ctx context.Context, id int64) (*model.DefaultTask, error) {
	var t model.DefaultTask
	result := r.db.WithContext(ctx).Raw(defaultTaskBaseQuery+` WHERE d.id = ?`, id).Scan(&t)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &t, nil
}

func (r *defaultTaskRepo) Create(ctx context.Context, req model.CreateDefaultTaskRequest) (*model.DefaultTask, error) {
	t := model.DefaultTask{
		Title:           req.Title,
		Description:     req.Description,
		DaysBeforeStart: req.DaysBeforeStart,
		HocaUserID:      normalizeNullableID(req.HocaUserID),
	}
	if err := r.db.WithContext(ctx).Create(&t).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, t.ID)
}

func (r *defaultTaskRepo) Update(ctx context.Context, id int64, req model.UpdateDefaultTaskRequest) (*model.DefaultTask, error) {
	result := r.db.WithContext(ctx).Model(&model.DefaultTask{}).Where("id = ?", id).Updates(map[string]any{
		"title":             req.Title,
		"description":       req.Description,
		"days_before_start": req.DaysBeforeStart,
		"hoca_user_id":      normalizeNullableID(req.HocaUserID),
	})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		// Значения могли не измениться — проверяем существование записи.
		if _, err := r.GetByID(ctx, id); err != nil {
			return nil, err
		}
	}
	return r.GetByID(ctx, id)
}

func (r *defaultTaskRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.DefaultTask{}, id).Error
}
