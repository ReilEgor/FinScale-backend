package middleware

import (
	sharedMiddleware "github.com/ReilEgor/FinScale-shared/middleware"
	"github.com/gin-gonic/gin"
	"log/slog"
	"time"
)

func SetupMiddleware(router *gin.Engine, logger *slog.Logger) {
	router.Use(gin.Recovery())
	router.Use(slogMiddleware(logger))
	router.Use(sharedMiddleware.AuthMiddleware())

}

func slogMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		logger.Info("request handled",
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("duration", time.Since(start)),
			slog.String("ip", c.ClientIP()),
		)
	}
}
