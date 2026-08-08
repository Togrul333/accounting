package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"accounting/internal/model"
	"accounting/internal/service"
	"accounting/internal/telegram"
)

type TelegramHandler struct {
	svc *service.TelegramService
}

func NewTelegramHandler(svc *service.TelegramService) *TelegramHandler {
	return &TelegramHandler{svc: svc}
}

// telegramError — ошибки настройки это 400, а не 500: их чинит пользователь.
func telegramError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrTelegramDisabled),
		errors.Is(err, service.ErrTelegramNoToken),
		errors.Is(err, service.ErrTelegramNoChat):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, telegram.ErrChatNotStarted):
		c.JSON(http.StatusBadRequest, gin.H{"error": telegram.ErrChatNotStarted.Error()})
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "hoca not found"})
	case telegram.IsAPIError(err):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (h *TelegramHandler) GetSettings(c *gin.Context) {
	settings, err := h.svc.GetSettings(c.Request.Context())
	if err != nil {
		telegramError(c, err)
		return
	}
	c.JSON(http.StatusOK, settings)
}

func (h *TelegramHandler) UpdateSettings(c *gin.Context) {
	var req model.UpdateTelegramSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	settings, err := h.svc.UpdateSettings(c.Request.Context(), req)
	if err != nil {
		telegramError(c, err)
		return
	}
	c.JSON(http.StatusOK, settings)
}

// LinkCode — разовый код и ссылка t.me для привязки сотрудника.
func (h *TelegramHandler) LinkCode(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	code, link, err := h.svc.LinkCode(c.Request.Context(), id)
	if err != nil {
		telegramError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": code, "link": link})
}

func (h *TelegramHandler) Unlink(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.Unlink(c.Request.Context(), id); err != nil {
		telegramError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Sync — разбирает новые /start и привязывает чаты к сотрудникам.
func (h *TelegramHandler) Sync(c *gin.Context) {
	linked, err := h.svc.SyncLinks(c.Request.Context())
	if err != nil {
		telegramError(c, err)
		return
	}
	if linked == nil {
		linked = []string{}
	}
	c.JSON(http.StatusOK, gin.H{"linked": linked, "count": len(linked)})
}

func (h *TelegramHandler) Test(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.SendTest(c.Request.Context(), id); err != nil {
		telegramError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
