package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"accounting/internal/model"
	"accounting/internal/service"
)

type PageHandler struct {
	accountSvc         *service.AccountService
	incomeCategorySvc  *service.IncomeCategoryService
	incomeSvc          *service.IncomeService
	expenseCategorySvc *service.ExpenseCategoryService
	expenseSvc         *service.ExpenseService
}

func NewPageHandler(
	accountSvc *service.AccountService,
	incomeCategorySvc *service.IncomeCategoryService,
	incomeSvc *service.IncomeService,
	expenseCategorySvc *service.ExpenseCategoryService,
	expenseSvc *service.ExpenseService,
) *PageHandler {
	return &PageHandler{
		accountSvc:         accountSvc,
		incomeCategorySvc:  incomeCategorySvc,
		incomeSvc:          incomeSvc,
		expenseCategorySvc: expenseCategorySvc,
		expenseSvc:         expenseSvc,
	}
}

func (h *PageHandler) Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func (h *PageHandler) Accounts(c *gin.Context) {
	accounts, err := h.accountSvc.GetAll(c.Request.Context())
	if err != nil {
		log.Printf("accounts page error: %v", err)
		accounts = []model.Account{}
	}
	if accounts == nil {
		accounts = []model.Account{}
	}

	var total float64
	for _, a := range accounts {
		total += a.Balance
	}

	c.HTML(http.StatusOK, "accounts.html", gin.H{
		"accounts": accounts,
		"total":    total,
		"active":   "accounts",
	})
}

func (h *PageHandler) IncomeCategories(c *gin.Context) {
	cats, err := h.incomeCategorySvc.GetAll(c.Request.Context())
	if err != nil {
		log.Printf("income categories page error: %v", err)
		cats = []model.IncomeCategory{}
	}
	if cats == nil {
		cats = []model.IncomeCategory{}
	}
	c.HTML(http.StatusOK, "income_categories.html", gin.H{
		"categories": cats,
		"active":     "income-categories",
	})
}

func (h *PageHandler) Incomes(c *gin.Context) {
	incomes, err := h.incomeSvc.GetAll(c.Request.Context())
	if err != nil {
		log.Printf("incomes page error: %v", err)
		incomes = []model.Income{}
	}
	if incomes == nil {
		incomes = []model.Income{}
	}

	cats, err := h.incomeCategorySvc.GetAll(c.Request.Context())
	if err != nil {
		cats = []model.IncomeCategory{}
	}

	accounts, err := h.accountSvc.GetAll(c.Request.Context())
	if err != nil {
		accounts = []model.Account{}
	}

	var total float64
	for _, inc := range incomes {
		total += inc.Amount
	}

	c.HTML(http.StatusOK, "incomes.html", gin.H{
		"incomes":    incomes,
		"categories": cats,
		"accounts":   accounts,
		"total":      total,
		"active":     "incomes",
	})
}

func (h *PageHandler) ExpenseCategories(c *gin.Context) {
	cats, err := h.expenseCategorySvc.GetAll(c.Request.Context())
	if err != nil {
		log.Printf("expense categories page error: %v", err)
		cats = []model.ExpenseCategory{}
	}
	if cats == nil {
		cats = []model.ExpenseCategory{}
	}
	c.HTML(http.StatusOK, "expense_categories.html", gin.H{
		"categories": cats,
		"active":     "expense-categories",
	})
}

func (h *PageHandler) Expenses(c *gin.Context) {
	expenses, err := h.expenseSvc.GetAll(c.Request.Context())
	if err != nil {
		log.Printf("expenses page error: %v", err)
		expenses = []model.Expense{}
	}
	if expenses == nil {
		expenses = []model.Expense{}
	}

	cats, err := h.expenseCategorySvc.GetAll(c.Request.Context())
	if err != nil {
		cats = []model.ExpenseCategory{}
	}

	accounts, err := h.accountSvc.GetAll(c.Request.Context())
	if err != nil {
		accounts = []model.Account{}
	}

	var total float64
	for _, exp := range expenses {
		total += exp.Amount
	}

	c.HTML(http.StatusOK, "expenses.html", gin.H{
		"expenses":   expenses,
		"categories": cats,
		"accounts":   accounts,
		"total":      total,
		"active":     "expenses",
	})
}
