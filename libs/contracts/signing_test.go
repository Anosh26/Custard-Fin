package contracts

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestSigningAndVerification_FullLifecycle(t *testing.T) {
	secret := []byte("super-secret-hmac-key")
	nodeID := "node-123"
	mediaID := "media-456"
	chunkPath := "1080p/segment_0001.ts"
	now := time.Now()
	expiresAt := now.Add(5 * time.Minute)

	query := SignChunkURL(secret, nodeID, mediaID, chunkPath, expiresAt)
	parts := strings.Split(query, "&")
	if len(parts) != 2 {
		t.Fatalf("unexpected query format: %s", query)
	}

	sig := strings.TrimPrefix(parts[0], "sig=")
	exp := strings.TrimPrefix(parts[1], "exp=")

	// 1. Valid verification
	if err := VerifyChunkSignature(secret, nodeID, mediaID, chunkPath, sig, exp, now); err != nil {
		t.Fatalf("expected signature to be valid: %v", err)
	}

	// 2. Tampered media ID
	if err := VerifyChunkSignature(secret, nodeID, "tampered-media", chunkPath, sig, exp, now); !errors.Is(err, ErrInvalidSignature) {
		t.Errorf("expected ErrInvalidSignature on tampered media ID, got: %v", err)
	}

	// 3. Tampered node ID
	if err := VerifyChunkSignature(secret, "tampered-node", mediaID, chunkPath, sig, exp, now); !errors.Is(err, ErrInvalidSignature) {
		t.Errorf("expected ErrInvalidSignature on tampered node ID, got: %v", err)
	}

	// 4. Tampered chunk path
	if err := VerifyChunkSignature(secret, nodeID, mediaID, "1080p/other.ts", sig, exp, now); !errors.Is(err, ErrInvalidSignature) {
		t.Errorf("expected ErrInvalidSignature on tampered chunk path, got: %v", err)
	}

	// 5. Wrong secret key
	if err := VerifyChunkSignature([]byte("wrong-key"), nodeID, mediaID, chunkPath, sig, exp, now); !errors.Is(err, ErrInvalidSignature) {
		t.Errorf("expected ErrInvalidSignature on wrong secret key, got: %v", err)
	}

	// 6. Expired signature
	futureNow := now.Add(10 * time.Minute)
	if err := VerifyChunkSignature(secret, nodeID, mediaID, chunkPath, sig, exp, futureNow); !errors.Is(err, ErrSignatureExpired) {
		t.Errorf("expected ErrSignatureExpired, got: %v", err)
	}

	// 7. Malformed exp parameter (non-integer)
	if err := VerifyChunkSignature(secret, nodeID, mediaID, chunkPath, sig, "not-a-number", now); !errors.Is(err, ErrMalformedURL) {
		t.Errorf("expected ErrMalformedURL on non-integer exp, got: %v", err)
	}
}
