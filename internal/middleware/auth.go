package middleware

import (
	"context"
	"shop/internal/repositories"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type AuthMiddleware struct {
	sessionRepo repositories.SessionRepository
	jwtSecret   []byte
	useSession  bool
}

func NewAuthMiddleware(sessionRepo repositories.SessionRepository, jwtSecret []byte, useSession bool) *AuthMiddleware {
	return &AuthMiddleware{
		sessionRepo: sessionRepo,
		jwtSecret:   jwtSecret,
		useSession:  useSession,
	}
}

func (m *AuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Missing token"})
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token format"})
			return
		}

		tokenStr := tokenParts[1]
		claims, err := m.validateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		if m.useSession {
			if valid, _ := m.validateSession(c.Request.Context(), claims.UserID, tokenStr); !valid {
				c.AbortWithStatusJSON(401, gin.H{"error": "Invalid session"})
				return
			}
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("isAdmin", claims.IsAdmin)
		c.Next()
	}
}

func (m *AuthMiddleware) validateToken(tokenStr string) (*models.Claims, error) {
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return m.jwtSecret, nil
	})
	
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}

func (m *AuthMiddleware) validateSession(ctx context.Context, userID int, tokenStr string) (bool, error) {
	storedToken, err := m.sessionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return false, err
	}
	return storedToken == tokenStr, nil
}