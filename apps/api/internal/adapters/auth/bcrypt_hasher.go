package auth

import (
	"errors"
	"fmt"

	"github.com/anosh26/custardfin/apps/api/internal/domain"
	"github.com/anosh26/custardfin/apps/api/internal/ports"
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
	cost int
}

// NewBcryptHasher creates a new BcryptHasher with specified cost (defaulting to bcrypt.DefaultCost if <= 0).
func NewBcryptHasher(cost int) ports.PasswordHasher {
	if cost <= 0 {
		cost = bcrypt.DefaultCost
	}
	return &BcryptHasher{cost: cost}
}

func (h *BcryptHasher) HashPassword(password string) (string, error) {
	if len(password) < 8 {
		return "", fmt.Errorf("%w: password must be at least 8 characters", domain.ErrValidation)
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

func (h *BcryptHasher) ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil || errors.Is(err, nil)
}
