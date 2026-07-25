package handler

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"

	"accounting/internal/model"
	"accounting/internal/service"
)

var ibanRe = regexp.MustCompile(`^AZ[A-Z0-9]{20,}$`)

type AccountHandler struct {
	svc        *service.AccountService
	incomeSvc  *service.IncomeService
	expenseSvc *service.ExpenseService
	orderSvc   *service.OrderService
}

func NewAccountHandler(svc *service.AccountService, incomeSvc *service.IncomeService, expenseSvc *service.ExpenseService, orderSvc *service.OrderService) *AccountHandler {
	return &AccountHandler{svc: svc, incomeSvc: incomeSvc, expenseSvc: expenseSvc, orderSvc: orderSvc}
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

	xl, err := openExcelFromRequest(c)
	if err != nil {
		return
	}
	defer xl.Close()

	preview, err := parseStatement(xl, strings.TrimSpace(account.AccountNumber))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// Hansı sətirlərin artıq bu hesaba əlavə edildiyini işarələyirik ki, önizləmədə
	// istifadəçiyə göstərək (referansa görə).
	if existingIncomeRefs, err := h.incomeSvc.GetBankRefsByAccountID(c.Request.Context(), id); err == nil {
		for i := range preview.Gelirler {
			if preview.Gelirler[i].Ref != "" && existingIncomeRefs[preview.Gelirler[i].Ref] {
				preview.Gelirler[i].AlreadyImported = true
			}
		}
	}
	if existingExpenseRefs, err := h.expenseSvc.GetBankRefsByAccountID(c.Request.Context(), id); err == nil {
		for i := range preview.Giderler {
			if preview.Giderler[i].Ref != "" && existingExpenseRefs[preview.Giderler[i].Ref] {
				preview.Giderler[i].AlreadyImported = true
			}
		}
	}

	c.JSON(http.StatusOK, preview)
}

// openExcelFromRequest открывает Excel из multipart-поля "file" и пишет ошибку в c сам.
func openExcelFromRequest(c *gin.Context) (*excelize.File, error) {
	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return nil, err
	}
	f, err := fh.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot open file"})
		return nil, err
	}
	defer f.Close()
	xl, err := excelize.OpenReader(f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid Excel file"})
		return nil, err
	}
	return xl, nil
}

// parseStatement разбирает Excel-выписку и возвращает StatementPreview.
func parseStatement(xl *excelize.File, expectedIBAN string) (*model.StatementPreview, error) {
	sheetName := xl.GetSheetName(0)
	rows, err := xl.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	// Ищем IBAN
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
		return nil, fmt.Errorf("IBAN not found in file")
	}
	if fileIban != expectedIBAN {
		return nil, fmt.Errorf("IBAN mismatch: file has %s, account has %s", fileIban, expectedIBAN)
	}

	// Ищем строку заголовка
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
		return nil, fmt.Errorf("unrecognised Excel format")
	}

	get := func(row []string, idx int) string {
		if idx < 0 || idx >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[idx])
	}
	getF := func(row []string, idx int) float64 {
		v := strings.ReplaceAll(get(row, idx), ",", "")
		if v == "" {
			return 0
		}
		var f float64
		fmt.Sscanf(v, "%f", &f)
		return f
	}

	preview := &model.StatementPreview{IBAN: fileIban}
	for _, row := range rows[hdrIdx+1:] {
		debit := getF(row, colDebit)
		credit := getF(row, colCredit)
		date := get(row, colDate)
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
	return preview, nil
}

func parseStatementDate(s string) (time.Time, error) {
	return time.Parse("02.01.2006", s)
}

// ImportIncomes bir banka ekstresi önizlemesindeki gelir satırlarını, her satır için
// ayrı seçilmiş kategori/tur/müşteri ile birlikte kaydeder. Tur+müşteri seçilmişse
// ilgili sipariş bulunur veya oluşturulur ve gelir ona bağlanır.
func (h *AccountHandler) ImportIncomes(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req model.ImportGelirlerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	existingRefs, err := h.incomeSvc.GetBankRefsByAccountID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	skipped := 0
	seen := map[string]bool{}
	orderCache := map[string]int64{}
	var reqs []model.CreateIncomeRequest
	for _, row := range req.Rows {
		if row.IncomeCategoryID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "income_category_id is required for all rows"})
			return
		}
		if row.Ref != "" && (existingRefs[row.Ref] || seen[row.Ref]) {
			skipped++
			continue
		}
		if row.Ref != "" {
			seen[row.Ref] = true
		}

		date, err := parseStatementDate(row.Date)
		if err != nil {
			skipped++
			continue
		}

		var orderID *int64
		if row.TourID != nil && row.ClientID != nil {
			cacheKey := fmt.Sprintf("%d:%d", *row.ClientID, *row.TourID)
			if cachedID, ok := orderCache[cacheKey]; ok {
				orderID = &cachedID
			} else {
				order, err := h.orderSvc.FindOrCreateByClientAndTour(c.Request.Context(), *row.ClientID, *row.TourID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "order: " + err.Error()})
					return
				}
				orderCache[cacheKey] = order.ID
				orderID = &order.ID
			}
		}

		name := row.Desc
		if name == "" {
			name = row.Ref
		}
		reqs = append(reqs, model.CreateIncomeRequest{
			Name:              name,
			Amount:            row.Credit,
			Date:              date.Format("2006-01-02"),
			IncomeCategoryID:  row.IncomeCategoryID,
			AccountID:         id,
			TourID:            row.TourID,
			OrderID:           orderID,
			BankRef:           row.Ref,
			Counterparty:      row.CP,
			CounterpartyTaxID: row.Tax,
		})
	}

	imported := 0
	if len(reqs) > 0 {
		created, err := h.incomeSvc.BulkCreate(c.Request.Context(), reqs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		imported = len(created)
	}

	c.JSON(http.StatusOK, gin.H{
		"imported_incomes":   imported,
		"skipped_duplicates": skipped,
	})
}

// ImportExpenses bir banka ekstresi önizlemesindeki gider satırlarını, tüm satırlar
// için tek bir kategori seçimiyle kaydeder.
func (h *AccountHandler) ImportExpenses(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req model.ImportGiderlerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if req.ExpenseCategoryID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expense_category_id is required"})
		return
	}

	existingRefs, err := h.expenseSvc.GetBankRefsByAccountID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	skipped := 0
	seen := map[string]bool{}
	var reqs []model.CreateExpenseRequest
	for _, row := range req.Rows {
		if row.Ref != "" && (existingRefs[row.Ref] || seen[row.Ref]) {
			skipped++
			continue
		}
		if row.Ref != "" {
			seen[row.Ref] = true
		}

		date, err := parseStatementDate(row.Date)
		if err != nil {
			skipped++
			continue
		}

		name := row.Desc
		if name == "" {
			name = row.Ref
		}
		reqs = append(reqs, model.CreateExpenseRequest{
			Name:              name,
			Amount:            row.Debit,
			Date:              date.Format("2006-01-02"),
			ExpenseCategoryID: req.ExpenseCategoryID,
			AccountID:         id,
			BankRef:           row.Ref,
			Counterparty:      row.CP,
			CounterpartyTaxID: row.Tax,
		})
	}

	imported := 0
	if len(reqs) > 0 {
		created, err := h.expenseSvc.BulkCreate(c.Request.Context(), reqs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		imported = len(created)
	}

	c.JSON(http.StatusOK, gin.H{
		"imported_expenses":  imported,
		"skipped_duplicates": skipped,
	})
}
