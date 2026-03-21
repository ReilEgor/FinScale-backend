package handlers

import (
	"log/slog"

	"github.com/ReilEgor/FinScale-backend/TransactionService/internal/domain"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc     domain.TransactionUseCase
	logger *slog.Logger
}

func NewHandler(uc domain.TransactionUseCase) *Handler {
	return &Handler{
		uc:     uc,
		logger: slog.With(slog.String("component", "handler")),
	}
}

func (h *Handler) InitRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		currency := api.Group("/transaction")
		{
			currency.POST("/record", h.RecordTransaction)
		}
	}
}
