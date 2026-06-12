package handler

import (
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"accounting/internal/model"
	"accounting/internal/service"
)

type dashboardTransaction struct {
	IsIncome     bool
	Name         string
	CategoryName string
	Amount       float64
	Date         time.Time
}

type PageHandler struct {
	accountSvc         *service.AccountService
	incomeCategorySvc  *service.IncomeCategoryService
	incomeSvc          *service.IncomeService
	expenseCategorySvc *service.ExpenseCategoryService
	expenseSvc         *service.ExpenseService
	tourCategorySvc    *service.TourCategoryService
	roomSvc            *service.RoomService
	tourSvc            *service.TourService
	clientSvc          *service.ClientService
	settingSvc         *service.SettingService
	userSvc            *service.UserService
}

func NewPageHandler(
	accountSvc *service.AccountService,
	incomeCategorySvc *service.IncomeCategoryService,
	incomeSvc *service.IncomeService,
	expenseCategorySvc *service.ExpenseCategoryService,
	expenseSvc *service.ExpenseService,
	tourCategorySvc *service.TourCategoryService,
	roomSvc *service.RoomService,
	tourSvc *service.TourService,
	clientSvc *service.ClientService,
	settingSvc *service.SettingService,
	userSvc *service.UserService,
) *PageHandler {
	return &PageHandler{
		accountSvc:         accountSvc,
		incomeCategorySvc:  incomeCategorySvc,
		incomeSvc:          incomeSvc,
		expenseCategorySvc: expenseCategorySvc,
		expenseSvc:         expenseSvc,
		tourCategorySvc:    tourCategorySvc,
		roomSvc:            roomSvc,
		tourSvc:            tourSvc,
		clientSvc:          clientSvc,
		settingSvc:         settingSvc,
		userSvc:            userSvc,
	}
}

func (h *PageHandler) Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func (h *PageHandler) Profile(c *gin.Context) {
	user, err := h.userSvc.GetByID(c.Request.Context(), 1)
	if err != nil {
		log.Printf("profile page error: %v", err)
		c.Redirect(http.StatusFound, "/")
		return
	}
	c.HTML(http.StatusOK, "profile.html", gin.H{
		"active": "",
		"user":   user,
	})
}

func (h *PageHandler) Settings(c *gin.Context) {
	rates, err := h.settingSvc.GetRates(c.Request.Context())
	if err != nil {
		rates = model.ExchangeRates{}
	}
	c.HTML(http.StatusOK, "settings.html", gin.H{
		"active": "",
		"rates":  rates,
	})
}

func (h *PageHandler) Dashboard(c *gin.Context) {
	ctx := c.Request.Context()

	accounts, _ := h.accountSvc.GetAll(ctx)
	if accounts == nil {
		accounts = []model.Account{}
	}
	var totalBalance float64
	for _, a := range accounts {
		totalBalance += a.Balance
	}

	incomes, _ := h.incomeSvc.GetAll(ctx)
	if incomes == nil {
		incomes = []model.Income{}
	}
	var totalIncome float64
	for _, inc := range incomes {
		totalIncome += inc.Amount
	}

	expenses, _ := h.expenseSvc.GetAll(ctx)
	if expenses == nil {
		expenses = []model.Expense{}
	}
	var totalExpense float64
	for _, exp := range expenses {
		totalExpense += exp.Amount
	}

	tours, _ := h.tourSvc.GetAll(ctx)
	if tours == nil {
		tours = []model.Tour{}
	}
	now := time.Now()
	var activeTourCount int
	for _, t := range tours {
		if !now.Before(t.StartDate) && !now.After(t.EndDate) {
			activeTourCount++
		}
	}

	clients, _ := h.clientSvc.GetAll(ctx)
	if clients == nil {
		clients = []model.Client{}
	}

	// Son 5 tur (başlangıç tarihine göre azalan)
	sortedTours := make([]model.Tour, len(tours))
	copy(sortedTours, tours)
	sort.Slice(sortedTours, func(i, j int) bool {
		return sortedTours[i].StartDate.After(sortedTours[j].StartDate)
	})
	if len(sortedTours) > 5 {
		sortedTours = sortedTours[:5]
	}

	// Son 10 işlem (gelir + gider karışık, tarihe göre azalan)
	var txs []dashboardTransaction
	for _, inc := range incomes {
		txs = append(txs, dashboardTransaction{
			IsIncome:     true,
			Name:         inc.Name,
			CategoryName: inc.IncomeCategoryName,
			Amount:       inc.Amount,
			Date:         inc.Date,
		})
	}
	for _, exp := range expenses {
		txs = append(txs, dashboardTransaction{
			IsIncome:     false,
			Name:         exp.Name,
			CategoryName: exp.ExpenseCategoryName,
			Amount:       exp.Amount,
			Date:         exp.Date,
		})
	}
	sort.Slice(txs, func(i, j int) bool {
		return txs[i].Date.After(txs[j].Date)
	})
	if len(txs) > 8 {
		txs = txs[:8]
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"active":               "dashboard",
		"totalBalance":         totalBalance,
		"accountCount":         len(accounts),
		"totalIncome":          totalIncome,
		"incomeCount":          len(incomes),
		"totalExpense":         totalExpense,
		"expenseCount":         len(expenses),
		"net":                  totalIncome - totalExpense,
		"tourCount":            len(tours),
		"activeTourCount":      activeTourCount,
		"clientCount":          len(clients),
		"recentTours":          sortedTours,
		"recentTransactions":   txs,
	})
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

func (h *PageHandler) AccountEdit(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Redirect(http.StatusFound, "/accounts")
		return
	}
	account, err := h.accountSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("account edit page error: %v", err)
		c.Redirect(http.StatusFound, "/accounts")
		return
	}
	incomes, err := h.incomeSvc.GetByAccountID(c.Request.Context(), id)
	if err != nil {
		log.Printf("account incomes error: %v", err)
		incomes = []model.Income{}
	}
	if incomes == nil {
		incomes = []model.Income{}
	}
	expenses, err := h.expenseSvc.GetByAccountID(c.Request.Context(), id)
	if err != nil {
		log.Printf("account expenses error: %v", err)
		expenses = []model.Expense{}
	}
	if expenses == nil {
		expenses = []model.Expense{}
	}
	c.HTML(http.StatusOK, "account_edit.html", gin.H{
		"account":  account,
		"incomes":  incomes,
		"expenses": expenses,
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

func (h *PageHandler) TourCategories(c *gin.Context) {
	cats, err := h.tourCategorySvc.GetAll(c.Request.Context())
	if err != nil {
		log.Printf("tour categories page error: %v", err)
		cats = []model.TourCategory{}
	}
	if cats == nil {
		cats = []model.TourCategory{}
	}
	c.HTML(http.StatusOK, "tour_categories.html", gin.H{
		"categories": cats,
		"active":     "tour-categories",
	})
}

func (h *PageHandler) Rooms(c *gin.Context) {
	rooms, err := h.roomSvc.GetAll(c.Request.Context())
	if err != nil {
		log.Printf("rooms page error: %v", err)
		rooms = []model.Room{}
	}
	if rooms == nil {
		rooms = []model.Room{}
	}
	c.HTML(http.StatusOK, "rooms.html", gin.H{
		"rooms":  rooms,
		"active": "rooms",
	})
}

func (h *PageHandler) Tours(c *gin.Context) {
	tours, err := h.tourSvc.GetAll(c.Request.Context())
	if err != nil {
		log.Printf("tours page error: %v", err)
		tours = []model.Tour{}
	}
	if tours == nil {
		tours = []model.Tour{}
	}

	cats, err := h.tourCategorySvc.GetAll(c.Request.Context())
	if err != nil {
		cats = []model.TourCategory{}
	}

	rooms, err := h.roomSvc.GetAll(c.Request.Context())
	if err != nil {
		rooms = []model.Room{}
	}

	c.HTML(http.StatusOK, "tours.html", gin.H{
		"tours":      tours,
		"categories": cats,
		"rooms":      rooms,
		"active":     "tours",
	})
}

func (h *PageHandler) Clients(c *gin.Context) {
	clients, err := h.clientSvc.GetAll(c.Request.Context())
	if err != nil {
		log.Printf("clients page error: %v", err)
		clients = []model.Client{}
	}
	if clients == nil {
		clients = []model.Client{}
	}
	c.HTML(http.StatusOK, "clients.html", gin.H{
		"clients": clients,
		"active":  "clients",
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
