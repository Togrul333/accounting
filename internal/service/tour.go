package service

import (
	"context"
	"errors"
	"log"

	"accounting/internal/model"
	"accounting/internal/repository"
)

var (
	ErrFlightRequired = errors.New("en az bir uçuş seçimi zorunludur")
	ErrRoomRequired   = errors.New("en az bir oda seçimi zorunludur")
)

type TourService struct {
	repo            repository.TourRepository
	defaultTaskRepo repository.DefaultTaskRepository
	taskRepo        repository.TaskRepository
}

func NewTourService(
	repo repository.TourRepository,
	defaultTaskRepo repository.DefaultTaskRepository,
	taskRepo repository.TaskRepository,
) *TourService {
	return &TourService{repo: repo, defaultTaskRepo: defaultTaskRepo, taskRepo: taskRepo}
}

func (s *TourService) GetAll(ctx context.Context) ([]model.Tour, error) {
	return s.repo.GetAll(ctx)
}

func (s *TourService) GetByID(ctx context.Context, id int64) (*model.Tour, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TourService) Create(ctx context.Context, req model.CreateTourRequest) (*model.Tour, error) {
	if len(req.Rooms) == 0 {
		return nil, ErrRoomRequired
	}
	if len(req.FlightIDs) == 0 {
		return nil, ErrFlightRequired
	}
	tour, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	s.createDefaultTasks(ctx, tour)
	return tour, nil
}

// createDefaultTasks — создаёт задачи из шаблонов и привязывает их к туру.
// Ошибка здесь не отменяет создание тура: тур уже создан, задачи можно добавить вручную.
func (s *TourService) createDefaultTasks(ctx context.Context, tour *model.Tour) {
	templates, err := s.defaultTaskRepo.GetAll(ctx)
	if err != nil {
		log.Printf("varsayılan görevler okunamadı (tur %d): %v", tour.ID, err)
		return
	}
	for _, tpl := range templates {
		dueDate := tour.StartDate.AddDate(0, 0, -tpl.DaysBeforeStart)
		_, err := s.taskRepo.Create(ctx, model.CreateTaskRequest{
			Title:       tpl.Title,
			Description: tpl.Description,
			Status:      "todo",
			DueDate:     dueDate.Format("2006-01-02"),
			HocaUserID:  tpl.HocaUserID,
			TourID:      &tour.ID,
		})
		if err != nil {
			log.Printf("varsayılan görev oluşturulamadı (tur %d, şablon %d): %v", tour.ID, tpl.ID, err)
		}
	}
}

func (s *TourService) Update(ctx context.Context, id int64, req model.UpdateTourRequest) (*model.Tour, error) {
	if len(req.Rooms) == 0 {
		return nil, ErrRoomRequired
	}
	if len(req.FlightIDs) == 0 {
		return nil, ErrFlightRequired
	}
	return s.repo.Update(ctx, id, req)
}

func (s *TourService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
