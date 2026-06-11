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
	pages *PageHandler,
	tmpl *template.Template,
) *gin.Engine {
	r := gin.Default()
	r.SetHTMLTemplate(tmpl)
	r.Static("/static", "web/static")

	r.GET("/", func(c *gin.Context) { c.Redirect(http.StatusFound, "/login") })
	r.GET("/login", pages.Login)
	r.GET("/accounts", pages.Accounts)
	r.GET("/incomes", pages.Incomes)
	r.GET("/income-categories", pages.IncomeCategories)

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
		api.GET("/incomes/:id", incomes.GetByID)
		api.PUT("/incomes/:id", incomes.Update)
		api.DELETE("/incomes/:id", incomes.Delete)
	}

	return r
}
