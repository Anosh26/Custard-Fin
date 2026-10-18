package auth_test

import (
	"errors"
	"testing"
	"time"

	"github.com/anosh26/custardfin/apps/api/internal/adapters/auth"
	"github.com/anosh26/custardfin/apps/api/internal/domain"
)

func TestNewJWTService_SecretLength(t *testing.T) {
	_, err := auth.NewJWTService([]byte("short"))
	if err == nil {
		t.Fatalf("expected error for secret < 16 bytes")
	}
	if !errors.Is(err, domain.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}

	svc, err := auth.NewJWTService([]byte("1234567890123456"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc == nil {
		t.Fatalf("expected non-nil service")
	}
}

func TestJWTService_GenerateAndValidate_Success(t *testing.T) {
	secret := []byte("super-secret-key-that-is-at-least-32-bytes-long!")
	svc, err := auth.NewJWTService(secret)
	if err != nil {
		t.Fatalf("failed to create JWT service: %v", err)
	}

	token, err := svc.GenerateToken("user-123", "alice", time.Hour)
	if err != nil {
		t.Fatalf("expected no error generating token, got %v", err)
	}

	if token == "" {
		t.Fatalf("expected non-empty token")
	}

	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("expected valid token, got %v", err)
	}

	if claims.UserID != "user-123" {
		t.Errorf("expected UserID user-123, got %s", claims.UserID)
	}
	if claims.Username != "alice" {
		t.Errorf("expected Username alice, got %s", claims.Username)
	}
	if claims.ExpiresAt <= claims.IssuedAt {
		t.Errorf("expected ExpiresAt > IssuedAt")
	}
}

func TestJWTService_ExpiredToken(t *testing.T) {
	secret := []byte("super-secret-key-that-is-at-least-32-bytes-long!")
	svc, err := auth.NewJWTService(secret)
	if err != nil {
		t.Fatalf("failed to create JWT service: %v", err)
	}

	// Generate token with negative duration so it is already expired
	token, err := svc.GenerateToken("user-123", "alice", -time.Second)
	// Actually GenerateToken validates ttl > 0, let's verify that validation first
	if err == nil {
		t.Fatalf("expected error for negative ttl")
	}
	if !errors.Is(err, domain.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}

	// Generate with 1 millisecond and wait 10 milliseconds
	token, err = svc.GenerateToken("user-123", "alice", time.Millisecond)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	time.Sleep(10 * time.Millisecond)

	_, err = svc.ValidateToken(token)
	if err == nil {
		t.Fatalf("expected error for expired token")
	}
	if !errors.Is(err, domain.ErrUnauthorized) {
		t.Errorf("expected ErrUnauthorized for expired token, got %v", err)
	}
}

func TestJWTService_TamperedToken(t *testing.T) {
	secret := []byte("super-secret-key-that-is-at-least-32-bytes-long!")
	svc, err := auth.NewJWTService(secret)
	if err != nil {
		t.Fatalf("failed to create JWT service: %v", err)
	}

	token, err := svc.GenerateToken("user-123", "alice", time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	// Tamper with secret by validating with different key
	svc2, _ := auth.NewJWTService([]byte("completely-different-signing-key-123456"))
	_, err = svc2.ValidateToken(token)
	if err == nil {
		t.Fatalf("expected error validating with wrong key")
	}
	if !errors.Is(err, domain.ErrUnauthorized) {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}

	// Empty token
	_, err = svc.ValidateToken("")
	if err == nil {
		t.Fatalf("expected error for empty token")
	}
	if !errors.Is(err, domain.ErrUnauthorized) {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}
