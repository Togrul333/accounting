package repository

import (
	"context"

	"gorm.io/gorm"

	"accounting/internal/model"
)

type TaskCommentRepository interface {
	GetByTaskID(ctx context.Context, taskID int64) ([]model.TaskComment, error)
	GetByID(ctx context.Context, id int64) (*model.TaskComment, error)
	Create(ctx context.Context, taskID int64, userID *int64, body string) (*model.TaskComment, error)
	Update(ctx context.Context, id int64, body string) (*model.TaskComment, error)
	Delete(ctx context.Context, id int64) error
}

type taskCommentRepo struct {
	db *gorm.DB
}

func NewTaskCommentRepository(db *gorm.DB) TaskCommentRepository {
	return &taskCommentRepo{db: db}
}

// taskCommentBaseQuery — комментарий вместе с именем автора.
const taskCommentBaseQuery = `
	SELECT tc.id, tc.task_id, tc.user_id, u.name AS user_name,
	       tc.body, tc.created_at, tc.updated_at
	FROM task_comments tc
	LEFT JOIN users u ON u.id = tc.user_id`

func (r *taskCommentRepo) GetByTaskID(ctx context.Context, taskID int64) ([]model.TaskComment, error) {
	var comments []model.TaskComment
	err := r.db.WithContext(ctx).
		Raw(taskCommentBaseQuery+` WHERE tc.task_id = ? ORDER BY tc.created_at ASC, tc.id ASC`, taskID).
		Scan(&comments).Error
	if comments == nil {
		comments = []model.TaskComment{}
	}
	return comments, err
}

func (r *taskCommentRepo) GetByID(ctx context.Context, id int64) (*model.TaskComment, error) {
	var c model.TaskComment
	result := r.db.WithContext(ctx).Raw(taskCommentBaseQuery+` WHERE tc.id = ?`, id).Scan(&c)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &c, nil
}

func (r *taskCommentRepo) Create(ctx context.Context, taskID int64, userID *int64, body string) (*model.TaskComment, error) {
	c := model.TaskComment{
		TaskID: taskID,
		UserID: normalizeNullableID(userID),
		Body:   body,
	}
	if err := r.db.WithContext(ctx).Create(&c).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, c.ID)
}

func (r *taskCommentRepo) Update(ctx context.Context, id int64, body string) (*model.TaskComment, error) {
	result := r.db.WithContext(ctx).Model(&model.TaskComment{}).Where("id = ?", id).Update("body", body)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		// Текст мог не измениться — проверяем существование записи.
		if _, err := r.GetByID(ctx, id); err != nil {
			return nil, err
		}
	}
	return r.GetByID(ctx, id)
}

func (r *taskCommentRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.TaskComment{}, id).Error
}
