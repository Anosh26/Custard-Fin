package memory

import (
	"context"
	"strings"
	"sync"

	"github.com/anosh26/custardfin/apps/api/internal/domain"
	"github.com/anosh26/custardfin/apps/api/internal/ports"
)

type InMemoryCatalogRepository struct {
	mu          sync.RWMutex
	items       map[string]*domain.CatalogItem // key: nodeID + ":" + mediaID
	byMediaID   map[string]*domain.CatalogItem // key: mediaID
	SimulateErr error
}

func NewInMemoryCatalogRepository() *InMemoryCatalogRepository {
	return &InMemoryCatalogRepository{
		items:     make(map[string]*domain.CatalogItem),
		byMediaID: make(map[string]*domain.CatalogItem),
	}
}

var _ ports.CatalogRepository = (*InMemoryCatalogRepository)(nil)

func (r *InMemoryCatalogRepository) UpsertManifest(ctx context.Context, item *domain.CatalogItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.SimulateErr != nil {
		return r.SimulateErr
	}

	key := item.NodeID + ":" + item.MediaID
	r.items[key] = item
	r.byMediaID[item.MediaID] = item
	return nil
}

func (r *InMemoryCatalogRepository) DeleteManifest(ctx context.Context, nodeID, mediaID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.SimulateErr != nil {
		return r.SimulateErr
	}

	key := nodeID + ":" + mediaID
	if _, exists := r.items[key]; !exists {
		return domain.ErrNotFound
	}

	delete(r.items, key)
	delete(r.byMediaID, mediaID)
	return nil
}

func (r *InMemoryCatalogRepository) Search(ctx context.Context, query string, limit, offset int) ([]*domain.CatalogItem, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.SimulateErr != nil {
		return nil, 0, r.SimulateErr
	}

	query = strings.ToLower(strings.TrimSpace(query))
	var matched []*domain.CatalogItem

	for _, item := range r.items {
		if !item.IsPublic {
			continue
		}
		if query == "" || strings.Contains(strings.ToLower(item.Title), query) || strings.Contains(strings.ToLower(item.Description), query) {
			matched = append(matched, item)
		}
	}

	total := len(matched)
	if offset >= total {
		return []*domain.CatalogItem{}, total, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return matched[offset:end], total, nil
}

func (r *InMemoryCatalogRepository) FindByMediaID(ctx context.Context, mediaID string) (*domain.CatalogItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.SimulateErr != nil {
		return nil, r.SimulateErr
	}

	item, exists := r.byMediaID[mediaID]
	if !exists {
		return nil, domain.ErrNotFound
	}
	return item, nil
}
