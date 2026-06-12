package model

type ExchangeRates struct {
	USD float64 `json:"usd"`
	EUR float64 `json:"eur"`
	GBP float64 `json:"gbp"`
}

type UpdateRatesRequest struct {
	USD float64 `json:"usd"`
	EUR float64 `json:"eur"`
	GBP float64 `json:"gbp"`
}
