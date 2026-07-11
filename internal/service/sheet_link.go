package service

import (
	"context"

	"accounting/internal/model"
	"accounting/internal/repository"
)

type SheetLinkService struct {
	repo repository.SheetLinkRepository
}

func NewSheetLinkService(repo repository.SheetLinkRepository) *SheetLinkService {
	return &SheetLinkService{repo: repo}
}

func (s *SheetLinkService) GetAll(ctx context.Context) ([]model.SheetLink, error) {
	return s.repo.GetAll(ctx)
}

func (s *SheetLinkService) Upsert(ctx context.Context, url, spreadsheetID string) (*model.SheetLink, error) {
	return s.repo.Upsert(ctx, url, spreadsheetID)
}
