package model

// ExchangeRates — 1 birim yabancı para = kaç AZN (manat).
type ExchangeRates struct {
	USD float64 `json:"usd"`
	EUR float64 `json:"eur"`
	GBP float64 `json:"gbp"`
	TRY float64 `json:"try"`
}

// RateFor bir para biriminin manat (AZN) karşılığı kurunu döner. AZN için her zaman 1,
// Ayarlar sayfasında girilmemiş (0) veya desteklenmeyen (örn. RUB) para birimleri için
// ok=false — bu durumda o para birimi manat toplamına dahil edilmemeli.
func (r ExchangeRates) RateFor(currency string) (rate float64, ok bool) {
	switch currency {
	case "AZN":
		return 1, true
	case "USD":
		return r.USD, r.USD > 0
	case "EUR":
		return r.EUR, r.EUR > 0
	case "GBP":
		return r.GBP, r.GBP > 0
	case "TRY":
		return r.TRY, r.TRY > 0
	default:
		return 0, false
	}
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
