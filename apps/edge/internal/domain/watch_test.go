package domain

import (
	"errors"
	"testing"
)

func TestNewWatchProgress_Success(t *testing.T) {
	wp, err := NewWatchProgress("media-123", 450, false)
	if err != nil {
		t.Fatalf("unexpected error creating watch progress: %v", err)
	}

	if wp.MediaID != "media-123" {
		t.Errorf("expected mediaID 'media-123', got %q", wp.MediaID)
	}
	if wp.PositionS != 450 {
		t.Errorf("expected position 450, got %d", wp.PositionS)
	}
	if wp.Completed {
		t.Error("expected completed to be false")
	}
	if wp.LastWatchedAt.IsZero() {
		t.Error("expected non-zero last_watched_at")
	}
}

func TestNewWatchProgress_ValidationErrors(t *testing.T) {
	tests := []struct {
		name        string
		mediaID     string
		position    int
		completed   bool
		expectedErr error
	}{
		{"empty mediaID", "", 10, false, ErrInvalidMediaID},
		{"whitespace mediaID", "   ", 10, false, ErrInvalidMediaID},
		{"negative position", "media-123", -1, false, ErrNegativePosition},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewWatchProgress(tt.mediaID, tt.position, tt.completed)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}
