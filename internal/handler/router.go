package handler

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(accounts *AccountHandler, pages *PageHandler, tmpl *template.Template) *gin.Engine {
	r := gin.Default()
	r.SetHTMLTemplate(tmpl)
	r.Static("/static", "web/static")

	r.GET("/", func(c *gin.Context) { c.Redirect(http.StatusFound, "/login") })
	r.GET("/login", pages.Login)
	r.GET("/accounts", pages.Accounts)

	api := r.Group("/api")
	{
		api.GET("/accounts", accounts.GetAll)
		api.POST("/accounts", accounts.Create)
		api.GET("/accounts/:id", accounts.GetByID)
		api.PUT("/accounts/:id", accounts.Update)
		api.DELETE("/accounts/:id", accounts.Delete)
	}

	return r
}
