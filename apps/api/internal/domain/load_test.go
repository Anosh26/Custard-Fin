package domain_test

import (
	"testing"

	"github.com/anosh26/custardfin/apps/api/internal/domain"
)

func TestLoadLevel_AllowedResolution(t *testing.T) {
	tests := []struct {
		name       string
		level      domain.LoadLevel
		resolution string
		want       bool
	}{
		{"Normal allows 1080p", domain.LoadNormal, "1080p", true},
		{"Normal allows 720p", domain.LoadNormal, "720p", true},
		{"Normal allows 480p", domain.LoadNormal, "480p", true},
		{"Moderate blocks 1080p", domain.LoadModerate, "1080p", false},
		{"Moderate blocks 1920x1080", domain.LoadModerate, "1920x1080", false},
		{"Moderate allows 720p", domain.LoadModerate, "720p", true},
		{"Moderate allows 480p", domain.LoadModerate, "480p", true},
		{"Saturated blocks 1080p", domain.LoadSaturated, "1080p", false},
		{"Saturated blocks 720p", domain.LoadSaturated, "720p", false},
		{"Saturated allows 480p", domain.LoadSaturated, "480p", true},
		{"Saturated allows 854x480", domain.LoadSaturated, "854x480", true},
		{"Exhausted blocks 1080p", domain.LoadExhausted, "1080p", false},
		{"Exhausted blocks 720p", domain.LoadExhausted, "720p", false},
		{"Exhausted blocks 480p", domain.LoadExhausted, "480p", false},
		{"Unknown defaults to true", domain.LoadLevel("UNKNOWN"), "1080p", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.level.AllowedResolution(tt.resolution)
			if got != tt.want {
				t.Errorf("AllowedResolution(%q) = %v, want %v for level %s", tt.resolution, got, tt.want, tt.level)
			}
		})
	}
}
