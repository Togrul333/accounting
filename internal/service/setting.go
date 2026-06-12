package service

import (
	"context"

	"accounting/internal/model"
	"accounting/internal/repository"
)

type SettingService struct {
	repo repository.SettingRepository
}

func NewSettingService(repo repository.SettingRepository) *SettingService {
	return &SettingService{repo: repo}
}

func (s *SettingService) GetRates(ctx context.Context) (model.ExchangeRates, error) {
	return s.repo.GetRates(ctx)
}

func (s *SettingService) UpdateRates(ctx context.Context, req model.UpdateRatesRequest) error {
	return s.repo.UpdateRates(ctx, req)
}
