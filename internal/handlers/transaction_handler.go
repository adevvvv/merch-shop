// internal/handlers/transaction_handler.go
package handlers

import (
	"net/http"
	"shop/internal/models"
	"shop/internal/services"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service services.TransactionService
}

func NewTransactionHandler(service services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) SendCoin(c *gin.Context) {
	var req models.SendCoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Errors: "invalid request"})
		return
	}

	userID, _ := c.Get("userID")
	if err := h.service.SendCoin(c.Request.Context(), userID.(int), req); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Errors: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "coins sent successfully"})
}