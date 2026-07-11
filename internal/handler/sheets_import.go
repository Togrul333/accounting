package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"accounting/internal/googlesheets"
	"accounting/internal/service"
)

type SheetsImportHandler struct {
	sheets    *googlesheets.Client
	linkSvc   *service.SheetLinkService
	tourSvc   *service.TourService
	clientSvc *service.ClientService
	orderSvc  *service.OrderService
}

func NewSheetsImportHandler(sheets *googlesheets.Client, linkSvc *service.SheetLinkService, tourSvc *service.TourService, clientSvc *service.ClientService, orderSvc *service.OrderService) *SheetsImportHandler {
	return &SheetsImportHandler{sheets: sheets, linkSvc: linkSvc, tourSvc: tourSvc, clientSvc: clientSvc, orderSvc: orderSvc}
}

// Links əvvəllər uğurla açılmış Google Sheets linklərini qaytarır (son istifadə tarixinə görə).
func (h *SheetsImportHandler) Links(c *gin.Context) {
	links, err := h.linkSvc.GetAll(c.Request.Context())
	if err != nil {
		log.Printf("sheet links error: %v", err)
		links = nil
	}
	c.JSON(http.StatusOK, gin.H{"links": links})
}

type sheetsTabsRequest struct {
	URL string `json:"url" binding:"required"`
}

// Tabs cədvəldəki bütün vərəqləri (tab) qaytarır ki, istifadəçi hansını çəkəcəyini seçə bilsin.
func (h *SheetsImportHandler) Tabs(c *gin.Context) {
	var req sheetsTabsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "link tələb olunur"})
		return
	}

	if h.sheets == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Google Sheets bağlantısı quraşdırılmayıb"})
		return
	}

	spreadsheetID, gid, hasGID, err := googlesheets.ParseSpreadsheetURL(req.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tabs, err := h.sheets.ListTabs(spreadsheetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := h.linkSvc.Upsert(c.Request.Context(), req.URL, spreadsheetID); err != nil {
		log.Printf("sheet link upsert error: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"spreadsheet_id": spreadsheetID,
		"tabs":           tabs,
		"selected_gid":   gid,
		"has_gid":        hasGID,
	})
}

type sheetsPreviewRequest struct {
	SpreadsheetID string `json:"spreadsheet_id" binding:"required"`
	GID           int64  `json:"gid"`
}

// Preview cədvəldəki bütün sətrləri geri qaytarır (hələ DB-yə yazmadan, sadəcə önizləmə üçün).
func (h *SheetsImportHandler) Preview(c *gin.Context) {
	var req sheetsPreviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cədvəl və vərəq tələb olunur"})
		return
	}

	if h.sheets == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Google Sheets bağlantısı quraşdırılmayıb"})
		return
	}

	sheetTitle, rows, err := h.sheets.FetchRowsByGID(req.SpreadsheetID, req.GID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Vərəqlərin çoxunda sətr 1 başlıq deyil (düymələr, boş sətrlər və s. ola bilər),
	// ona görə burada başlıq təxmin edilmir — sadəcə Google Sheets-dəki kimi A/B/C... sütun adları ilə xam sətrlər qaytarılır.
	maxCols := 0
	for _, row := range rows {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}
	headers := make([]string, maxCols)
	for i := range headers {
		headers[i] = columnLetter(i)
	}

	c.JSON(http.StatusOK, gin.H{
		"sheet_title": sheetTitle,
		"headers":     headers,
		"rows":        rows,
	})
}

// columnLetter Google Sheets-dəki kimi sütun adı yaradır: 0->A, 1->B, ..., 25->Z, 26->AA...
func columnLetter(i int) string {
	letters := ""
	for i >= 0 {
		letters = string(rune('A'+i%26)) + letters
		i = i/26 - 1
	}
	return letters
}
