package ports

import (
	"context"

	"github.com/anosh26/custardfin/apps/api/internal/domain"
)

// NodeRepository provides persistence operations for Edge Node entities.
type NodeRepository interface {
	Create(ctx context.Context, node *domain.Node) error
	FindByID(ctx context.Context, id string) (*domain.Node, error)
	FindByTokenHash(ctx context.Context, tokenHash string) (*domain.Node, error)
	ListByUserID(ctx context.Context, userID string) ([]*domain.Node, error)
	UpdateOwner(ctx context.Context, nodeID, userID string) error
}
