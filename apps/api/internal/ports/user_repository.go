package ports

import (
	"context"

	"github.com/anosh26/custardfin/apps/api/internal/domain"
)

// UserRepository provides persistence operations for User entities.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
}
