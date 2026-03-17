package handlers

import (
	"log/slog"

	"github.com/ReilEgor/FinScale-backend/UserService/internal/domain"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc     domain.UserUseCase
	logger *slog.Logger
}

func NewHandler(uc domain.UserUseCase) *Handler {
	return &Handler{
		uc:     uc,
		logger: slog.With(slog.String("handler", "NewHandler")),
	}
}

func (h *Handler) InitRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.POST("", h.CreateUser)
			users.GET("/:id", h.GetUser)
			users.PUT("/:id", h.UpdateUser)
			users.DELETE("/:id", h.DeleteUser)
			users.POST("/sync", h.SyncUser)
		}
	}
}
