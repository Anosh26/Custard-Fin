package ports

import (
	"time"
)

// Claims represents claims stored in an authenticated JWT token.
type Claims struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
}

// PasswordHasher handles secure password hashing and verification.
type PasswordHasher interface {
	HashPassword(password string) (string, error)
	ComparePassword(password, hash string) bool
}

// TokenService issues and verifies JWT bearer tokens.
type TokenService interface {
	GenerateToken(userID, username string, ttl time.Duration) (string, error)
	ValidateToken(tokenStr string) (*Claims, error)
}
