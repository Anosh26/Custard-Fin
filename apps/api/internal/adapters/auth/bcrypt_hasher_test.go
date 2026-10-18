package auth_test

import (
	"errors"
	"testing"

	"github.com/anosh26/custardfin/apps/api/internal/adapters/auth"
	"github.com/anosh26/custardfin/apps/api/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func TestBcryptHasher_HashAndCompare(t *testing.T) {
	hasher := auth.NewBcryptHasher(bcrypt.MinCost) // MinCost for fast tests

	password := "correct_horse_battery_staple"
	hash, err := hasher.HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hash == "" || hash == password {
		t.Errorf("expected non-empty, non-plaintext hash")
	}

	// Verify correct password matches
	if !hasher.ComparePassword(password, hash) {
		t.Errorf("expected ComparePassword to succeed for correct password")
	}

	// Verify incorrect password fails
	if hasher.ComparePassword("wrong_password_123", hash) {
		t.Errorf("expected ComparePassword to fail for incorrect password")
	}
}

func TestBcryptHasher_Validation_MinLength(t *testing.T) {
	hasher := auth.NewBcryptHasher(bcrypt.MinCost)

	tests := []struct {
		name     string
		password string
	}{
		{"empty", ""},
		{"short 3 chars", "abc"},
		{"short 7 chars", "1234567"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := hasher.HashPassword(tt.password)
			if err == nil {
				t.Fatalf("expected validation error for short password")
			}
			if !errors.Is(err, domain.ErrValidation) {
				t.Errorf("expected ErrValidation, got %v", err)
			}
		})
	}
}
