package models

import (
	"time"
	"github.com/golang-jwt/jwt/v4"
)

type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Errors string `json:"errors"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Coins    int    `json:"coins"`
	IsAdmin  bool   `json:"is_admin"`
}

type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type CoinTransactionReceived struct {
	FromUser string    `json:"fromUser"`
	Amount   int       `json:"amount"`
	Date     time.Time `json:"date"`
}

type CoinTransactionSent struct {
	ToUser string    `json:"toUser"`
	Amount int       `json:"amount"`
	Date   time.Time `json:"date"`
}

type InfoResponse struct {
	Coins       int                      `json:"coins"`
	Inventory   []InventoryItem          `json:"inventory"`
	CoinHistory struct {
		Received []CoinTransactionReceived `json:"received"`
		Sent     []CoinTransactionSent     `json:"sent"`
	} `json:"coinHistory"`
}

type SendCoinRequest struct {
	ToUser string `json:"toUser" binding:"required"`
	Amount int    `json:"amount" binding:"required"`
}

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

var MerchItems = map[string]int{
	"t-shirt":    80,
	"cup":        20,
	"book":       50,
	"pen":        10,
	"powerbank": 200,
	"hoody":      300,
	"umbrella":   200,
	"socks":      10,
	"wallet":     50,
	"pink-hoody": 500,
}