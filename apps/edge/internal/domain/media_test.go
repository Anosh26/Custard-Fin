package domain

import (
	"errors"
	"regexp"
	"testing"
	"time"
)

var uuidRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func TestNewMedia_Success(t *testing.T) {
	modTime := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)
	m, err := NewMedia(
		"/videos/folder/Avatar.The.Way.Of.Water.2022.mkv",
		1572864000,
		modTime,
		"a3f5b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6",
		11520,
		3840,
		2160,
		"hevc",
		"eac3",
	)
	if err != nil {
		t.Fatalf("unexpected error creating media: %v", err)
	}

	if m.Title != "Avatar.The.Way.Of.Water.2022" {
		t.Errorf("expected title 'Avatar.The.Way.Of.Water.2022', got %q", m.Title)
	}
	if m.DurationS != 11520 {
		t.Errorf("expected duration 11520, got %d", m.DurationS)
	}
	if m.Width != 3840 || m.Height != 2160 {
		t.Errorf("expected resolution 3840x2160, got %dx%d", m.Width, m.Height)
	}
	if m.VideoCodec != "hevc" || m.AudioCodec != "eac3" {
		t.Errorf("expected codecs hevc/eac3, got %s/%s", m.VideoCodec, m.AudioCodec)
	}
	if !m.ModTime.Equal(modTime) {
		t.Errorf("expected modTime %v, got %v", modTime, m.ModTime)
	}

	// UUID validation (RFC4122 v4)
	if !uuidRegex.MatchString(m.ID) {
		t.Errorf("media ID is not a valid RFC4122 v4 UUID: %s", m.ID)
	}
}

func TestNewMedia_ValidationErrors(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name        string
		path        string
		hash        string
		duration    int
		expectedErr error
	}{
		{"empty path", "", "hash", 100, ErrEmptyPath},
		{"whitespace path", "   ", "hash", 100, ErrEmptyPath},
		{"empty hash", "/path/video.mp4", "", 100, ErrEmptyContentHash},
		{"whitespace hash", "/path/video.mp4", "   ", 100, ErrEmptyContentHash},
		{"zero duration", "/path/video.mp4", "hash", 0, ErrInvalidDuration},
		{"negative duration", "/path/video.mp4", "hash", -5, ErrInvalidDuration},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewMedia(tt.path, 100, now, tt.hash, tt.duration, 1920, 1080, "h264", "aac")
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestNewUUID_Uniqueness(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id, err := newUUID()
		if err != nil {
			t.Fatalf("newUUID returned error: %v", err)
		}
		if seen[id] {
			t.Fatalf("collision detected for UUID: %s", id)
		}
		seen[id] = true
		if !uuidRegex.MatchString(id) {
			t.Errorf("generated UUID %q failed RFC4122 v4 regex", id)
		}
	}
}
