package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/domain"
	"github.com/gin-gonic/gin"
)

// ConvertCurrencyResponse defines the successful output of a conversion operation.
type ConvertCurrencyResponse struct {
	From   string    `json:"from" example:"USD"`                // Source currency code (ISO 4217 format).
	To     string    `json:"to" example:"EUR"`                  // Target currency code (ISO 4217 format).
	Amount float64   `json:"amount" example:"100.0"`            // Initial value provided by the user for conversion.
	Result float64   `json:"result" example:"92.5"`             // Calculated value in the target currency based on the latest exchange rate.
	At     time.Time `json:"at" example:"2026-03-22T15:04:05Z"` // Timestamp when the conversion was processed (UTC).
}

// ConvertCurrencyRequest defines the input payload for currency conversion.
type ConvertCurrencyRequest struct {
	From   string  `json:"from" binding:"required,alpha" example:"USD"`    // Base currency code (ISO 4217)
	To     string  `json:"to" binding:"required,alpha" example:"EUR"`      // Target currency code (ISO 4217)
	Amount float64 `json:"amount" binding:"required,gt=0" example:"100.0"` // Amount to convert
}

// ConvertCurrency calculates the value of one currency in terms of another.
//
// @Summary      Convert Currency Amount
// @Description  Performs a real-time currency conversion using cached rates or external providers (CryptoCompare).
// @Tags         currency
// @Accept       json
// @Produce      json
// @Security     OAuth2AccessCode
// @Param        request  body      ConvertCurrencyRequest  true  "Conversion parameters (from, to, amount)"
// @Success      200      {object}  ConvertCurrencyResponse "Returns the converted amount and processing timestamp"
// @Failure      400      {object}  map[string]string       "Validation error: invalid currency codes or negative amount"
// @Failure      404      {object}  map[string]string       "Rate not found: the requested currency pair is unsupported"
// @Failure      502      {object}  map[string]string       "Bad Gateway: upstream exchange rate provider failed"
// @Failure      500      {object}  map[string]string       "Internal Server Error: something went wrong on our end"
// @Router       /currency/convert [post]
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
