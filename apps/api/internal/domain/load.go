package domain

// LoadLevel represents the saturation state of a home edge node.
type LoadLevel string

const (
	// LoadNormal represents 0-3 active streams: full ladder (1080p, 720p, 480p) enabled.
	LoadNormal LoadLevel = "NORMAL"
	// LoadModerate represents 4-6 active streams: 1080p pruned, serving 720p and 480p.
	LoadModerate LoadLevel = "MODERATE"
	// LoadSaturated represents 7-9 active streams: 1080p and 720p pruned, serving 480p only.
	LoadSaturated LoadLevel = "SATURATED"
	// LoadExhausted represents 10+ active streams: capacity exhausted (Invariant 10).
	LoadExhausted LoadLevel = "EXHAUSTED"
)

// AllowedResolution reports whether a specific resolution variant is permitted under the current load level.
func (l LoadLevel) AllowedResolution(resolution string) bool {
	switch l {
	case LoadNormal:
		return true
	case LoadModerate:
		return resolution != "1080p" && resolution != "1920x1080"
	case LoadSaturated:
		return resolution == "480p" || resolution == "854x480" || resolution == "640x480"
	case LoadExhausted:
		return false
	default:
		return true
	}
}
