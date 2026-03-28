// Package jwt provides JWT token generation and validation services.
package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// Claims represents the JWT claims for access tokens.
type Claims struct {
	jwt.RegisteredClaims
	UserID          uuid.UUID `json:"user_id"`
	Email           string    `json:"email"`
	IsPlatformAdmin bool      `json:"is_platform_admin,omitempty"`
}

// Service handles JWT token operations.
type Service struct {
	secret          []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	issuer          string
}

// NewService creates a new JWT service.
func NewService(secret string, accessTTLSeconds, refreshTTLSeconds int) *Service {
	return &Service{
		secret:          []byte(secret),
		accessTokenTTL:  time.Duration(accessTTLSeconds) * time.Second,
		refreshTokenTTL: time.Duration(refreshTTLSeconds) * time.Second,
		issuer:          "goauth",
	}
}

// NewServiceWithIssuer creates a new JWT service with a custom issuer.
func NewServiceWithIssuer(secret string, accessTTLSeconds, refreshTTLSeconds int, issuer string) *Service {
	svc := NewService(secret, accessTTLSeconds, refreshTTLSeconds)
	svc.issuer = issuer
	return svc
}

// TokenPair represents an access and refresh token pair.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int // seconds until access token expires
}

// GenerateTokenPair generates a new access and refresh token pair.
func (s *Service) GenerateTokenPair(userID uuid.UUID, email string, isPlatformAdmin bool) (*TokenPair, error) {
	now := time.Now()
	accessExpiry := now.Add(s.accessTokenTTL)

	// Access token claims
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
		UserID:          userID,
		Email:           email,
		IsPlatformAdmin: isPlatformAdmin,
	}

	// Create access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err := accessToken.SignedString(s.secret)
	if err != nil {
		return nil, err
	}

	// Generate refresh token (opaque token)
	refreshToken, err := GenerateRandomToken(32)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken,
		ExpiresIn:    int(s.accessTokenTTL.Seconds()),
	}, nil
}

// GenerateAccessToken generates only an access token (useful for token refresh).
func (s *Service) GenerateAccessToken(userID uuid.UUID, email string, isPlatformAdmin bool) (string, int, error) {
	now := time.Now()
	accessExpiry := now.Add(s.accessTokenTTL)

	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
		UserID:          userID,
		Email:           email,
		IsPlatformAdmin: isPlatformAdmin,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err := accessToken.SignedString(s.secret)
	if err != nil {
		return "", 0, err
	}

	return accessTokenString, int(s.accessTokenTTL.Seconds()), nil
}

// ValidateAccessToken validates an access token and returns the claims.
func (s *Service) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// RefreshTokenTTL returns the refresh token TTL.
func (s *Service) RefreshTokenTTL() time.Duration {
	return s.refreshTokenTTL
}

// AccessTokenTTL returns the access token TTL.
func (s *Service) AccessTokenTTL() time.Duration {
	return s.accessTokenTTL
}

// GenerateRandomToken generates a cryptographically secure random token.
func GenerateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
