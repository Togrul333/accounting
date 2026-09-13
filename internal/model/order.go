package model

import (
	"sort"
	"time"
)

// baseCurrency — indirimlerin ve tur/oda fiyatlarının her zaman kayıtlı olduğu
// para birimi (Discount ve Room modellerinde ayrı bir currency alanı yok).
const baseCurrency = "AZN"

// CurrencyAmount tek bir para birimindeki toplamı temsil eder —
// gelirler farklı hesaplardan (farklı para birimlerinden) geldiği için kullanılır.
type CurrencyAmount struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

// CurrencyAmountAZN bir CurrencyAmount'ı, Ayarlar → Döviz Kurları'ndaki kura göre
// manat (AZN) karşılığıyla birlikte taşır — ekranda "= X ₼ (kurs Y)" bilgisi için.
// Bu sadece bilgi amaçlıdır; asıl tutar hâlâ kendi para biriminde (Amount/Currency) tutulur.
type CurrencyAmountAZN struct {
	Currency string
	Amount   float64
	Rate     float64
	AZN      float64
	HasRate  bool
}

// BreakdownWithAZN her satırı manat karşılığıyla birlikte döner ve toplam manat
// karşılığını hesaplar. Kuru Ayarlar'da girilmemiş para birimleri (Rate=0) toplama
// dahil edilmez — sessizce yanlış toplam göstermektense o satır AZN'siz kalır.
func BreakdownWithAZN(amounts []CurrencyAmount, rates ExchangeRates) (items []CurrencyAmountAZN, totalAZN float64) {
	items = make([]CurrencyAmountAZN, 0, len(amounts))
	for _, a := range amounts {
		item := CurrencyAmountAZN{Currency: a.Currency, Amount: a.Amount}
		if rate, ok := rates.RateFor(a.Currency); ok {
			item.Rate = rate
			item.AZN = a.Amount * rate
			item.HasRate = true
			totalAZN += item.AZN
		}
		items = append(items, item)
	}
	return items, totalAZN
}

type Order struct {
	ID                int64      `json:"id" gorm:"primaryKey"`
	ClientID          int64      `json:"client_id"`
	ClientName        string     `json:"client_name,omitempty" gorm:"<-:false"`
	TourID            int64      `json:"tour_id"`
	TourCode          string     `json:"tour_code,omitempty" gorm:"<-:false"`
	TourCategoryName  string     `json:"tour_category_name,omitempty" gorm:"<-:false"`
	RoomID            *int64     `json:"room_id"`
	RoomCode          string     `json:"room_code,omitempty" gorm:"<-:false"`
	RoomPrice         float64    `json:"room_price,omitempty" gorm:"<-:false"`
	TourPrice         float64    `json:"tour_price" gorm:"-"`
	IncomeCount       int        `json:"income_count" gorm:"<-:false"`
	IncomeTotal       float64    `json:"income_total" gorm:"<-:false"`
	Incomes           []Income   `json:"incomes,omitempty" gorm:"-"`
	DiscountCount     int        `json:"discount_count" gorm:"<-:false"`
	DiscountTotal     float64    `json:"discount_total" gorm:"<-:false"`
	Discounts         []Discount `json:"discounts,omitempty" gorm:"-"`
	NetTotal          float64    `json:"net_total" gorm:"-"`
	DiscountedPrice   float64    `json:"discounted_price" gorm:"-"`
	RemainingTotal    float64    `json:"remaining_total" gorm:"-"`
	// IncomeByCurrency / NetByCurrency — gelirler birden fazla para biriminden
	// (hesaptan) gelebildiği için IncomeTotal/NetTotal'ı para birimine göre ayırır.
	IncomeByCurrency []CurrencyAmount `json:"income_by_currency,omitempty" gorm:"-"`
	NetByCurrency    []CurrencyAmount `json:"net_by_currency,omitempty" gorm:"-"`
	// IncomeByCurrencyAZN / NetByCurrencyAZN — yalnızca bilgi amaçlı: Ayarlar'daki
	// döviz kuruna göre manat karşılığı (ApplyExchangeRates çağrılmadıkça boştur).
	IncomeByCurrencyAZN []CurrencyAmountAZN `json:"-" gorm:"-"`
	IncomeTotalAZN      float64             `json:"-" gorm:"-"`
	NetByCurrencyAZN    []CurrencyAmountAZN `json:"-" gorm:"-"`
	NetTotalAZN         float64             `json:"-" gorm:"-"`
	CreatedAt           time.Time           `json:"created_at"`
	UpdatedAt           time.Time           `json:"updated_at"`
}

func (o *Order) ComputeNet() {
	o.TourPrice = o.RoomPrice
	o.NetTotal = o.IncomeTotal - o.DiscountTotal
	o.DiscountedPrice = o.TourPrice - o.DiscountTotal
	o.RemainingTotal = o.DiscountedPrice - o.IncomeTotal
}

// ApplyCurrencyBreakdown gelirleri hesabın para birimine göre ayırır ve
// IncomeByCurrency/NetByCurrency/RemainingTotal alanlarını buna göre günceller.
// İndirim her zaman AZN cinsinden girildiği için (Discount modelinde para birimi
// alanı yok) yalnızca AZN satırından düşülür; diğer para birimlerindeki gelirler
// AZN'ye çevrilmeden, oldukları gibi ayrı satırlarda gösterilir. Bu, ComputeNet
// içindeki eski davranışın (tüm para birimlerinin tek sayıda toplanması) yerini alır.
func (o *Order) ApplyCurrencyBreakdown(incomeByCurrency map[string]float64) {
	currencySet := make(map[string]bool, len(incomeByCurrency)+1)
	for cur := range incomeByCurrency {
		currencySet[cur] = true
	}
	if o.DiscountTotal > 0 {
		currencySet[baseCurrency] = true
	}
	currencies := make([]string, 0, len(currencySet))
	for cur := range currencySet {
		currencies = append(currencies, cur)
	}
	sort.Strings(currencies)

	o.IncomeByCurrency = make([]CurrencyAmount, 0, len(incomeByCurrency))
	o.NetByCurrency = make([]CurrencyAmount, 0, len(currencies))
	for _, cur := range currencies {
		amt := incomeByCurrency[cur]
		if amt != 0 {
			o.IncomeByCurrency = append(o.IncomeByCurrency, CurrencyAmount{Currency: cur, Amount: amt})
		}
		net := amt
		if cur == baseCurrency {
			net -= o.DiscountTotal
		}
		o.NetByCurrency = append(o.NetByCurrency, CurrencyAmount{Currency: cur, Amount: net})
	}
	o.RemainingTotal = o.DiscountedPrice - incomeByCurrency[baseCurrency]
}

// ApplyExchangeRates IncomeByCurrency/NetByCurrency listelerinin manat karşılığını
// hesaplar (yalnızca bilgi amaçlı gösterim için — bkz. CurrencyAmountAZN).
func (o *Order) ApplyExchangeRates(rates ExchangeRates) {
	o.IncomeByCurrencyAZN, o.IncomeTotalAZN = BreakdownWithAZN(o.IncomeByCurrency, rates)
	o.NetByCurrencyAZN, o.NetTotalAZN = BreakdownWithAZN(o.NetByCurrency, rates)
}

type CreateOrderRequest struct {
	ClientID int64                 `json:"client_id"`
	TourID   int64                 `json:"tour_id"`
	RoomID   *int64                `json:"room_id"`
	Incomes  []CreateIncomeRequest `json:"incomes"`
}

type UpdateOrderRequest struct {
	ClientID int64  `json:"client_id"`
	TourID   int64  `json:"tour_id"`
	RoomID   *int64 `json:"room_id"`
}
