package domain

import (
	"crypto/rand"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var validUsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,64}$`)

// User represents a registered CustardFin user.
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewUser constructs and validates a new User entity.
func NewUser(username, passwordHash string) (*User, error) {
	username = strings.TrimSpace(username)
	if !validUsernameRegex.MatchString(username) {
		return nil, fmt.Errorf("%w: username must be 3-64 characters and contain only alphanumeric, underscore, or hyphen", ErrValidation)
	}

	if strings.TrimSpace(passwordHash) == "" {
		return nil, fmt.Errorf("%w: password hash cannot be empty", ErrValidation)
	}

	now := time.Now().UTC()
	return &User{
		ID:           newUUIDv4(),
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// newUUIDv4 generates an RFC 4122 compliant v4 UUID.
func newUUIDv4() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // RFC 4122 variant
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
