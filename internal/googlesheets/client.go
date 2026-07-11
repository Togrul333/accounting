package googlesheets

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var idPattern = regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)

// Client оборачивает sheets.Service для чтения строк из Google Sheets.
type Client struct {
	svc *sheets.Service
}

// Tab — одна вкладка (лист) внутри Google Sheets документа.
type Tab struct {
	GID   int64  `json:"gid"`
	Title string `json:"title"`
}

// NewClient создаёт клиент, используя JSON-ключ сервисного аккаунта по пути credentialsPath.
func NewClient(ctx context.Context, credentialsPath string) (*Client, error) {
	svc, err := sheets.NewService(ctx, option.WithCredentialsFile(credentialsPath), option.WithScopes(sheets.SpreadsheetsReadonlyScope))
	if err != nil {
		return nil, fmt.Errorf("sheets servisi yaradıla bilmədi: %w", err)
	}
	return &Client{svc: svc}, nil
}

// CredentialsPath переменной окружения GOOGLE_CREDENTIALS_PATH (по умолчанию ./credentials.json).
func CredentialsPath() string {
	if p := os.Getenv("GOOGLE_CREDENTIALS_PATH"); p != "" {
		return p
	}
	return "./credentials.json"
}

// ParseSpreadsheetURL извлекает SPREADSHEET_ID и gid (id листа, если указан в ссылке).
func ParseSpreadsheetURL(rawURL string) (spreadsheetID string, gid int64, hasGID bool, err error) {
	m := idPattern.FindStringSubmatch(rawURL)
	if m == nil {
		return "", 0, false, fmt.Errorf("keçərsiz Google Sheets linki")
	}
	spreadsheetID = m[1]

	u, err := url.Parse(rawURL)
	if err == nil {
		if g := u.Query().Get("gid"); g != "" {
			gid, _ = strconv.ParseInt(g, 10, 64)
			hasGID = true
		} else if u.Fragment != "" {
			if frag, ferr := url.ParseQuery(u.Fragment); ferr == nil {
				if g := frag.Get("gid"); g != "" {
					gid, _ = strconv.ParseInt(g, 10, 64)
					hasGID = true
				}
			}
		}
	}
	return spreadsheetID, gid, hasGID, nil
}

// ListTabs возвращает все вкладки (листы) документа.
func (c *Client) ListTabs(spreadsheetID string) ([]Tab, error) {
	ss, err := c.svc.Spreadsheets.Get(spreadsheetID).Fields("sheets.properties").Do()
	if err != nil {
		return nil, fmt.Errorf("cədvələ giriş alınmadı: %w", err)
	}
	tabs := make([]Tab, 0, len(ss.Sheets))
	for _, sh := range ss.Sheets {
		tabs = append(tabs, Tab{GID: sh.Properties.SheetId, Title: sh.Properties.Title})
	}
	return tabs, nil
}

// FetchRowsByGID читает все строки конкретной вкладки по её gid.
func (c *Client) FetchRowsByGID(spreadsheetID string, gid int64) (sheetTitle string, rows [][]string, err error) {
	tabs, err := c.ListTabs(spreadsheetID)
	if err != nil {
		return "", nil, err
	}
	title := ""
	for _, t := range tabs {
		if t.GID == gid {
			title = t.Title
			break
		}
	}
	if title == "" {
		if len(tabs) == 0 {
			return "", nil, fmt.Errorf("cədvəldə heç bir vərəq tapılmadı")
		}
		title = tabs[0].Title
	}

	resp, err := c.svc.Spreadsheets.Values.Get(spreadsheetID, "'"+title+"'").Do()
	if err != nil {
		return "", nil, fmt.Errorf("sətrlər oxuna bilmədi: %w", err)
	}

	rows = make([][]string, 0, len(resp.Values))
	for _, row := range resp.Values {
		strRow := make([]string, len(row))
		for i, cell := range row {
			strRow[i] = fmt.Sprintf("%v", cell)
		}
		rows = append(rows, strRow)
	}
	return title, rows, nil
}
