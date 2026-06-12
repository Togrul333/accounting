package repository

import (
	"context"
	"strconv"

	"gorm.io/gorm"

	"accounting/internal/model"
)

type SettingRepository interface {
	GetRates(ctx context.Context) (model.ExchangeRates, error)
	UpdateRates(ctx context.Context, req model.UpdateRatesRequest) error
}

type settingRepo struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) SettingRepository {
	return &settingRepo{db: db}
}

func (r *settingRepo) GetRates(ctx context.Context) (model.ExchangeRates, error) {
	type row struct {
		Key   string `gorm:"column:key"`
		Value string `gorm:"column:value"`
	}
	var rows []row
	if err := r.db.WithContext(ctx).Raw("SELECT `key`, `value` FROM settings WHERE `key` IN ('rate_usd','rate_eur','rate_gbp')").Scan(&rows).Error; err != nil {
		return model.ExchangeRates{}, err
	}
	var rates model.ExchangeRates
	for _, rv := range rows {
		v, _ := strconv.ParseFloat(rv.Value, 64)
		switch rv.Key {
		case "rate_usd":
			rates.USD = v
		case "rate_eur":
			rates.EUR = v
		case "rate_gbp":
			rates.GBP = v
		}
	}
	return rates, nil
}

func (r *settingRepo) UpdateRates(ctx context.Context, req model.UpdateRatesRequest) error {
	updates := map[string]float64{
		"rate_usd": req.USD,
		"rate_eur": req.EUR,
		"rate_gbp": req.GBP,
	}
	for key, val := range updates {
		if err := r.db.WithContext(ctx).Exec(
			"INSERT INTO settings (`key`, `value`) VALUES (?, ?) ON DUPLICATE KEY UPDATE `value` = ?",
			key, strconv.FormatFloat(val, 'f', 4, 64), strconv.FormatFloat(val, 'f', 4, 64),
		).Error; err != nil {
			return err
		}
	}
	return nil
}
