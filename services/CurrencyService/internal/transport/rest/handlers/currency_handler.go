package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/domain"
	"github.com/gin-gonic/gin"
)

type ConvertCurrencyResponse struct {
	From   string    `json:"from"`
	To     string    `json:"to"`
	Amount float64   `json:"amount"`
	Result float64   `json:"result"`
	At     time.Time `json:"at"`
}
type ConvertCurrencyRequest struct {
	From   string  `json:"from" binding:"required,alpha"`
	To     string  `json:"to" binding:"required,alpha"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

func (h *Handler) ConvertCurrency(c *gin.Context) {
	ctx := c.Request.Context()
	var req ConvertCurrencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	converted, err := h.uc.ConvertCurrency(ctx, req.From, req.To, req.Amount)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRateNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "currency pair not found"})
		case errors.Is(err, domain.ErrFetchFromExternalAPI):
			c.JSON(http.StatusBadGateway, gin.H{"error": "external API unavailable"})
		default:
			h.logger.Error("failed to convert currency", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
	c.JSON(http.StatusOK, ConvertCurrencyResponse{
		From:   req.From,
		To:     req.To,
		Amount: req.Amount,
		Result: converted,
		At:     time.Now(),
	})
}
