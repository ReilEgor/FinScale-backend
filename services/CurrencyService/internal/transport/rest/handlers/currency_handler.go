package handlers

import (
	"net/http"
	"time"

	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/domain"
	"github.com/gin-gonic/gin"
)

func (h *Handler) ConvertCurrency(c *gin.Context) {
	ctx := c.Request.Context()
	var req domain.ConvertCurrencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	converted, err := h.uc.ConvertCurrency(ctx, req.From, req.To, req.Amount)
	if err != nil {
		h.logger.Error("conversion error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert currency"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"from":   req.From,
		"to":     req.To,
		"amount": req.Amount,
		"result": converted,
		"at":     time.Now().Format(time.RFC3339),
	})
	return
}
