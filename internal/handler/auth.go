package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"accounting/internal/auth"
	"accounting/internal/service"
)

const sessionCookie = "session_token"

type AuthHandler struct {
	userSvc  *service.UserService
	sessions *auth.Store
}

func NewAuthHandler(userSvc *service.UserService, sessions *auth.Store) *AuthHandler {
	return &AuthHandler{userSvc: userSvc, sessions: sessions}
}

func (h *AuthHandler) Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	user, err := h.userSvc.Login(c.Request.Context(), email, password)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "E-posta veya şifre hatalı"})
		return
	}

	token := h.sessions.Create(user.ID)
	c.SetCookie(sessionCookie, token, int((7 * 24 * time.Hour).Seconds()), "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

func (h *AuthHandler) Logout(c *gin.Context) {
	if token, err := c.Cookie(sessionCookie); err == nil && token != "" {
		h.sessions.Delete(token)
	}
	c.SetCookie(sessionCookie, "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}
