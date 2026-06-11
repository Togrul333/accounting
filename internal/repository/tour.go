package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"accounting/internal/model"
)

type TourRepository interface {
	GetAll(ctx context.Context) ([]model.Tour, error)
	GetByID(ctx context.Context, id int64) (*model.Tour, error)
	Create(ctx context.Context, req model.CreateTourRequest) (*model.Tour, error)
	Update(ctx context.Context, id int64, req model.UpdateTourRequest) (*model.Tour, error)
	Delete(ctx context.Context, id int64) error
}

type tourRepo struct {
	db *gorm.DB
}

func NewTourRepository(db *gorm.DB) TourRepository {
	return &tourRepo{db: db}
}

const tourSelectQuery = `
	SELECT t.id, t.code, t.start_date, t.end_date,
	       t.tour_category_id, tc.name AS tour_category_name, tc.price AS tour_category_price,
	       t.room_id, r.price AS room_price, r.beds_count AS room_beds_count,
	       t.created_at, t.updated_at
	FROM tours t
	JOIN tour_categories tc ON tc.id = t.tour_category_id
	JOIN rooms r ON r.id = t.room_id`

func (r *tourRepo) GetAll(ctx context.Context) ([]model.Tour, error) {
	var tours []model.Tour
	err := r.db.WithContext(ctx).Raw(tourSelectQuery + ` ORDER BY t.id`).Scan(&tours).Error
	if tours == nil {
		tours = []model.Tour{}
	}
	return tours, err
}

func (r *tourRepo) GetByID(ctx context.Context, id int64) (*model.Tour, error) {
	var t model.Tour
	result := r.db.WithContext(ctx).Raw(tourSelectQuery+` WHERE t.id = ?`, id).Scan(&t)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &t, nil
}

func (r *tourRepo) Create(ctx context.Context, req model.CreateTourRequest) (*model.Tour, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, err
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, err
	}
	t := model.Tour{
		Code:           req.Code,
		StartDate:      startDate,
		EndDate:        endDate,
		TourCategoryID: req.TourCategoryID,
		RoomID:         req.RoomID,
	}
	if err := r.db.WithContext(ctx).Create(&t).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, t.ID)
}

func (r *tourRepo) Update(ctx context.Context, id int64, req model.UpdateTourRequest) (*model.Tour, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, err
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, err
	}
	result := r.db.WithContext(ctx).Model(&model.Tour{}).Where("id = ?", id).Updates(map[string]any{
		"code":             req.Code,
		"start_date":       startDate,
		"end_date":         endDate,
		"tour_category_id": req.TourCategoryID,
		"room_id":          req.RoomID,
	})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return r.GetByID(ctx, id)
}

func (r *tourRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Tour{}, id).Error
}
