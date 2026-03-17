package rest

import (
	"github.com/ReilEgor/FinScale-backend/UserService/internal/domain"
	"github.com/ReilEgor/FinScale-backend/UserService/internal/transport/rest/handlers"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type GinServer struct {
	router *gin.Engine
	uc     domain.UserUseCase
	logger *slog.Logger
}

func NewGinServer(uc domain.UserUseCase) *GinServer {
	router := gin.New()
	logger := slog.With(slog.String("component", "gin_server"))

	SetupMiddleware(router, logger)

	s := &GinServer{
		router: router,
		uc:     uc,
		logger: logger,
	}

	h := handlers.NewHandler(uc)
	h.InitRoutes(s.router)

	return s
}

func (s *GinServer) Run(addr string) error {
	s.logger.Info("Starting Gin server", slog.String("address", addr))

	return s.router.Run(addr)
}
