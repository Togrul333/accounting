package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"accounting/internal/model"
	"accounting/internal/service"
)

type PageHandler struct {
	accountSvc *service.AccountService
}

func NewPageHandler(accountSvc *service.AccountService) *PageHandler {
	return &PageHandler{accountSvc: accountSvc}
}

func (h *PageHandler) Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func (h *PageHandler) Accounts(c *gin.Context) {
	accounts, err := h.accountSvc.GetAll(c.Request.Context())
	if err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", nil)
		return
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
