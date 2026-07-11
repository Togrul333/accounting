package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type turlarItem struct {
	RoomCode     string  `json:"room_code"`
	CategoryName string  `json:"category_name"`
	Price        float64 `json:"price"`
}

type turlarTour struct {
	Code      string       `json:"code"`
	StartDate string       `json:"start_date"`
	EndDate   string       `json:"end_date"`
	Exists    bool         `json:"exists"`
	Items     []turlarItem `json:"items"`
}

// TurlarCandidates "Turlar" vərəqini oxuyub tur kodlarına görə qruplaşdırır və bizim bazada
// artıq olub-olmadığını (Code üzrə) yoxlayır — import düymələri üçün.
func (h *SheetsImportHandler) TurlarCandidates(c *gin.Context) {
	var req sheetsPreviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cədvəl və vərəq tələb olunur"})
		return
	}

	if h.sheets == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Google Sheets bağlantısı quraşdırılmayıb"})
		return
	}

	_, rows, err := h.sheets.FetchRowsByGID(req.SpreadsheetID, req.GID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	headerIdx, codeCol, startCol, endCol, roomCol, categoryCols := findTurlarHeader(rows)
	if headerIdx == -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "\"Tur Kodu\" başlıq sətri tapılmadı"})
		return
	}
	header := rows[headerIdx]

	existingCodes, err := h.existingTourCodes(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	order := make([]string, 0)
	byCode := make(map[string]*turlarTour)

	for _, row := range rows[headerIdx+1:] {
		code := cellAt(row, codeCol)
		if code == "" {
			continue
		}
		startISO, ok1 := parseDMY(cellAt(row, startCol))
		endISO, ok2 := parseDMY(cellAt(row, endCol))
		if !ok1 || !ok2 {
			continue
		}
		roomCode := cellAt(row, roomCol)

		tour, found := byCode[code]
		if !found {
			tour = &turlarTour{
				Code:      code,
				StartDate: startISO,
				EndDate:   endISO,
				Exists:    existingCodes[code],
				Items:     []turlarItem{},
			}
			byCode[code] = tour
			order = append(order, code)
		}

		for _, colIdx := range categoryCols {
			priceRaw := cellAt(row, colIdx)
			price, ok := parsePrice(priceRaw)
			if !ok {
				continue
			}
			tour.Items = append(tour.Items, turlarItem{
				RoomCode:     roomCode,
				CategoryName: header[colIdx],
				Price:        price,
			})
		}
	}

	tours := make([]*turlarTour, 0, len(order))
	for _, code := range order {
		tours = append(tours, byCode[code])
	}

	c.JSON(http.StatusOK, gin.H{"tours": tours})
}

func (h *SheetsImportHandler) existingTourCodes(c *gin.Context) (map[string]bool, error) {
	tours, err := h.tourSvc.GetAll(c.Request.Context())
	if err != nil {
		return nil, err
	}
	codes := make(map[string]bool, len(tours))
	for _, t := range tours {
		codes[strings.TrimSpace(t.Code)] = true
	}
	return codes, nil
}

func cellAt(row []string, idx int) string {
	if idx < 0 || idx >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[idx])
}

// findTurlarHeader "Tur Kodu" olan sətri axtarır və lazımi sütunların indekslərini müəyyən edir.
func findTurlarHeader(rows [][]string) (headerIdx, codeCol, startCol, endCol, roomCol int, categoryCols []int) {
	headerIdx, codeCol, startCol, endCol, roomCol = -1, -1, -1, -1, -1

	for ri, row := range rows {
		for ci, cell := range row {
			if strings.EqualFold(strings.TrimSpace(cell), "Tur Kodu") {
				headerIdx = ri
				codeCol = ci
				break
			}
		}
		if headerIdx != -1 {
			break
		}
	}
	if headerIdx == -1 {
		return
	}

	header := rows[headerIdx]
	for i, cell := range header {
		lower := strings.ToLower(strings.TrimSpace(cell))
		switch {
		case i == codeCol:
			continue
		case strings.Contains(lower, "başlan") || strings.Contains(lower, "baslan"):
			startCol = i
		case strings.Contains(lower, "bitiş") || strings.Contains(lower, "bitis"):
			endCol = i
		case strings.Contains(lower, "otaq") || strings.Contains(lower, "oda"):
			roomCol = i
		case lower != "":
			categoryCols = append(categoryCols, i)
		}
	}
	return
}

// parseDMY "24.06.2026" formatını "2026-06-24" (ISO) formatına çevirir.
func parseDMY(s string) (string, bool) {
	t, err := time.Parse("02.01.2006", s)
	if err != nil {
		return "", false
	}
	return t.Format("2006-01-02"), true
}

// parsePrice "$1,250.00" kimi mətni float-a çevirir.
func parsePrice(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	s = strings.NewReplacer("$", "", ",", "", " ", "").Replace(s)
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || v == 0 {
		return 0, false
	}
	return v, true
}
