package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"accounting/internal/model"
)

type SheetLinkRepository interface {
	GetAll(ctx context.Context) ([]model.SheetLink, error)
	Upsert(ctx context.Context, url, spreadsheetID string) (*model.SheetLink, error)
}

type sheetLinkRepo struct {
	db *gorm.DB
}

func NewSheetLinkRepository(db *gorm.DB) SheetLinkRepository {
	return &sheetLinkRepo{db: db}
}

func (r *sheetLinkRepo) GetAll(ctx context.Context) ([]model.SheetLink, error) {
	var links []model.SheetLink
	err := r.db.WithContext(ctx).Order("updated_at DESC").Find(&links).Error
	if links == nil {
		links = []model.SheetLink{}
	}
	return links, err
}

// Upsert eyni spreadsheet_id üçün linki yaradır ya da mövcud olanı yeniləyir (son istifadə tarixi üstə çıxsın deyə).
func (r *sheetLinkRepo) Upsert(ctx context.Context, url, spreadsheetID string) (*model.SheetLink, error) {
	link := model.SheetLink{URL: url, SpreadsheetID: spreadsheetID}
	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "spreadsheet_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"url"}),
	}).Create(&link).Error
	if err != nil {
		return nil, err
	}

	var saved model.SheetLink
	if err := r.db.WithContext(ctx).Where("spreadsheet_id = ?", spreadsheetID).First(&saved).Error; err != nil {
		return nil, err
	}
	return &saved, nil
}
