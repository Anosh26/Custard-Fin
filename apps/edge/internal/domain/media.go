package domain

import (
	"crypto/rand"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

var (
	ErrEmptyPath         = errors.New("media path cannot be empty")
	ErrInvalidDuration   = errors.New("media duration must be positive")
	ErrEmptyContentHash  = errors.New("content hash is required")
)

// Media represents a discovered local video file and its metadata.
type Media struct {
	ID           string    `json:"id"`
	Path         string    `json:"path"`
	FileSize     int64     `json:"file_size"`
	ModTime      time.Time `json:"mod_time"`
	ContentHash  string    `json:"content_hash"` // SHA-256
	Title        string    `json:"title"`
	DurationS    int       `json:"duration_s"`
	Width        int       `json:"width"`
	Height       int       `json:"height"`
	VideoCodec   string    `json:"video_codec"`
	AudioCodec   string    `json:"audio_codec"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewMedia constructs a new Media entity with a generated UUID.
func NewMedia(filePath string, fileSize int64, modTime time.Time, hash string, durationS int, width, height int, videoCodec, audioCodec string) (*Media, error) {
	if strings.TrimSpace(filePath) == "" {
		return nil, ErrEmptyPath
	}
	if strings.TrimSpace(hash) == "" {
		return nil, ErrEmptyContentHash
	}
	if durationS <= 0 {
		return nil, ErrInvalidDuration
	}

	cleanPath := filepath.Clean(filePath)
	base := filepath.Base(cleanPath)
	ext := filepath.Ext(base)
	title := strings.TrimSuffix(base, ext)

	id, err := newUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate media id: %w", err)
	}

	now := time.Now().UTC()
	return &Media{
		ID:          id,
		Path:        cleanPath,
		FileSize:    fileSize,
		ModTime:     modTime.UTC(),
		ContentHash: hash,
		Title:       title,
		DurationS:   durationS,
		Width:       width,
		Height:      height,
		VideoCodec:  videoCodec,
		AudioCodec:  audioCodec,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// newUUID generates a RFC4122 v4 UUID string using cryptographically secure random bytes.
func newUUID() (string, error) {
	var uuid [16]byte
	if _, err := rand.Read(uuid[:]); err != nil {
		return "", err
	}
	// Set version to 4
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	// Set variant to RFC4122
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16]), nil
}
