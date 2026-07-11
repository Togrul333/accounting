package repository

import (
	"context"

	"gorm.io/gorm"

	"accounting/internal/model"
)

type RoomRepository interface {
	GetAll(ctx context.Context) ([]model.Room, error)
	GetByID(ctx context.Context, id int64) (*model.Room, error)
	Create(ctx context.Context, req model.CreateRoomRequest) (*model.Room, error)
	Update(ctx context.Context, id int64, req model.UpdateRoomRequest) (*model.Room, error)
	Delete(ctx context.Context, id int64) error
}

type roomRepo struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepo{db: db}
}

func (r *roomRepo) GetAll(ctx context.Context) ([]model.Room, error) {
	var rooms []model.Room
	err := r.db.WithContext(ctx).Order("id").Find(&rooms).Error
	return rooms, err
}

func (r *roomRepo) GetByID(ctx context.Context, id int64) (*model.Room, error) {
	var rm model.Room
	err := r.db.WithContext(ctx).First(&rm, id).Error
	if err != nil {
		return nil, err
	}
	return &rm, nil
}

func (r *roomRepo) Create(ctx context.Context, req model.CreateRoomRequest) (*model.Room, error) {
	rm := model.Room{Code: req.Code, Price: req.Price, BedsCount: req.BedsCount}
	if err := r.db.WithContext(ctx).Create(&rm).Error; err != nil {
		return nil, err
	}
	return &rm, nil
}

func (r *roomRepo) Update(ctx context.Context, id int64, req model.UpdateRoomRequest) (*model.Room, error) {
	result := r.db.WithContext(ctx).Model(&model.Room{}).Where("id = ?", id).Updates(map[string]any{
		"code":       req.Code,
		"price":      req.Price,
		"beds_count": req.BedsCount,
	})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return r.GetByID(ctx, id)
}

func (r *roomRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Room{}, id).Error
}
