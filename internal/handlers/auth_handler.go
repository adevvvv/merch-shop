package handlers

import (
	"net/http"
	"shop/internal/models"
	"shop/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) AuthHandler(c *gin.Context) {
	var req models.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Errors: "Invalid request"})
		return
	}

	resp, err := h.authService.Authenticate(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Errors: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}