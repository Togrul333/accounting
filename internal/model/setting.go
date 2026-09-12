package model

// ExchangeRates — 1 birim yabancı para = kaç AZN (manat).
type ExchangeRates struct {
	USD float64 `json:"usd"`
	EUR float64 `json:"eur"`
	GBP float64 `json:"gbp"`
	TRY float64 `json:"try"`
}

type UpdateRatesRequest struct {
	USD float64 `json:"usd"`
	EUR float64 `json:"eur"`
	GBP float64 `json:"gbp"`
	TRY float64 `json:"try"`
}

// Ключи настроек Telegram в таблице settings.
const (
	SettingTelegramBotToken    = "telegram_bot_token"
	SettingTelegramBotUsername = "telegram_bot_username"
	SettingTelegramEnabled     = "telegram_enabled"
	SettingTelegramOffset      = "telegram_update_offset"
)

// TelegramSettings — состояние интеграции. Токен наружу не отдаётся,
// вместо него флаг HasToken.
type TelegramSettings struct {
	Enabled     bool   `json:"enabled"`
	HasToken    bool   `json:"has_token"`
	BotUsername string `json:"bot_username"`
}

type UpdateTelegramSettingsRequest struct {
	BotToken string `json:"bot_token"`
	Enabled  bool   `json:"enabled"`
}
