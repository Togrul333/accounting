package service

import (
	"context"

	"accounting/internal/model"
	"accounting/internal/repository"
)

type ReferansUserService struct {
	repo repository.ReferansUserRepository
}

func NewReferansUserService(repo repository.ReferansUserRepository) *ReferansUserService {
	return &ReferansUserService{repo: repo}
}

func (s *ReferansUserService) GetAll(ctx context.Context) ([]model.ReferansUser, error) {
	return s.repo.GetAll(ctx)
}

func (s *ReferansUserService) GetByID(ctx context.Context, id int64) (*model.ReferansUser, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ReferansUserService) Create(ctx context.Context, req model.CreateReferansUserRequest) (*model.ReferansUser, error) {
	return s.repo.Create(ctx, req)
}

func (s *ReferansUserService) Update(ctx context.Context, id int64, req model.UpdateReferansUserRequest) (*model.ReferansUser, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *ReferansUserService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// Candidates — заказы-кандидаты по совпадению имени/фамилии с reference_name клиента.
func (s *ReferansUserService) Candidates(ctx context.Context, id int64, limit int) ([]model.ReferansOrder, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.repo.Candidates(ctx, *u, limit)
}

func (s *ReferansUserService) SearchOrders(ctx context.Context, id int64, query string, limit int) ([]model.ReferansOrder, error) {
	return s.repo.SearchOrders(ctx, id, query, limit)
}

func (s *ReferansUserService) Referrals(ctx context.Context, id int64) ([]model.ReferansOrder, error) {
	return s.repo.Referrals(ctx, id)
}

func (s *ReferansUserService) AddReferral(ctx context.Context, id, orderID int64) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return err
	}
	return s.repo.AddReferral(ctx, id, orderID)
}

func (s *ReferansUserService) RemoveReferral(ctx context.Context, id, orderID int64) error {
	return s.repo.RemoveReferral(ctx, id, orderID)
}
