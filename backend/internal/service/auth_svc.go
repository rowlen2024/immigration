package service

import (
	"errors"
	"fmt"
	"time"

	"mygo-immigration/backend/internal/config"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication and token management.
type AuthService struct {
	repo repository.UserRepository
	cfg  *config.Config
}

// NewAuthService creates a new AuthService with the given repository and config.
func NewAuthService(repo repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{repo: repo, cfg: cfg}
}

// TokenPair represents an access and refresh token pair returned on login.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// Login verifies the username and password, then returns JWT access and refresh tokens.
func (s *AuthService) Login(username, password string) (*TokenPair, error) {
	if username == "" || password == "" {
		return nil, errors.New("username and password are required")
	}

	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("invalid username or password: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	if user.Status != 1 {
		return nil, errors.New("account is disabled")
	}

	accessToken, err := s.generateToken(user.ID, user.Username, user.Role, s.cfg.JWTAccessExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateToken(user.ID, user.Username, user.Role, s.cfg.JWTRefreshExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.cfg.JWTAccessExpiry.Seconds()),
	}, nil
}

// RefreshToken validates the refresh token and generates a new access token.
func (s *AuthService) RefreshToken(refreshToken string) (*TokenPair, error) {
	if refreshToken == "" {
		return nil, errors.New("refresh token is required")
	}

	claims := &model.JWTClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired refresh token")
	}

	newAccessToken, err := s.generateToken(claims.UserID, claims.Username, claims.Role, s.cfg.JWTAccessExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &TokenPair{
		AccessToken: newAccessToken,
		ExpiresIn:   int64(s.cfg.JWTAccessExpiry.Seconds()),
	}, nil
}

func (s *AuthService) generateToken(userID uint64, username, role string, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := &model.JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}
