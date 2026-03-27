package handlers

import (
	_ "github.com/ReilEgor/FinScale-backend/TransactionService/api/docs"
	"github.com/ReilEgor/FinScale-backend/TransactionService/internal/domain"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"log/slog"
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
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := router.Group("/api/v1")
	{
		currency := api.Group("/transaction")
		{
			currency.POST("/record", h.RecordTransaction)
		}
	}
}
