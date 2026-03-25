package rest

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	handler "github.com/ReilEgor/FinScale-backend/CurrencyService/internal/transport/rest/handlers"
	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/transport/rest/middleware"
	usecase "github.com/ReilEgor/FinScale-backend/CurrencyService/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type GinServer struct {
	router *gin.Engine
	uc     *usecase.CurrencyUseCase
	logger *slog.Logger
}

func NewGinServer(uc *usecase.CurrencyUseCase, redisClient *redis.Client) *GinServer {
	router := gin.New()
	logger := slog.With(slog.String("component", "gin_server"))
	middleware.SetupMiddleware(router, logger, redisClient)

	s := &GinServer{
		router: router,
		uc:     uc,
		logger: logger,
	}

	h := handler.NewHandler(uc)
	h.InitRoutes(s.router)

	return s
}

func (s *GinServer) Run(ctx context.Context, port string) error {
	srv := &http.Server{Addr: port, Handler: s.router}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			s.logger.Error("forced shutdown", "error", err)
		}
	}()

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
