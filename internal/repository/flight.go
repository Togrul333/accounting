package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"accounting/internal/model"
)

type FlightRepository interface {
	GetAll(ctx context.Context) ([]model.Flight, error)
	GetByID(ctx context.Context, id int64) (*model.Flight, error)
	Create(ctx context.Context, req model.CreateFlightRequest) (*model.Flight, error)
	Update(ctx context.Context, id int64, req model.UpdateFlightRequest) (*model.Flight, error)
	Delete(ctx context.Context, id int64) error
}

type flightRepo struct {
	db *gorm.DB
}

func NewFlightRepository(db *gorm.DB) FlightRepository {
	return &flightRepo{db: db}
}

const flightDepartureTimeLayout = "2006-01-02T15:04"

func (r *flightRepo) GetAll(ctx context.Context) ([]model.Flight, error) {
	var flights []model.Flight
	err := r.db.WithContext(ctx).Order("departure_time desc").Find(&flights).Error
	return flights, err
}

func (r *flightRepo) GetByID(ctx context.Context, id int64) (*model.Flight, error) {
	var f model.Flight
	err := r.db.WithContext(ctx).First(&f, id).Error
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *flightRepo) Create(ctx context.Context, req model.CreateFlightRequest) (*model.Flight, error) {
	departureTime, err := time.Parse(flightDepartureTimeLayout, req.DepartureTime)
	if err != nil {
		return nil, err
	}
	f := model.Flight{
		PNR:           req.PNR,
		DepartureTime: departureTime,
		Price:         req.Price,
		Deposit:       req.Deposit,
		PaxCount:      req.PaxCount,
	}
	if err := r.db.WithContext(ctx).Create(&f).Error; err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *flightRepo) Update(ctx context.Context, id int64, req model.UpdateFlightRequest) (*model.Flight, error) {
	departureTime, err := time.Parse(flightDepartureTimeLayout, req.DepartureTime)
	if err != nil {
		return nil, err
	}
	result := r.db.WithContext(ctx).Model(&model.Flight{}).Where("id = ?", id).Updates(map[string]any{
		"pnr":            req.PNR,
		"departure_time": departureTime,
		"price":          req.Price,
		"deposit":        req.Deposit,
		"pax_count":      req.PaxCount,
	})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return r.GetByID(ctx, id)
}

func (r *flightRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Flight{}, id).Error
}
