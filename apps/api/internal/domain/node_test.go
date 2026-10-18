package domain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/anosh26/custardfin/apps/api/internal/domain"
)

func TestNewNode_Success(t *testing.T) {
	userID := "00000000-0000-0000-0000-000000000001"
	name := "Living Room RTX Server"
	tokenHash := "a1b2c3d4e5f60718293a4b5c6d7e8f90a1b2c3d4e5f60718293a4b5c6d7e8f90"
	frpToken := "frp_secret_key_123"

	node, err := domain.NewNode(userID, name, tokenHash, frpToken)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if node.ID == "" {
		t.Errorf("expected non-empty node ID")
	}
	if node.UserID != userID {
		t.Errorf("expected userID %s, got %s", userID, node.UserID)
	}
	if node.Name != name {
		t.Errorf("expected name %s, got %s", name, node.Name)
	}
	if node.NodeTokenHash != tokenHash {
		t.Errorf("expected tokenHash %s, got %s", tokenHash, node.NodeTokenHash)
	}
	if node.FrpAuthToken != frpToken {
		t.Errorf("expected frpToken %s, got %s", frpToken, node.FrpAuthToken)
	}
	if node.CreatedAt.IsZero() {
		t.Errorf("expected non-zero CreatedAt")
	}
	if node.UpdatedAt.IsZero() {
		t.Errorf("expected non-zero UpdatedAt")
	}
}

func TestNewNode_ValidationErrors(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		nodeName      string
		nodeTokenHash string
		frpAuthToken  string
		wantErr       error
	}{
		{
			name:          "empty userID",
			userID:        "",
			nodeName:      "Node 1",
			nodeTokenHash: "hash",
			frpAuthToken:  "frp",
			wantErr:       domain.ErrValidation,
		},
		{
			name:          "whitespace userID",
			userID:        "   ",
			nodeName:      "Node 1",
			nodeTokenHash: "hash",
			frpAuthToken:  "frp",
			wantErr:       domain.ErrValidation,
		},
		{
			name:          "empty node name",
			userID:        "user-1",
			nodeName:      "",
			nodeTokenHash: "hash",
			frpAuthToken:  "frp",
			wantErr:       domain.ErrValidation,
		},
		{
			name:          "whitespace node name",
			userID:        "user-1",
			nodeName:      "   ",
			nodeTokenHash: "hash",
			frpAuthToken:  "frp",
			wantErr:       domain.ErrValidation,
		},
		{
			name:          "node name too long (> 128 chars)",
			userID:        "user-1",
			nodeName:      strings.Repeat("n", 129),
			nodeTokenHash: "hash",
			frpAuthToken:  "frp",
			wantErr:       domain.ErrValidation,
		},
		{
			name:          "empty nodeTokenHash",
			userID:        "user-1",
			nodeName:      "Node 1",
			nodeTokenHash: "",
			frpAuthToken:  "frp",
			wantErr:       domain.ErrValidation,
		},
		{
			name:          "empty frpAuthToken",
			userID:        "user-1",
			nodeName:      "Node 1",
			nodeTokenHash: "hash",
			frpAuthToken:  "",
			wantErr:       domain.ErrValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := domain.NewNode(tt.userID, tt.nodeName, tt.nodeTokenHash, tt.frpAuthToken)
			if err == nil {
				t.Fatalf("expected error, got nil node: %+v", node)
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected error wrapping %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestGenerateNodeCredentials(t *testing.T) {
	rawToken, tokenHash, frpToken, err := domain.GenerateNodeCredentials()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(rawToken) != 64 {
		t.Errorf("expected 64-char hex raw token, got %d", len(rawToken))
	}
	if len(tokenHash) != 64 {
		t.Errorf("expected 64-char hex sha256 token hash, got %d", len(tokenHash))
	}
	if len(frpToken) != 32 {
		t.Errorf("expected 32-char hex frp token, got %d", len(frpToken))
	}

	expectedHash := domain.HashNodeToken(rawToken)
	if tokenHash != expectedHash {
		t.Errorf("expected hash %s, got %s", expectedHash, tokenHash)
	}
}
