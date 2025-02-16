// internal/handlers/user_handler.go
package handlers

import (
	"net/http"
	"shop/internal/models"
	"shop/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetInfo(c *gin.Context) {
	userID, _ := c.Get("userID")
	info, err := h.service.GetUserInfo(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Errors: err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

func (h *UserHandler) BuyItem(c *gin.Context) {
	item := c.Param("item")
	userID, _ := c.Get("userID")
	
	if err := h.service.BuyItem(c.Request.Context(), userID.(int), item); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Errors: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "item purchased successfully"})
}