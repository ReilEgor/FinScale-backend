package rest

import (
	"context"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupMiddleware(router *gin.Engine, logger *slog.Logger) {
	router.Use(gin.Recovery())
	router.Use(slogMiddleware(logger))
	router.Use(Timeout(5 * time.Second))
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
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
