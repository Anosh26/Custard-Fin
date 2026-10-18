package domain

import (
	"strings"
	"testing"
	"time"
)

func TestNewMedia(t *testing.T) {
	now := time.Now()
	m, err := NewMedia("/videos/tears_of_steel.mp4", 1048576, now, "dummyhash", 120, 1920, 1080, "h264", "aac")
	if err != nil {
		t.Fatalf("unexpected error creating media: %v", err)
	}

	if m.Title != "tears_of_steel" {
		t.Errorf("expected title 'tears_of_steel', got %q", m.Title)
	}
	if len(m.ID) != 36 {
		t.Errorf("expected 36-char UUID, got %q", m.ID)
	}
	if !strings.Contains(m.ID, "-") {
		t.Errorf("expected formatted UUID, got %q", m.ID)
	}

	// Validation checks
	if _, err := NewMedia("", 100, now, "hash", 100, 1920, 1080, "h264", "aac"); err != ErrEmptyPath {
		t.Errorf("expected ErrEmptyPath, got: %v", err)
	}
	if _, err := NewMedia("test.mp4", 100, now, "", 100, 1920, 1080, "h264", "aac"); err != ErrEmptyContentHash {
		t.Errorf("expected ErrEmptyContentHash, got: %v", err)
	}
	if _, err := NewMedia("test.mp4", 100, now, "hash", 0, 1920, 1080, "h264", "aac"); err != ErrInvalidDuration {
		t.Errorf("expected ErrInvalidDuration, got: %v", err)
	}
}
