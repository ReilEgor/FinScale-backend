package rest

import (
	"log/slog"

	handler "github.com/ReilEgor/FinScale-backend/CurrencyService/internal/transport/rest/handlers"
	usecase "github.com/ReilEgor/FinScale-backend/CurrencyService/internal/usecase"
	"github.com/gin-gonic/gin"
)

type GinServer struct {
	router *gin.Engine
	uc     *usecase.CurrencyUseCase
	logger *slog.Logger
}

func NewGinServer(uc *usecase.CurrencyUseCase) *GinServer {
	router := gin.New()
	logger := slog.With(slog.String("component", "gin_server"))

	SetupMiddleware(router, logger)

	s := &GinServer{
		router: router,
		uc:     uc,
		logger: logger,
	}

	h := handler.NewHandler(uc)
	h.InitRoutes(s.router)

	return s
}

func (s *GinServer) Run(port string) error {
	return s.router.Run(port)
}
