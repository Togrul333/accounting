package handler

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	accounts *AccountHandler,
	incomeCategories *IncomeCategoryHandler,
	incomes *IncomeHandler,
	expenseCategories *ExpenseCategoryHandler,
	expenses *ExpenseHandler,
	tourCategories *TourCategoryHandler,
	rooms *RoomHandler,
	tours *TourHandler,
	clients *ClientHandler,
	settings *SettingHandler,
	users *UserHandler,
	discountCategories *DiscountCategoryHandler,
	discounts *DiscountHandler,
	pages *PageHandler,
	tmpl *template.Template,
) *gin.Engine {
	r := gin.Default()
	r.SetHTMLTemplate(tmpl)
	r.Static("/static", "web/static")

	r.GET("/", pages.Dashboard)
	r.GET("/login", pages.Login)
	r.GET("/profile", pages.Profile)
	r.GET("/settings", pages.Settings)
	r.POST("/logout", func(c *gin.Context) { c.Redirect(http.StatusFound, "/login") })
	r.GET("/accounts", pages.Accounts)
	r.GET("/accounts/:id/edit", pages.AccountEdit)
	r.GET("/incomes", pages.Incomes)
	r.GET("/income-categories", pages.IncomeCategories)
	r.GET("/expenses", pages.Expenses)
	r.GET("/expense-categories", pages.ExpenseCategories)
	r.GET("/tours", pages.Tours)
	r.GET("/tour-categories", pages.TourCategories)
	r.GET("/rooms", pages.Rooms)
	r.GET("/clients", pages.Clients)
	r.GET("/discounts", pages.Discounts)
	r.GET("/discount-categories", pages.DiscountCategories)

	api := r.Group("/api")
	{
		api.GET("/accounts", accounts.GetAll)
		api.POST("/accounts", accounts.Create)
		api.GET("/accounts/:id", accounts.GetByID)
		api.PUT("/accounts/:id", accounts.Update)
		api.DELETE("/accounts/:id", accounts.Delete)

		api.GET("/income-categories", incomeCategories.GetAll)
		api.POST("/income-categories", incomeCategories.Create)
		api.GET("/income-categories/:id", incomeCategories.GetByID)
		api.PUT("/income-categories/:id", incomeCategories.Update)
		api.DELETE("/income-categories/:id", incomeCategories.Delete)

		api.GET("/incomes", incomes.GetAll)
		api.POST("/incomes", incomes.Create)
		api.POST("/incomes/bulk", incomes.BulkCreate)
		api.GET("/incomes/:id", incomes.GetByID)
		api.PUT("/incomes/:id", incomes.Update)
		api.DELETE("/incomes/:id", incomes.Delete)

		api.GET("/expense-categories", expenseCategories.GetAll)
		api.POST("/expense-categories", expenseCategories.Create)
		api.GET("/expense-categories/:id", expenseCategories.GetByID)
		api.PUT("/expense-categories/:id", expenseCategories.Update)
		api.DELETE("/expense-categories/:id", expenseCategories.Delete)

		api.GET("/expenses", expenses.GetAll)
		api.POST("/expenses", expenses.Create)
		api.POST("/expenses/bulk", expenses.BulkCreate)
		api.GET("/expenses/:id", expenses.GetByID)
		api.PUT("/expenses/:id", expenses.Update)
		api.DELETE("/expenses/:id", expenses.Delete)

		api.GET("/tour-categories", tourCategories.GetAll)
		api.POST("/tour-categories", tourCategories.Create)
		api.GET("/tour-categories/:id", tourCategories.GetByID)
		api.PUT("/tour-categories/:id", tourCategories.Update)
		api.DELETE("/tour-categories/:id", tourCategories.Delete)

		api.GET("/rooms", rooms.GetAll)
		api.POST("/rooms", rooms.Create)
		api.GET("/rooms/:id", rooms.GetByID)
		api.PUT("/rooms/:id", rooms.Update)
		api.DELETE("/rooms/:id", rooms.Delete)

		api.GET("/tours", tours.GetAll)
		api.POST("/tours", tours.Create)
		api.GET("/tours/:id", tours.GetByID)
		api.PUT("/tours/:id", tours.Update)
		api.DELETE("/tours/:id", tours.Delete)

		api.GET("/clients", clients.GetAll)
		api.POST("/clients", clients.Create)
		api.GET("/clients/:id", clients.GetByID)
		api.PUT("/clients/:id", clients.Update)
		api.DELETE("/clients/:id", clients.Delete)

		api.PUT("/settings/rates", settings.UpdateRates)

		api.PUT("/profile", users.UpdateProfile)
		api.PUT("/profile/password", users.UpdatePassword)

		api.GET("/discount-categories", discountCategories.GetAll)
		api.POST("/discount-categories", discountCategories.Create)
		api.GET("/discount-categories/:id", discountCategories.GetByID)
		api.PUT("/discount-categories/:id", discountCategories.Update)
		api.DELETE("/discount-categories/:id", discountCategories.Delete)

		api.GET("/discounts", discounts.GetAll)
		api.POST("/discounts", discounts.Create)
		api.GET("/discounts/:id", discounts.GetByID)
		api.PUT("/discounts/:id", discounts.Update)
		api.DELETE("/discounts/:id", discounts.Delete)
	}

	return r
}
