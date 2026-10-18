package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/anosh26/custardfin/apps/api/internal/domain"
	"github.com/anosh26/custardfin/apps/api/internal/ports"
)

type InMemoryUserRepository struct {
	mu       sync.RWMutex
	byID     map[string]*domain.User
	byName   map[string]*domain.User
	SimulateErr error
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		byID:   make(map[string]*domain.User),
		byName: make(map[string]*domain.User),
	}
}

var _ ports.UserRepository = (*InMemoryUserRepository)(nil)

func (r *InMemoryUserRepository) Create(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.SimulateErr != nil {
		return r.SimulateErr
	}

	if _, exists := r.byName[user.Username]; exists {
		return fmt.Errorf("%w: user with username %s already exists", domain.ErrConflict, user.Username)
	}

	r.byID[user.ID] = user
	r.byName[user.Username] = user
	return nil
}

func (r *InMemoryUserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.SimulateErr != nil {
		return nil, r.SimulateErr
	}

	u, exists := r.byName[username]
	if !exists {
		return nil, domain.ErrNotFound
	}
	return u, nil
}

func (r *InMemoryUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.SimulateErr != nil {
		return nil, r.SimulateErr
	}

	u, exists := r.byID[id]
	if !exists {
		return nil, domain.ErrNotFound
	}
	return u, nil
}
