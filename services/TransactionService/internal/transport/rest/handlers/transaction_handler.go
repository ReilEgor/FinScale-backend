package handlers

import (
	"encoding/json"
	"github.com/ReilEgor/FinScale-backend/TransactionService/internal/domain"
	sharedContextUtil "github.com/ReilEgor/FinScale-shared/pkg/contextutil"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	formFieldFile = "receipt"
	formFieldData = "data"
)

// RecordTransaction godoc
// @Summary      Record a new transaction with receipt
// @Description  Uploads a receipt image and saves transaction details.
// @Description  The 'data' field must contain a JSON string representing the transaction object.
// @Tags         transactions
// @Accept       mpfd
// @Produce      json
// @Param        X-User-ID  header    string  true  "User ID"
// @Param        receipt    formData  file    true  "Receipt image file"
// @Param        data       formData  string  true  "Transaction JSON data (e.g. {\"amount\": 100, \"currency\": \"USD\"})"
// @Success      201        {object}  map[string]interface{} "status: success, receipt_url: string"
// @Failure      400        {object}  map[string]interface{} "error: invalid data format or missing file"
// @Failure      401        {object}  map[string]interface{} "error: unauthorized"
// @Failure      500        {object}  map[string]interface{} "error: internal server error"
// @Router       /transactions [post]
func (h *Handler) RecordTransaction(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := sharedContextUtil.UserIDFromContext(ctx)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	file, err := c.FormFile(formFieldFile)
	if err != nil {
		h.logger.Error("receipt file is required", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "receipt file is required"})
		return
	}

	var req domain.Transaction
	if err := json.Unmarshal([]byte(c.PostForm(formFieldData)), &req); err != nil {
		h.logger.Error("invalid transaction data", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data format"})
		return
	}

	fileContent, err := file.Open()
	if err != nil {
		h.logger.Error("failed to open receipt file", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open receipt file"})
		return
	}
	defer fileContent.Close()

	receiptURL, err := h.uc.RecordReceipt(ctx, fileContent, file.Filename)
	if err != nil {
		h.logger.Error("failed to upload receipt", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload receipt"})
		return
	}

	req.UserID = userID
	req.ReceiptURL = receiptURL

	transactionID, err := h.uc.RecordTransaction(ctx, req)
	if err != nil {
		h.logger.Error("failed to record transaction",
			"user_id", userID,
			"error", err,
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to record transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":         "success",
		"receipt_url":    receiptURL,
		"transaction_id": transactionID.String(),
	})
}
