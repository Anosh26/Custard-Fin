package domain

import "errors"

var (
	// ErrNotFound is returned when a requested entity does not exist.
	ErrNotFound = errors.New("resource not found")

	// ErrConflict is returned when an entity with unique fields already exists.
	ErrConflict = errors.New("resource already exists")

	// ErrUnauthorized is returned when authentication fails or is missing.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden is returned when the authenticated caller lacks permission.
	ErrForbidden = errors.New("forbidden")

	// ErrValidation is returned when input data fails domain validation rules.
	ErrValidation = errors.New("validation failed")

	// ErrManifestTooLarge is returned when manifest exceeds 3 KB cap.
	ErrManifestTooLarge = errors.New("manifest exceeds 3 KB limit")
)
