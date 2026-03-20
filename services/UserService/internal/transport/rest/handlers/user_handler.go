package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	KetCloakID string `json:"keycloak_id" binding:"required"`
}
type UpdateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

func (h *Handler) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.uc.CreateUser(ctx, req.KetCloakID, req.Name, req.Email)
	c.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
}

func (h *Handler) GetUser(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("missing user id").Error()})
		return
	}
	user, err := h.uc.GetUser(ctx, userID)
	if err != nil {
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("missing user id").Error()})
		return
	}
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.uc.UpdateUser(ctx, userID, req.Name, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("missing user id").Error()})
		return
	}
	atoi, err := strconv.Atoi(userID)
	if err != nil {
		return
	}
	err = h.uc.DeleteUser(ctx, int64(atoi))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

func (h *Handler) SyncUser(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetHeader("X-User-ID")
	email := c.GetHeader("X-User-Email")
	name := c.GetHeader("X-User-Name")

	if userID == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.uc.GetUser(ctx, userID)

	if user == nil {
		_, err = h.uc.CreateUser(ctx, userID, name, email)
	}
	if err != nil {
		h.logger.Error("failed to sync user", slog.Any("err", err))
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, user)
}
