package service

import (
	"context"

	"accounting/internal/model"
	"accounting/internal/repository"
)

type AccountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) GetAll(ctx context.Context) ([]model.Account, error) {
	return s.repo.GetAll(ctx)
}

func (s *AccountService) GetByID(ctx context.Context, id int64) (*model.Account, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *AccountService) Create(ctx context.Context, req model.CreateAccountRequest) (*model.Account, error) {
	return s.repo.Create(ctx, req)
}

func (s *AccountService) Update(ctx context.Context, id int64, req model.UpdateAccountRequest) (*model.Account, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *AccountService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
