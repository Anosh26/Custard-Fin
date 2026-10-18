package domain_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/anosh26/custardfin/apps/api/internal/domain"
)

func TestNewUser_Success(t *testing.T) {
	username := "alice_123"
	passwordHash := "$2a$10$e8wV46L9z7eX.p.w/s.X.u3YjF9V7fVb0v/p4/yX1E0J7/j3YjF9V"

	user, err := domain.NewUser(username, passwordHash)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID == "" {
		t.Errorf("expected non-empty user ID")
	}
	if user.Username != username {
		t.Errorf("expected username %s, got %s", username, user.Username)
	}
	if user.PasswordHash != passwordHash {
		t.Errorf("expected password hash %s, got %s", passwordHash, user.PasswordHash)
	}
	if user.CreatedAt.IsZero() {
		t.Errorf("expected non-zero CreatedAt")
	}
	if user.UpdatedAt.IsZero() {
		t.Errorf("expected non-zero UpdatedAt")
	}
}

func TestNewUser_ValidationErrors(t *testing.T) {
	tests := []struct {
		name         string
		username     string
		passwordHash string
		wantErr      error
	}{
		{
			name:         "empty username",
			username:     "",
			passwordHash: "hash123",
			wantErr:      domain.ErrValidation,
		},
		{
			name:         "username too short (< 3 chars)",
			username:     "al",
			passwordHash: "hash123",
			wantErr:      domain.ErrValidation,
		},
		{
			name:         "username too long (> 64 chars)",
			username:     strings.Repeat("a", 65),
			passwordHash: "hash123",
			wantErr:      domain.ErrValidation,
		},
		{
			name:         "username with invalid characters",
			username:     "alice@example.com",
			passwordHash: "hash123",
			wantErr:      domain.ErrValidation,
		},
		{
			name:         "username with spaces",
			username:     "alice smith",
			passwordHash: "hash123",
			wantErr:      domain.ErrValidation,
		},
		{
			name:         "empty password hash",
			username:     "alice",
			passwordHash: "",
			wantErr:      domain.ErrValidation,
		},
		{
			name:         "whitespace password hash",
			username:     "alice",
			passwordHash: "   ",
			wantErr:      domain.ErrValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := domain.NewUser(tt.username, tt.passwordHash)
			if err == nil {
				t.Fatalf("expected error, got nil user: %+v", user)
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected error wrapping %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestNewUUIDv4_Uniqueness(t *testing.T) {
	seen := make(map[string]struct{})
	for i := 0; i < 50; i++ {
		u, err := domain.NewUser(fmt.Sprintf("user_%03d", i), "somehash")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, exists := seen[u.ID]; exists {
			t.Fatalf("collision detected for UUID: %s", u.ID)
		}
		seen[u.ID] = struct{}{}
		if len(u.ID) != 36 {
			t.Errorf("expected 36-char UUID, got %d (%s)", len(u.ID), u.ID)
		}
	}
}
