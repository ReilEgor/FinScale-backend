package handlers

import (
	"log/slog"
	"net/http"

	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/domain"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc     domain.CurrencyUseCase
	logger *slog.Logger
}

func NewHandler(uc domain.CurrencyUseCase) *Handler {
	return &Handler{
		uc:     uc,
		logger: slog.With(slog.String("component", "handler")),
	}
}

func (h *Handler) InitRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api/v1")
	{
		currency := api.Group("/currency")
		{
			currency.POST("/convert", h.ConvertCurrency)
		}
	}
}
