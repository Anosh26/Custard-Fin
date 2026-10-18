package auth

import (
	"fmt"
	"time"

	"github.com/anosh26/custardfin/apps/api/internal/domain"
	"github.com/anosh26/custardfin/apps/api/internal/ports"
	"github.com/golang-jwt/jwt/v5"
)

type customClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secret []byte
}

// NewJWTService creates a new HMAC-SHA256 JWT service with the given secret.
func NewJWTService(secret []byte) (ports.TokenService, error) {
	if len(secret) < 16 {
		return nil, fmt.Errorf("%w: jwt secret must be at least 16 bytes", domain.ErrValidation)
	}
	return &JWTService{secret: secret}, nil
}

func (s *JWTService) GenerateToken(userID, username string, ttl time.Duration) (string, error) {
	if userID == "" || username == "" {
		return "", fmt.Errorf("%w: userID and username are required for token generation", domain.ErrValidation)
	}
	if ttl <= 0 {
		return "", fmt.Errorf("%w: token ttl must be positive", domain.ErrValidation)
	}

	now := time.Now().UTC()
	claims := customClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			Issuer:    "custardfin-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}

func (s *JWTService) ValidateToken(tokenStr string) (*ports.Claims, error) {
	if tokenStr == "" {
		return nil, fmt.Errorf("%w: empty token", domain.ErrUnauthorized)
	}

	token, err := jwt.ParseWithClaims(tokenStr, &customClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("%w: invalid or expired token: %v", domain.ErrUnauthorized, err)
	}

	claims, ok := token.Claims.(*customClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("%w: invalid token claims", domain.ErrUnauthorized)
	}

	var exp, iat int64
	if claims.ExpiresAt != nil {
		exp = claims.ExpiresAt.Unix()
	}
	if claims.IssuedAt != nil {
		iat = claims.IssuedAt.Unix()
	}

	return &ports.Claims{
		UserID:    claims.UserID,
		Username:  claims.Username,
		ExpiresAt: exp,
		IssuedAt:  iat,
	}, nil
}
