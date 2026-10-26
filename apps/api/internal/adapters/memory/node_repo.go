package memory

import (
	"context"
	"strings"
	"sync"

	"github.com/anosh26/custardfin/apps/api/internal/domain"
	"github.com/anosh26/custardfin/apps/api/internal/ports"
)

type InMemoryNodeRepository struct {
	mu          sync.RWMutex
	byID        map[string]*domain.Node
	byTokenHash map[string]*domain.Node
	SimulateErr error
}

func NewInMemoryNodeRepository() *InMemoryNodeRepository {
	return &InMemoryNodeRepository{
		byID:        make(map[string]*domain.Node),
		byTokenHash: make(map[string]*domain.Node),
	}
}

var _ ports.NodeRepository = (*InMemoryNodeRepository)(nil)

func (r *InMemoryNodeRepository) Create(ctx context.Context, node *domain.Node) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.SimulateErr != nil {
		return r.SimulateErr
	}

	r.byID[node.ID] = node
	r.byTokenHash[node.NodeTokenHash] = node
	return nil
}

func (r *InMemoryNodeRepository) FindByID(ctx context.Context, id string) (*domain.Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.SimulateErr != nil {
		return nil, r.SimulateErr
	}

	node, exists := r.byID[id]
	if !exists {
		return nil, domain.ErrNotFound
	}
	return node, nil
}

func (r *InMemoryNodeRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*domain.Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.SimulateErr != nil {
		return nil, r.SimulateErr
	}

	node, exists := r.byTokenHash[tokenHash]
	if !exists {
		return nil, domain.ErrNotFound
	}
	return node, nil
}

func (r *InMemoryNodeRepository) ListByUserID(ctx context.Context, userID string) ([]*domain.Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.SimulateErr != nil {
		return nil, r.SimulateErr
	}

	var nodes []*domain.Node
	for _, n := range r.byID {
		if n.UserID == userID {
			nodes = append(nodes, n)
		}
	}
	return nodes, nil
}

func (r *InMemoryNodeRepository) UpdateOwner(ctx context.Context, nodeID, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.SimulateErr != nil {
		return r.SimulateErr
	}

	node, exists := r.byID[nodeID]
	if !exists {
		return domain.ErrNotFound
	}
	node.UserID = userID
	return nil
}

func (r *InMemoryNodeRepository) UpdateName(ctx context.Context, nodeID, name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.SimulateErr != nil {
		return r.SimulateErr
	}

	node, exists := r.byID[nodeID]
	if !exists {
		return domain.ErrNotFound
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return domain.ErrValidation
	}
	node.Name = name
	return nil
}

