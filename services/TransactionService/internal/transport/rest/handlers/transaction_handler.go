package handlers

import (
	"encoding/json"

	"github.com/ReilEgor/FinScale-backend/TransactionService/internal/domain"
	"github.com/gin-gonic/gin"
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

	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		h.logger.Error("user id is missing in headers")
		c.JSON(401, gin.H{"error": "unauthorized: user id missing"})
		return
	}

	file, err := c.FormFile("receipt")
	if err != nil {
		h.logger.Error("receipt is required", "error", err)
		c.JSON(400, gin.H{"error": "receipt file is required"})
		return
	}

	jsonData := c.PostForm("data")
	var req domain.Transaction
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		c.JSON(400, gin.H{"error": "invalid data format"})
		return
	}

	fileContent, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to open receipt file"})
		return
	}
	defer fileContent.Close()

	receiptURL, err := h.uc.RecordReceipt(ctx, fileContent, file.Filename)
	if err != nil {
		h.logger.Error("upload failed", "error", err)
		c.JSON(500, gin.H{"error": "failed to upload receipt"})
		return
	}

	req.ReceiptURL = receiptURL
	req.UserID = userID
	if err := h.uc.RecordTransaction(ctx, req); err != nil {
		h.logger.Error("db record failed", "error", err)
		c.JSON(500, gin.H{"error": "failed to record transaction"})
		return
	}

	c.JSON(201, gin.H{"status": "success", "receipt_url": receiptURL})
}
