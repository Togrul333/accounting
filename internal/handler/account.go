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
}

func NewAccountHandler(svc *service.AccountService, incomeSvc *service.IncomeService, expenseSvc *service.ExpenseService) *AccountHandler {
	return &AccountHandler{svc: svc, incomeSvc: incomeSvc, expenseSvc: expenseSvc}
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
	return preview, nil
}

func parseStatementDate(s string) (time.Time, error) {
	return time.Parse("02.01.2006", s)
}

func (h *AccountHandler) ImportStatement(c *gin.Context) {
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

	incomeCatID, err := strconv.ParseInt(c.PostForm("income_category_id"), 10, 64)
	if err != nil || incomeCatID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "income_category_id is required"})
		return
	}
	expenseCatID, err := strconv.ParseInt(c.PostForm("expense_category_id"), 10, 64)
	if err != nil || expenseCatID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expense_category_id is required"})
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

	preview, err := parseStatement(xl, strings.TrimSpace(account.AccountNumber))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// DD.MM.YYYY → YYYY-MM-DD
	convertDate := func(s string) string {
		t, err := parseStatementDate(s)
		if err != nil {
			return ""
		}
		return t.Format("2006-01-02")
	}

	var incomeReqs []model.CreateIncomeRequest
	for _, row := range preview.Gelirler {
		name := row.Desc
		if name == "" {
			name = row.Ref
		}
		incomeReqs = append(incomeReqs, model.CreateIncomeRequest{
			Name:              name,
			Amount:            row.Credit,
			Date:              convertDate(row.Date),
			IncomeCategoryID:  incomeCatID,
			AccountID:         id,
			BankRef:           row.Ref,
			Counterparty:      row.CP,
			CounterpartyTaxID: row.Tax,
		})
	}

	var expenseReqs []model.CreateExpenseRequest
	for _, row := range preview.Giderler {
		name := row.Desc
		if name == "" {
			name = row.Ref
		}
		expenseReqs = append(expenseReqs, model.CreateExpenseRequest{
			Name:              name,
			Amount:            row.Debit,
			Date:              convertDate(row.Date),
			ExpenseCategoryID: expenseCatID,
			AccountID:         id,
			BankRef:           row.Ref,
			Counterparty:      row.CP,
			CounterpartyTaxID: row.Tax,
		})
	}

	importedIncomes := 0
	if len(incomeReqs) > 0 {
		created, err := h.incomeSvc.BulkCreate(c.Request.Context(), incomeReqs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "incomes: " + err.Error()})
			return
		}
		importedIncomes = len(created)
	}

	importedExpenses := 0
	if len(expenseReqs) > 0 {
		created, err := h.expenseSvc.BulkCreate(c.Request.Context(), expenseReqs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "expenses: " + err.Error()})
			return
		}
		importedExpenses = len(created)
	}

	c.JSON(http.StatusOK, gin.H{
		"imported_incomes":  importedIncomes,
		"imported_expenses": importedExpenses,
	})
}
