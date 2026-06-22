package handler

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"

	"accounting/internal/model"
	"accounting/internal/service"
)

var ibanRe = regexp.MustCompile(`^AZ[A-Z0-9]{20,}$`)

type AccountHandler struct {
	svc *service.AccountService
}

func NewAccountHandler(svc *service.AccountService) *AccountHandler {
	return &AccountHandler{svc: svc}
}

func (h *AccountHandler) GetAll(c *gin.Context) {
	accounts, err := h.svc.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if accounts == nil {
		accounts = []model.Account{}
	}
	c.JSON(http.StatusOK, accounts)
}

func (h *AccountHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	account, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, account)
}

func (h *AccountHandler) Create(c *gin.Context) {
	var req model.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	account, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, account)
}

func (h *AccountHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req model.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	account, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, account)
}

func (h *AccountHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *AccountHandler) ParseStatement(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	account, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	f, err := fh.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot open file"})
		return
	}
	defer f.Close()

	xl, err := excelize.OpenReader(f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid Excel file"})
		return
	}
	defer xl.Close()

	sheetName := xl.GetSheetName(0)

	// Ищем IBAN — первая ячейка вида AZ...
	rows, err := xl.GetRows(sheetName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var fileIban string
	for _, row := range rows {
		for _, cell := range row {
			cell = strings.TrimSpace(cell)
			if ibanRe.MatchString(cell) {
				fileIban = cell
				break
			}
		}
		if fileIban != "" {
			break
		}
	}
	if fileIban == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "IBAN not found in file"})
		return
	}
	if fileIban != strings.TrimSpace(account.AccountNumber) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": fmt.Sprintf("IBAN mismatch: file has %s, account has %s", fileIban, account.AccountNumber),
		})
		return
	}

	// Ищем строку заголовка: та где есть 'Əməliyyatın tarixi' (точная фраза, не просто tarixi)
	hdrIdx := -1
	colDate, colRef, colCP, colDebit, colCredit, colDesc, colTax := -1, -1, -1, -1, -1, -1, -1
	for i, row := range rows {
		for _, cell := range row {
			if strings.Contains(cell, "Əməliyyatın tarixi") {
				hdrIdx = i
				for jj, h := range row {
					switch {
					case strings.Contains(h, "tarixi"):
						colDate = jj
					case strings.Contains(h, "referens"):
						colRef = jj
					case strings.Contains(h, "hesab"):
						colCP = jj
					case strings.Contains(h, "Debit"):
						colDebit = jj
					case strings.Contains(h, "Kredit"):
						colCredit = jj
					case strings.Contains(h, "yinat"):
						colDesc = jj
					case strings.Contains(h, "VÖEN"):
						colTax = jj
					}
				}
				break
			}
		}
		if hdrIdx >= 0 {
			break
		}
	}
	if hdrIdx < 0 || colDebit < 0 || colCredit < 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "unrecognised Excel format"})
		return
	}

	get := func(row []string, idx int) string {
		if idx < 0 || idx >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[idx])
	}
	// числа вида "50,000.00" → убираем запятые перед парсингом
	getF := func(row []string, idx int) float64 {
		v := strings.ReplaceAll(get(row, idx), ",", "")
		if v == "" {
			return 0
		}
		var f float64
		fmt.Sscanf(v, "%f", &f)
		return f
	}

	var preview model.StatementPreview
	preview.IBAN = fileIban

	for _, row := range rows[hdrIdx+1:] {
		debit  := getF(row, colDebit)
		credit := getF(row, colCredit)
		date   := get(row, colDate)
		if date == "" && debit == 0 && credit == 0 {
			continue
		}
		entry := model.StatementRow{
			Date:   date,
			Ref:    get(row, colRef),
			CP:     get(row, colCP),
			Debit:  debit,
			Credit: credit,
			Desc:   get(row, colDesc),
			Tax:    get(row, colTax),
		}
		if credit > 0 {
			preview.Gelirler = append(preview.Gelirler, entry)
			preview.TotalCredit += credit
		}
		if debit > 0 {
			preview.Giderler = append(preview.Giderler, entry)
			preview.TotalDebit += debit
		}
	}
	if preview.Gelirler == nil {
		preview.Gelirler = []model.StatementRow{}
	}
	if preview.Giderler == nil {
		preview.Giderler = []model.StatementRow{}
	}

	c.JSON(http.StatusOK, preview)
}
