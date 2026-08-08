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

type tourFlight struct {
	TourID   int64 `gorm:"column:tour_id"`
	FlightID int64 `gorm:"column:flight_id"`
}

func (tourFlight) TableName() string { return "tour_flights" }

type tourRoom struct {
	TourID int64   `gorm:"column:tour_id"`
	RoomID int64   `gorm:"column:room_id"`
	Price  float64 `gorm:"column:price"`
}

func (tourRoom) TableName() string { return "tour_rooms" }

type tourRoomRow struct {
	TourID    int64   `gorm:"column:tour_id"`
	RoomID    int64   `gorm:"column:room_id"`
	Code      string  `gorm:"column:code"`
	BedsCount int     `gorm:"column:beds_count"`
	Price     float64 `gorm:"column:price"`
}

type tourFlightRow struct {
	TourID        int64     `gorm:"column:tour_id"`
	ID            int64     `gorm:"column:id"`
	PNR           string    `gorm:"column:pnr"`
	DepartureTime time.Time `gorm:"column:departure_time"`
	Price         float64   `gorm:"column:price"`
	Deposit       float64   `gorm:"column:deposit"`
	PaxCount      int       `gorm:"column:pax_count"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

const tourSelectQuery = `
	SELECT t.id, t.code, t.start_date, t.end_date,
	       t.tour_category_id, tc.name AS tour_category_name, tc.price AS tour_category_price,
	       t.created_at, t.updated_at
	FROM tours t
	JOIN tour_categories tc ON tc.id = t.tour_category_id`

const tourRoomsSelectQuery = `
	SELECT tr.tour_id, tr.room_id, r.code, r.beds_count, tr.price
	FROM tour_rooms tr
	JOIN rooms r ON r.id = tr.room_id
	WHERE tr.tour_id IN ?
	ORDER BY r.beds_count, r.code`

const tourFlightsSelectQuery = `
	SELECT tf.tour_id, f.id, f.pnr, f.departure_time, f.price, f.deposit, f.pax_count, f.created_at, f.updated_at
	FROM tour_flights tf
	JOIN flights f ON f.id = tf.flight_id
	WHERE tf.tour_id IN ?
	ORDER BY f.departure_time`

func (r *tourRepo) loadFlights(ctx context.Context, tourIDs []int64) (map[int64][]model.Flight, error) {
	result := make(map[int64][]model.Flight)
	if len(tourIDs) == 0 {
		return result, nil
	}
	var rows []tourFlightRow
	if err := r.db.WithContext(ctx).Raw(tourFlightsSelectQuery, tourIDs).Scan(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		result[row.TourID] = append(result[row.TourID], model.Flight{
			ID:            row.ID,
			PNR:           row.PNR,
			DepartureTime: row.DepartureTime,
			Price:         row.Price,
			Deposit:       row.Deposit,
			PaxCount:      row.PaxCount,
			CreatedAt:     row.CreatedAt,
			UpdatedAt:     row.UpdatedAt,
		})
	}
	return result, nil
}

func (r *tourRepo) loadRooms(ctx context.Context, tourIDs []int64) (map[int64][]model.TourRoom, error) {
	result := make(map[int64][]model.TourRoom)
	if len(tourIDs) == 0 {
		return result, nil
	}
	var rows []tourRoomRow
	if err := r.db.WithContext(ctx).Raw(tourRoomsSelectQuery, tourIDs).Scan(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		result[row.TourID] = append(result[row.TourID], model.TourRoom{
			RoomID:    row.RoomID,
			Code:      row.Code,
			BedsCount: row.BedsCount,
			Price:     row.Price,
		})
	}
	return result, nil
}

func (r *tourRepo) GetAll(ctx context.Context) ([]model.Tour, error) {
	var tours []model.Tour
	if err := r.db.WithContext(ctx).Raw(tourSelectQuery + ` ORDER BY t.id`).Scan(&tours).Error; err != nil {
		return nil, err
	}
	if tours == nil {
		tours = []model.Tour{}
	}
	ids := make([]int64, len(tours))
	for i, t := range tours {
		ids[i] = t.ID
	}
	flightsByTour, err := r.loadFlights(ctx, ids)
	if err != nil {
		return nil, err
	}
	roomsByTour, err := r.loadRooms(ctx, ids)
	if err != nil {
		return nil, err
	}
	for i := range tours {
		tours[i].Flights = flightsByTour[tours[i].ID]
		if tours[i].Flights == nil {
			tours[i].Flights = []model.Flight{}
		}
		tours[i].Rooms = roomsByTour[tours[i].ID]
		if tours[i].Rooms == nil {
			tours[i].Rooms = []model.TourRoom{}
		}
	}
	return tours, nil
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
	flightsByTour, err := r.loadFlights(ctx, []int64{id})
	if err != nil {
		return nil, err
	}
	t.Flights = flightsByTour[id]
	if t.Flights == nil {
		t.Flights = []model.Flight{}
	}
	roomsByTour, err := r.loadRooms(ctx, []int64{id})
	if err != nil {
		return nil, err
	}
	t.Rooms = roomsByTour[id]
	if t.Rooms == nil {
		t.Rooms = []model.TourRoom{}
	}
	return &t, nil
}

func (r *tourRepo) setFlights(tx *gorm.DB, tourID int64, flightIDs []int64) error {
	if err := tx.Where("tour_id = ?", tourID).Delete(&tourFlight{}).Error; err != nil {
		return err
	}
	if len(flightIDs) == 0 {
		return nil
	}
	links := make([]tourFlight, len(flightIDs))
	for i, fid := range flightIDs {
		links[i] = tourFlight{TourID: tourID, FlightID: fid}
	}
	return tx.Create(&links).Error
}

func (r *tourRepo) setRooms(tx *gorm.DB, tourID int64, rooms []model.TourRoomInput) error {
	if err := tx.Where("tour_id = ?", tourID).Delete(&tourRoom{}).Error; err != nil {
		return err
	}
	if len(rooms) == 0 {
		return nil
	}
	// Одна и та же комната может прийти дважды (например, повторяющаяся строка в листе Turlar).
	// Без схлопывания это конфликт по составному первичному ключу — берём последнюю цену.
	links := make([]tourRoom, 0, len(rooms))
	seen := make(map[int64]int, len(rooms))
	for _, rm := range rooms {
		if i, ok := seen[rm.RoomID]; ok {
			links[i].Price = rm.Price
			continue
		}
		seen[rm.RoomID] = len(links)
		links = append(links, tourRoom{TourID: tourID, RoomID: rm.RoomID, Price: rm.Price})
	}
	return tx.Create(&links).Error
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
	}
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&t).Error; err != nil {
			return err
		}
		if err := r.setRooms(tx, t.ID, req.Rooms); err != nil {
			return err
		}
		return r.setFlights(tx, t.ID, req.FlightIDs)
	})
	if err != nil {
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
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Существование проверяем отдельным запросом: RowsAffected здесь не показатель —
		// MySQL вернёт 0, если поля тура не изменились, а правили только комнаты или рейсы.
		var count int64
		if err := tx.Model(&model.Tour{}).Where("id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return gorm.ErrRecordNotFound
		}
		err := tx.Model(&model.Tour{}).Where("id = ?", id).Updates(map[string]any{
			"code":             req.Code,
			"start_date":       startDate,
			"end_date":         endDate,
			"tour_category_id": req.TourCategoryID,
		}).Error
		if err != nil {
			return err
		}
		if err := r.setRooms(tx, id, req.Rooms); err != nil {
			return err
		}
		return r.setFlights(tx, id, req.FlightIDs)
	})
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *tourRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Tour{}, id).Error
}
