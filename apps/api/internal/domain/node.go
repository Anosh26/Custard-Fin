package domain

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// Node represents a registered Home Edge Node.
type Node struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	Name          string    `json:"name"`
	NodeTokenHash string    `json:"-"`
	FrpAuthToken  string    `json:"-"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// NewNode constructs and validates a new Node entity.
func NewNode(userID, name, nodeTokenHash, frpAuthToken string) (*Node, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, fmt.Errorf("%w: user ID cannot be empty", ErrValidation)
	}

	name = strings.TrimSpace(name)
	if name == "" || len(name) > 128 {
		return nil, fmt.Errorf("%w: node name must be between 1 and 128 characters", ErrValidation)
	}

	nodeTokenHash = strings.TrimSpace(nodeTokenHash)
	if nodeTokenHash == "" {
		return nil, fmt.Errorf("%w: node token hash cannot be empty", ErrValidation)
	}

	frpAuthToken = strings.TrimSpace(frpAuthToken)
	if frpAuthToken == "" {
		return nil, fmt.Errorf("%w: frp auth token cannot be empty", ErrValidation)
	}

	now := time.Now().UTC()
	return &Node{
		ID:            newUUIDv4(),
		UserID:        userID,
		Name:          name,
		NodeTokenHash: nodeTokenHash,
		FrpAuthToken:  frpAuthToken,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

// GenerateNodeCredentials creates a cryptographically secure 32-byte node token,
// its SHA-256 hash (for database storage), and an FRP tunnel auth token.
func GenerateNodeCredentials() (rawToken, tokenHash, frpToken string, err error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", "", "", fmt.Errorf("failed to generate random token: %w", err)
	}
	rawToken = hex.EncodeToString(tokenBytes)
	tokenHash = HashNodeToken(rawToken)

	frpBytes := make([]byte, 16)
	if _, err := rand.Read(frpBytes); err != nil {
		return "", "", "", fmt.Errorf("failed to generate frp token: %w", err)
	}
	frpToken = hex.EncodeToString(frpBytes)

	return rawToken, tokenHash, frpToken, nil
}

// HashNodeToken computes the SHA-256 hex digest of a raw node token.
func HashNodeToken(rawToken string) string {
	h := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(h[:])
}
