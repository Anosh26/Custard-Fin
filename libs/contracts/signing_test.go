package contracts

import (
	"strings"
	"testing"
	"time"
)

func TestSigningAndVerification(t *testing.T) {
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

	// Valid verification
	if err := VerifyChunkSignature(secret, nodeID, mediaID, chunkPath, sig, exp, now); err != nil {
		t.Fatalf("expected signature to be valid: %v", err)
	}

	// Tampered media ID
	if err := VerifyChunkSignature(secret, nodeID, "tampered-media", chunkPath, sig, exp, now); err == nil {
		t.Fatalf("expected signature verification failure on tampered media ID")
	}

	// Expired signature
	futureNow := now.Add(10 * time.Minute)
	if err := VerifyChunkSignature(secret, nodeID, mediaID, chunkPath, sig, exp, futureNow); err != ErrSignatureExpired {
		t.Fatalf("expected ErrSignatureExpired, got: %v", err)
	}
}
