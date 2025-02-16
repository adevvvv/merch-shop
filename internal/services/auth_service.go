package services

import (
	"context"
	"database/sql"
	"errors"
	"shop/internal/models"
	"shop/internal/repositories"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AuthService struct {
	userRepo    repositories.UserRepository
	sessionRepo repositories.SessionRepository
	jwtSecret   []byte
	useSession  bool
}

func NewAuthService(
	userRepo repositories.UserRepository,
	sessionRepo repositories.SessionRepository,
	jwtSecret []byte,
	useSession bool,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtSecret:   jwtSecret,
		useSession:  useSession,
	}
}

func (s *AuthService) Authenticate(ctx context.Context, req models.AuthRequest) (models.AuthResponse, error) {
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if errors.Is(err, sql.ErrNoRows) {
		newUser := &models.User{
			Username: req.Username,
			Password: req.Password,
			Coins:    1000,
			IsAdmin:  false,
		}
		if err := s.userRepo.Create(ctx, newUser); err != nil {
			return models.AuthResponse{}, errors.New("user creation failed")
		}
		user = newUser
	} else if err != nil {
		return models.AuthResponse{}, errors.New("authentication error")
	}

	token, err := s.generateJWT(user)
	if err != nil {
		return models.AuthResponse{}, err
	}

	if s.useSession {
		if err := s.sessionRepo.Create(ctx, user.ID, token, time.Now().Add(time.Hour)); err != nil {
			return models.AuthResponse{}, errors.New("session creation failed")
		}
	}

	return models.AuthResponse{Token: token}, nil
}

func (s *AuthService) generateJWT(user *models.User) (string, error) {
	claims := &models.Claims{
		UserID:   user.ID,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}