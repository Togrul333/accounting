// Package meta — клиент Meta Marketing API (Graph API) для чтения расходов
// по рекламным кампаниям Facebook / Instagram.
//
// Используется только чтение (права ads_read + read_insights), поэтому
// App Review со стороны Meta не требуется — достаточно system user токена
// внутри того же Business Manager, которому принадлежит рекламный кабинет.
package meta

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	graphBaseURL      = "https://graph.facebook.com"
	defaultAPIVersion = "v23.0"
	maxPages          = 100      // страховка от бесконечной пагинации
	maxBodyBytes      = 32 << 20 // 32 MB — потолок на размер ответа
)

// Client — тонкая обёртка над Graph API.
type Client struct {
	token   string
	version string
	http    *http.Client
}

// NewClient создаёт клиент с указанным access token.
func NewClient(token string) *Client {
	return &Client{
		token:   token,
		version: APIVersion(),
		http:    &http.Client{Timeout: 60 * time.Second},
	}
}

// APIVersion — версия Graph API (META_API_VERSION, по умолчанию v23.0).
func APIVersion() string {
	if v := os.Getenv("META_API_VERSION"); v != "" {
		return v
	}
	return defaultAPIVersion
}

// EnvToken — системный токен из окружения (META_ACCESS_TOKEN).
func EnvToken() string {
	return os.Getenv("META_ACCESS_TOKEN")
}

// EnvAdAccountID — рекламный кабинет из окружения (META_AD_ACCOUNT_ID).
func EnvAdAccountID() string {
	return NormalizeAdAccountID(os.Getenv("META_AD_ACCOUNT_ID"))
}

// NormalizeAdAccountID приводит id кабинета к виду act_1234567890.
func NormalizeAdAccountID(id string) string {
	id = strings.TrimSpace(id)
	if id == "" || strings.HasPrefix(id, "act_") {
		return id
	}
	return "act_" + id
}

// AccountInfo — базовые данные рекламного кабинета.
type AccountInfo struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Currency      string `json:"currency"`
	AccountStatus int    `json:"account_status"`
	AmountSpent   string `json:"amount_spent"`
}

// Insight — статистика одной кампании за один день.
type Insight struct {
	CampaignID   string  `json:"campaign_id"`
	CampaignName string  `json:"campaign_name"`
	Date         string  `json:"date"` // YYYY-MM-DD
	Spend        float64 `json:"spend"`
	Impressions  int64   `json:"impressions"`
	Clicks       int64   `json:"clicks"`
	Reach        int64   `json:"reach"`
	CPC          float64 `json:"cpc"`
	CTR          float64 `json:"ctr"`
	Currency     string  `json:"currency"`
}

// graphError — формат ошибки Graph API.
type graphError struct {
	Message   string `json:"message"`
	Type      string `json:"type"`
	Code      int    `json:"code"`
	Subcode   int    `json:"error_subcode"`
	UserTitle string `json:"error_user_title"`
	UserMsg   string `json:"error_user_msg"`
}

func (e *graphError) Error() string {
	if e.UserMsg != "" {
		return fmt.Sprintf("Meta API: %s (code %d)", e.UserMsg, e.Code)
	}
	return fmt.Sprintf("Meta API: %s (code %d)", e.Message, e.Code)
}

// AccountInfo проверяет токен и возвращает данные кабинета.
// Удобно для кнопки «Bağlantıyı test et» в UI.
func (c *Client) AccountInfo(ctx context.Context, adAccountID string) (*AccountInfo, error) {
	if c.token == "" {
		return nil, fmt.Errorf("Meta access token boşdur")
	}
	adAccountID = NormalizeAdAccountID(adAccountID)
	if adAccountID == "" {
		return nil, fmt.Errorf("reklam hesabı id-si boşdur")
	}

	q := url.Values{}
	q.Set("fields", "name,currency,account_status,amount_spent")
	q.Set("access_token", c.token)
	endpoint := fmt.Sprintf("%s/%s/%s?%s", graphBaseURL, c.version, adAccountID, q.Encode())

	var info AccountInfo
	if err := c.getJSON(ctx, endpoint, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// insightsResponse — сырой ответ /insights: все числа приходят строками.
type insightsResponse struct {
	Data []struct {
		CampaignID      string `json:"campaign_id"`
		CampaignName    string `json:"campaign_name"`
		DateStart       string `json:"date_start"`
		Spend           string `json:"spend"`
		Impressions     string `json:"impressions"`
		Clicks          string `json:"clicks"`
		Reach           string `json:"reach"`
		CPC             string `json:"cpc"`
		CTR             string `json:"ctr"`
		AccountCurrency string `json:"account_currency"`
	} `json:"data"`
	Paging struct {
		Next string `json:"next"`
	} `json:"paging"`
}

// CampaignInsights возвращает расходы по кампаниям с разбивкой по дням
// за период [since, until] включительно (даты в формате YYYY-MM-DD).
func (c *Client) CampaignInsights(ctx context.Context, adAccountID, since, until string) ([]Insight, error) {
	if c.token == "" {
		return nil, fmt.Errorf("Meta access token boşdur")
	}
	adAccountID = NormalizeAdAccountID(adAccountID)
	if adAccountID == "" {
		return nil, fmt.Errorf("reklam hesabı id-si boşdur")
	}

	timeRange, err := json.Marshal(map[string]string{"since": since, "until": until})
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Set("level", "campaign")
	q.Set("time_increment", "1") // разбивка по дням
	q.Set("time_range", string(timeRange))
	q.Set("fields", "campaign_id,campaign_name,spend,impressions,clicks,reach,cpc,ctr,account_currency")
	q.Set("limit", "500")
	q.Set("access_token", c.token)

	endpoint := fmt.Sprintf("%s/%s/%s/insights?%s", graphBaseURL, c.version, adAccountID, q.Encode())

	result := make([]Insight, 0, 64)
	for page := 0; endpoint != "" && page < maxPages; page++ {
		var resp insightsResponse
		if err := c.getJSON(ctx, endpoint, &resp); err != nil {
			return nil, err
		}
		for _, row := range resp.Data {
			result = append(result, Insight{
				CampaignID:   row.CampaignID,
				CampaignName: row.CampaignName,
				Date:         row.DateStart,
				Spend:        parseFloat(row.Spend),
				Impressions:  parseInt(row.Impressions),
				Clicks:       parseInt(row.Clicks),
				Reach:        parseInt(row.Reach),
				CPC:          parseFloat(row.CPC),
				CTR:          parseFloat(row.CTR),
				Currency:     row.AccountCurrency,
			})
		}
		endpoint = resp.Paging.Next // next уже содержит access_token
	}

	return result, nil
}

// getJSON выполняет GET и разбирает ответ, отдельно вытаскивая ошибку Graph API.
func (c *Client) getJSON(ctx context.Context, endpoint string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("Meta API-yə qoşulmaq alınmadı: %w", err)
	}
	defer resp.Body.Close()

	// Тело читаем целиком: и успех, и ошибка приходят как JSON.
	var envelope struct {
		Error *graphError `json:"error"`
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes))
	if err != nil {
		return fmt.Errorf("Meta API cavabı oxuna bilmədi: %w", err)
	}
	if err := json.Unmarshal(body, &envelope); err == nil && envelope.Error != nil {
		return envelope.Error
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("Meta API xətası (HTTP %d)", resp.StatusCode)
	}
	return json.Unmarshal(body, out)
}

func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return v
}

func parseInt(s string) int64 {
	v, _ := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	return v
}
