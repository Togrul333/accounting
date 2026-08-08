package handler

import (
	"html/template"

	"github.com/gin-gonic/gin"

	"accounting/internal/auth"
)

func NewRouter(
	authHandler *AuthHandler,
	sessions *auth.Store,
	accounts *AccountHandler,
	incomeCategories *IncomeCategoryHandler,
	incomes *IncomeHandler,
	expenseCategories *ExpenseCategoryHandler,
	expenses *ExpenseHandler,
	tourCategories *TourCategoryHandler,
	rooms *RoomHandler,
	flights *FlightHandler,
	tours *TourHandler,
	clients *ClientHandler,
	settings *SettingHandler,
	referansUsers *ReferansUserHandler,
	hocaUsers *HocaUserHandler,
	defaultTasks *DefaultTaskHandler,
	users *UserHandler,
	discountCategories *DiscountCategoryHandler,
	discounts *DiscountHandler,
	orders *OrderHandler,
	tasks *TaskHandler,
	taskComments *TaskCommentHandler,
	telegram *TelegramHandler,
	pages *PageHandler,
	sheetsImport *SheetsImportHandler,
	metaAds *MetaAdsHandler,
	tmpl *template.Template,
) *gin.Engine {
	r := gin.Default()
	r.SetHTMLTemplate(tmpl)
	r.Static("/static", "web/static")
	r.Use(AuthMiddleware(sessions))

	r.GET("/", pages.Dashboard)
	r.GET("/login", pages.Login)
	r.POST("/login", authHandler.Login)
	r.GET("/profile", pages.Profile)
	r.GET("/settings", pages.Settings)
	r.GET("/referans-users/:id", pages.ReferansUserShow)
	r.POST("/logout", authHandler.Logout)
	r.GET("/accounts", pages.Accounts)
	r.GET("/accounts/:id/edit", pages.AccountEdit)
	r.GET("/accounts/:id/incomes", pages.AccountIncomes)
	r.GET("/accounts/:id/expenses", pages.AccountExpenses)
	r.GET("/incomes", pages.Incomes)
	r.GET("/income-categories", pages.IncomeCategories)
	r.GET("/expenses", pages.Expenses)
	r.GET("/expense-categories", pages.ExpenseCategories)
	r.GET("/tours", pages.Tours)
	r.GET("/tours/:id", pages.TourShow)
	r.GET("/tours/:id/edit", pages.TourEdit)
	r.GET("/tour-categories", pages.TourCategories)
	r.GET("/rooms", pages.Rooms)
	r.GET("/flights", pages.Flights)
	r.GET("/flights/:id", pages.FlightShow)
	r.GET("/clients", pages.Clients)
	r.GET("/discounts", pages.Discounts)
	r.GET("/discount-categories", pages.DiscountCategories)
	r.GET("/orders", pages.Orders)
	r.GET("/orders/:id", pages.OrderShow)
	r.GET("/orders/:id/edit", pages.OrderEdit)
	r.GET("/sheets-import", pages.SheetsImport)
	r.GET("/meta-ads", pages.MetaAds)
	r.GET("/tasks", pages.Tasks)

	api := r.Group("/api")
	{
		api.GET("/accounts", accounts.GetAll)
		api.POST("/accounts", accounts.Create)
		api.GET("/accounts/:id", accounts.GetByID)
		api.PUT("/accounts/:id", accounts.Update)
		api.DELETE("/accounts/:id", accounts.Delete)
		api.POST("/accounts/:id/statement-preview", accounts.ParseStatement)
		api.POST("/accounts/:id/statement-import-incomes", accounts.ImportIncomes)
		api.POST("/accounts/:id/statement-import-expenses", accounts.ImportExpenses)

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

		api.GET("/flights", flights.GetAll)
		api.POST("/flights", flights.Create)
		api.GET("/flights/:id", flights.GetByID)
		api.PUT("/flights/:id", flights.Update)
		api.DELETE("/flights/:id", flights.Delete)

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
		api.POST("/clients/:id/document", clients.UploadDocument)

		api.PUT("/settings/rates", settings.UpdateRates)

		api.GET("/referans-users", referansUsers.GetAll)
		api.POST("/referans-users", referansUsers.Create)
		api.GET("/referans-users/:id", referansUsers.GetByID)
		api.PUT("/referans-users/:id", referansUsers.Update)
		api.DELETE("/referans-users/:id", referansUsers.Delete)
		api.GET("/referans-users/:id/candidates", referansUsers.Candidates)
		api.GET("/referans-users/:id/order-search", referansUsers.SearchOrders)
		api.GET("/referans-users/:id/orders", referansUsers.Referrals)
		api.POST("/referans-users/:id/orders", referansUsers.AddReferral)
		api.DELETE("/referans-users/:id/orders/:order_id", referansUsers.RemoveReferral)

		api.GET("/hoca-users", hocaUsers.GetAll)
		api.POST("/hoca-users", hocaUsers.Create)
		api.GET("/hoca-users/:id", hocaUsers.GetByID)
		api.PUT("/hoca-users/:id", hocaUsers.Update)
		api.DELETE("/hoca-users/:id", hocaUsers.Delete)

		api.GET("/default-tasks", defaultTasks.GetAll)
		api.POST("/default-tasks", defaultTasks.Create)
		api.GET("/default-tasks/:id", defaultTasks.GetByID)
		api.PUT("/default-tasks/:id", defaultTasks.Update)
		api.DELETE("/default-tasks/:id", defaultTasks.Delete)

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

		api.GET("/orders", orders.GetAll)
		api.POST("/orders", orders.Create)
		api.GET("/orders/:id", orders.GetByID)
		api.PUT("/orders/:id", orders.Update)
		api.DELETE("/orders/:id", orders.Delete)
		api.POST("/orders/:id/incomes", orders.AddIncome)
		api.POST("/orders/:id/discounts", orders.AddDiscount)

		api.GET("/tasks", tasks.GetAll)
		api.POST("/tasks", tasks.Create)
		api.GET("/tasks/:id", tasks.GetByID)
		api.PUT("/tasks/:id", tasks.Update)
		api.PUT("/tasks/:id/status", tasks.UpdateStatus)
		api.DELETE("/tasks/:id", tasks.Delete)
		api.GET("/tasks/:id/comments", taskComments.GetByTask)
		api.POST("/tasks/:id/comments", taskComments.Create)
		api.PUT("/task-comments/:id", taskComments.Update)
		api.DELETE("/task-comments/:id", taskComments.Delete)

		api.GET("/telegram/settings", telegram.GetSettings)
		api.PUT("/telegram/settings", telegram.UpdateSettings)
		api.POST("/telegram/sync", telegram.Sync)
		api.POST("/telegram/hoca-users/:id/link-code", telegram.LinkCode)
		api.POST("/telegram/hoca-users/:id/test", telegram.Test)
		api.DELETE("/telegram/hoca-users/:id", telegram.Unlink)

		api.GET("/sheets-import/links", sheetsImport.Links)
		api.POST("/sheets-import/tabs", sheetsImport.Tabs)
		api.POST("/sheets-import/preview", sheetsImport.Preview)
		api.POST("/sheets-import/turlar-candidates", sheetsImport.TurlarCandidates)
		api.POST("/sheets-import/passenger-candidates", sheetsImport.PassengerCandidates)

		api.GET("/meta-ads/accounts", metaAds.GetAccounts)
		api.POST("/meta-ads/accounts", metaAds.CreateAccount)
		api.GET("/meta-ads/accounts/:id", metaAds.GetAccount)
		api.PUT("/meta-ads/accounts/:id", metaAds.UpdateAccount)
		api.DELETE("/meta-ads/accounts/:id", metaAds.DeleteAccount)
		api.POST("/meta-ads/accounts/:id/verify", metaAds.Verify)
		api.POST("/meta-ads/accounts/:id/sync", metaAds.Sync)
		api.GET("/meta-ads/spend", metaAds.Spend)
		api.GET("/meta-ads/summary", metaAds.Summary)
		api.PUT("/meta-ads/spend/:id/tour", metaAds.SetTour)
	}

	return r
}
