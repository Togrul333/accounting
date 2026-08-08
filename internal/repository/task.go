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

// taskBaseQuery — задачи вместе с подписями всех связанных сущностей:
// исполнитель (hoca_users), тур, заказ, клиент и число комментариев.
const taskBaseQuery = `
	SELECT t.id, t.title, t.description, t.status, t.priority, t.due_date, t.hoca_user_id,
	       TRIM(CONCAT(h.first_name, ' ', h.last_name)) AS hoca_user_name,
	       t.tour_id, tr.code AS tour_code,
	       t.order_id,
	       CASE WHEN t.order_id IS NULL THEN NULL
	            ELSE TRIM(CONCAT('#', o.id, ' ', COALESCE(oc.first_name, ''), ' ', COALESCE(oc.last_name, '')))
	       END AS order_label,
	       t.client_id, TRIM(CONCAT(c.first_name, ' ', c.last_name)) AS client_name,
	       (SELECT COUNT(*) FROM task_comments tcm WHERE tcm.task_id = t.id) AS comment_count,
	       t.created_at, t.updated_at
	FROM tasks t
	LEFT JOIN hoca_users h ON h.id = t.hoca_user_id
	LEFT JOIN tours tr     ON tr.id = t.tour_id
	LEFT JOIN orders o     ON o.id = t.order_id
	LEFT JOIN clients oc   ON oc.id = o.client_id
	LEFT JOIN clients c    ON c.id = t.client_id`

// taskOrder — сначала по приоритету, внутри приоритета по сроку (без срока — в конец).
const taskOrder = `
	ORDER BY FIELD(t.priority, 'urgent', 'high', 'medium', 'low'),
	         t.due_date IS NULL, t.due_date ASC, t.id DESC`

// normalizeTaskPriority — пустое значение из формы означает обычный приоритет.
func normalizeTaskPriority(p string) string {
	switch p {
	case "urgent", "high", "medium", "low":
		return p
	default:
		return "medium"
	}
}

// normalizeNullableID — 0 или nil из формы означают «не выбрано».
func normalizeNullableID(id *int64) *int64 {
	if id == nil || *id <= 0 {
		return nil
	}
	return id
}

func (r *taskRepo) GetAll(ctx context.Context) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.WithContext(ctx).Raw(taskBaseQuery + taskOrder).Scan(&tasks).Error
	if tasks == nil {
		tasks = []model.Task{}
	}
	return tasks, err
}

func (r *taskRepo) GetByID(ctx context.Context, id int64) (*model.Task, error) {
	var t model.Task
	result := r.db.WithContext(ctx).Raw(taskBaseQuery+` WHERE t.id = ?`, id).Scan(&t)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
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
		Priority:    normalizeTaskPriority(req.Priority),
		DueDate:     dueDate,
		HocaUserID:  normalizeNullableID(req.HocaUserID),
		TourID:      normalizeNullableID(req.TourID),
		OrderID:     normalizeNullableID(req.OrderID),
		ClientID:    normalizeNullableID(req.ClientID),
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
		"title":        req.Title,
		"description":  req.Description,
		"status":       req.Status,
		"priority":     normalizeTaskPriority(req.Priority),
		"due_date":     dueDate,
		"hoca_user_id": normalizeNullableID(req.HocaUserID),
		"tour_id":      normalizeNullableID(req.TourID),
		"order_id":     normalizeNullableID(req.OrderID),
		"client_id":    normalizeNullableID(req.ClientID),
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

func (r *taskRepo) UpdateStatus(ctx context.Context, id int64, status string) (*model.Task, error) {
	result := r.db.WithContext(ctx).Model(&model.Task{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		// Статус мог не измениться (перенос в ту же колонку) — это не 404.
		if _, err := r.GetByID(ctx, id); err != nil {
			return nil, err
		}
	}
	return r.GetByID(ctx, id)
}

func (r *taskRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Task{}, id).Error
}
