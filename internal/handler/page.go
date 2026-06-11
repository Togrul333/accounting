package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"accounting/internal/model"
	"accounting/internal/service"
)

type PageHandler struct {
	accountSvc        *service.AccountService
	incomeCategorySvc *service.IncomeCategoryService
	incomeSvc         *service.IncomeService
}

func NewPageHandler(
	accountSvc *service.AccountService,
	incomeCategorySvc *service.IncomeCategoryService,
	incomeSvc *service.IncomeService,
) *PageHandler {
	return &PageHandler{
		accountSvc:        accountSvc,
		incomeCategorySvc: incomeCategorySvc,
		incomeSvc:         incomeSvc,
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
