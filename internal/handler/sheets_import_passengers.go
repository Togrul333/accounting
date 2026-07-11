package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type passengerMatch struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

type passengerCandidate struct {
	FirstName    string          `json:"first_name"`
	LastName     string          `json:"last_name"`
	BirthYear    int             `json:"birth_year"`
	Phone        string          `json:"phone"`
	PassportNo   string          `json:"passport_no"`
	TourCode     string          `json:"tour_code"`
	RoomCode     string          `json:"room_code"`
	CategoryName string          `json:"category_name"`
	Cancelled    bool            `json:"cancelled"`
	ClientMatch  *passengerMatch `json:"client_match"`
	TourMatch    *int64          `json:"tour_match"`
	HasOrder     bool            `json:"has_order"`
}

// PassengerCandidates konkret tur kodu vərəqindəki (məsələn "AZ2606001") sərnişin siyahısını oxuyur,
// hər sətri bizim Clients/Tours/Orders bazası ilə tutuşdurur.
func (h *SheetsImportHandler) PassengerCandidates(c *gin.Context) {
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

	cols, headerIdx := findPassengerHeader(rows)
	if headerIdx == -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "\"AD\" / \"SOYAD\" başlıqları tapılmadı"})
		return
	}

	ctx := c.Request.Context()

	clients, err := h.clientSvc.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	clientByKey := make(map[string]passengerMatch, len(clients))
	for _, cl := range clients {
		clientByKey[clientKey(cl.FirstName, cl.LastName, cl.BirthYear)] = passengerMatch{ID: cl.ID, FirstName: cl.FirstName, LastName: cl.LastName}
	}

	tours, err := h.tourSvc.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tourByKey := make(map[string]int64, len(tours))
	for _, t := range tours {
		tourByKey[tourKey(t.Code, t.RoomCode, t.TourCategoryName)] = t.ID
	}

	orders, err := h.orderSvc.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	orderPairs := make(map[string]bool, len(orders))
	for _, o := range orders {
		orderPairs[orderKey(o.ClientID, o.TourID)] = true
	}

	candidates := make([]passengerCandidate, 0)
	for _, row := range rows[headerIdx+1:] {
		firstName := cellAt(row, cols.first)
		lastName := cellAt(row, cols.last)
		if firstName == "" && lastName == "" {
			continue
		}

		birthYear := parseYear(cellAt(row, cols.birth))
		cand := passengerCandidate{
			FirstName:    firstName,
			LastName:     lastName,
			BirthYear:    birthYear,
			Phone:        cellAt(row, cols.phone),
			PassportNo:   cellAt(row, cols.passport),
			TourCode:     cellAt(row, cols.code),
			RoomCode:     cellAt(row, cols.room),
			CategoryName: cellAt(row, cols.category),
			Cancelled:    cellAt(row, cols.cancelled) == "1",
		}

		if m, ok := clientByKey[clientKey(cand.FirstName, cand.LastName, cand.BirthYear)]; ok {
			match := m
			cand.ClientMatch = &match
		}

		if tourID, ok := tourByKey[tourKey(cand.TourCode, cand.RoomCode, cand.CategoryName)]; ok {
			id := tourID
			cand.TourMatch = &id
			if cand.ClientMatch != nil {
				cand.HasOrder = orderPairs[orderKey(cand.ClientMatch.ID, tourID)]
			}
		}

		candidates = append(candidates, cand)
	}

	c.JSON(http.StatusOK, gin.H{"passengers": candidates})
}

type passengerCols struct {
	code, room, category, first, last, birth, phone, passport, cancelled int
}

// findPassengerHeader "AD" və "SOYAD" xanaları olan sətri axtarır (sərnişin cədvəllərinin başlıq sətri).
func findPassengerHeader(rows [][]string) (passengerCols, int) {
	cols := passengerCols{-1, -1, -1, -1, -1, -1, -1, -1, -1}

	for ri, row := range rows {
		firstIdx, lastIdx := -1, -1
		for ci, cell := range row {
			t := strings.TrimSpace(cell)
			if strings.EqualFold(t, "AD") {
				firstIdx = ci
			}
			if strings.EqualFold(t, "SOYAD") {
				lastIdx = ci
			}
		}
		if firstIdx == -1 || lastIdx == -1 {
			continue
		}

		cols.first, cols.last = firstIdx, lastIdx
		for ci, cell := range row {
			lower := strings.ToLower(strings.TrimSpace(cell))
			switch {
			case ci == firstIdx || ci == lastIdx:
				continue
			case strings.Contains(lower, "tur kodu"):
				cols.code = ci
			case strings.Contains(lower, "otaq") || strings.Contains(lower, "oda tip"):
				cols.room = ci
			case strings.Contains(lower, "konsept"):
				cols.category = ci
			case strings.Contains(lower, "doğum") || strings.Contains(lower, "dogum"):
				cols.birth = ci
			case strings.Contains(lower, "telefon"):
				cols.phone = ci
			case strings.Contains(lower, "pasaport"):
				cols.passport = ci
			case strings.Contains(lower, "iptal"):
				cols.cancelled = ci
			}
		}
		return cols, ri
	}
	return cols, -1
}

func clientKey(firstName, lastName string, birthYear int) string {
	return strings.ToLower(strings.TrimSpace(firstName)) + "|" + strings.ToLower(strings.TrimSpace(lastName)) + "|" + strconv.Itoa(birthYear)
}

func tourKey(code, roomCode, categoryName string) string {
	return strings.ToLower(strings.TrimSpace(code)) + "|" + strings.ToLower(strings.TrimSpace(roomCode)) + "|" + strings.ToLower(strings.TrimSpace(categoryName))
}

func orderKey(clientID, tourID int64) string {
	return strconv.FormatInt(clientID, 10) + "|" + strconv.FormatInt(tourID, 10)
}

// parseYear doğum tarixi mətnindən (müxtəlif formatlarda) il-i çıxarır.
func parseYear(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	for _, layout := range []string{"2006-01-02", "02.01.2006", "2006-1-2", "02/01/2006"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t.Year()
		}
	}
	return 0
}
