package domain

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidMediaID  = errors.New("media_id cannot be empty")
	ErrNegativePosition = errors.New("watch position cannot be negative")
)

// WatchProgress records the local user's playback position for a media item.
// Invariant 8: This information lives purely in edge SQLite and is never sent to the cloud.
type WatchProgress struct {
	MediaID       string    `json:"media_id"`
	PositionS     int       `json:"position_s"`
	Completed     bool      `json:"completed"`
	LastWatchedAt time.Time `json:"last_watched_at"`
}

// NewWatchProgress constructs a validated WatchProgress entity.
func NewWatchProgress(mediaID string, positionS int, completed bool) (*WatchProgress, error) {
	if strings.TrimSpace(mediaID) == "" {
		return nil, ErrInvalidMediaID
	}
	if positionS < 0 {
		return nil, ErrNegativePosition
	}
	return &WatchProgress{
		MediaID:       mediaID,
		PositionS:     positionS,
		Completed:     completed,
		LastWatchedAt: time.Now().UTC(),
	}, nil
}
