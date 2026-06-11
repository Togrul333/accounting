package service

import (
	"context"

	"accounting/internal/model"
	"accounting/internal/repository"
)

type RoomService struct {
	repo repository.RoomRepository
}

func NewRoomService(repo repository.RoomRepository) *RoomService {
	return &RoomService{repo: repo}
}

func (s *RoomService) GetAll(ctx context.Context) ([]model.Room, error) {
	return s.repo.GetAll(ctx)
}

func (s *RoomService) GetByID(ctx context.Context, id int64) (*model.Room, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *RoomService) Create(ctx context.Context, req model.CreateRoomRequest) (*model.Room, error) {
	return s.repo.Create(ctx, req)
}

func (s *RoomService) Update(ctx context.Context, id int64, req model.UpdateRoomRequest) (*model.Room, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *RoomService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
