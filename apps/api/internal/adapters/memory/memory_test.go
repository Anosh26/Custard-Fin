package memory_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/anosh26/custardfin/apps/api/internal/adapters/memory"
	"github.com/anosh26/custardfin/apps/api/internal/domain"
)

func TestInMemoryUserRepository(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewInMemoryUserRepository()

	user, err := domain.NewUser("alice", "hashed_pwd_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 1. Create user
	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// 2. Duplicate username error
	dupUser, _ := domain.NewUser("alice", "another_hash")
	err = repo.Create(ctx, dupUser)
	if err == nil || !errors.Is(err, domain.ErrConflict) {
		t.Errorf("expected ErrConflict for duplicate user, got %v", err)
	}

	// 3. FindByUsername
	found, err := repo.FindByUsername(ctx, "alice")
	if err != nil {
		t.Fatalf("failed to find by username: %v", err)
	}
	if found.ID != user.ID {
		t.Errorf("expected ID %s, got %s", user.ID, found.ID)
	}

	// 4. FindByID
	foundByID, err := repo.FindByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to find by ID: %v", err)
	}
	if foundByID.Username != "alice" {
		t.Errorf("expected username alice, got %s", foundByID.Username)
	}

	// 5. Not found
	_, err = repo.FindByUsername(ctx, "nonexistent")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestInMemoryNodeRepository(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewInMemoryNodeRepository()

	node, err := domain.NewNode("user-1", "Living Room Node", "hash_abc", "frp_secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Create
	if err := repo.Create(ctx, node); err != nil {
		t.Fatalf("failed to create node: %v", err)
	}

	// FindByID
	n, err := repo.FindByID(ctx, node.ID)
	if err != nil {
		t.Fatalf("failed to find node: %v", err)
	}
	if n.Name != "Living Room Node" {
		t.Errorf("expected node name 'Living Room Node', got %s", n.Name)
	}

	// FindByTokenHash
	nToken, err := repo.FindByTokenHash(ctx, "hash_abc")
	if err != nil {
		t.Fatalf("failed to find node by token hash: %v", err)
	}
	if nToken.ID != node.ID {
		t.Errorf("expected node ID %s, got %s", node.ID, nToken.ID)
	}

	// ListByUserID
	nodes, err := repo.ListByUserID(ctx, "user-1")
	if err != nil {
		t.Fatalf("failed to list nodes: %v", err)
	}
	if len(nodes) != 1 {
		t.Errorf("expected 1 node, got %d", len(nodes))
	}
}

func TestInMemoryCatalogRepository(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewInMemoryCatalogRepository()

	item, err := domain.NewCatalogItem(
		"media-1",
		"node-1",
		"Open Movie Sintel",
		"Fantasy animated short film",
		900,
		"hash_123",
		true,
		[]byte(`{"schema_version":"1.0"}`),
		[]domain.Rendition{
			{Resolution: "1080p", BitrateBps: 5000000, PlaylistPath: "1080p.m3u8", ChunkDurationS: 4},
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Upsert
	if err := repo.UpsertManifest(ctx, item); err != nil {
		t.Fatalf("failed to upsert manifest: %v", err)
	}

	// FindByMediaID
	found, err := repo.FindByMediaID(ctx, "media-1")
	if err != nil {
		t.Fatalf("failed to find by media ID: %v", err)
	}
	if found.Title != "Open Movie Sintel" {
		t.Errorf("expected title 'Open Movie Sintel', got %s", found.Title)
	}

	// Search matching query
	results, total, err := repo.Search(ctx, "sintel", 10, 0, nil)
	if err != nil {
		t.Fatalf("search error: %v", err)
	}
	if total != 1 || len(results) != 1 {
		t.Errorf("expected 1 result, got total=%d len=%d", total, len(results))
	}

	// Search non-matching query
	results, total, err = repo.Search(ctx, "nonexistent", 10, 0, nil)
	if err != nil {
		t.Fatalf("search error: %v", err)
	}
	if total != 0 || len(results) != 0 {
		t.Errorf("expected 0 results, got total=%d len=%d", total, len(results))
	}

	// Test private item filtering
	privItem, err := domain.NewCatalogItem(
		"media-priv",
		"node-priv",
		"Private Video",
		"Private video description",
		30,
		"hash-priv",
		false,
		[]byte(`{"schema_version":"1.0"}`),
		[]domain.Rendition{
			{Resolution: "1080p", BitrateBps: 5000000, PlaylistPath: "1080p.m3u8", ChunkDurationS: 4},
		},
	)
	if err != nil {
		t.Fatalf("unexpected error creating private catalog item: %v", err)
	}
	if err := repo.UpsertManifest(ctx, privItem); err != nil {
		t.Fatalf("failed to upsert private manifest: %v", err)
	}

	// Unauthenticated / public search should not find private video
	results, total, err = repo.Search(ctx, "Private", 10, 0, nil)
	if err != nil {
		t.Fatalf("search error: %v", err)
	}
	if total != 0 || len(results) != 0 {
		t.Errorf("expected 0 results for public search of private video, got %d", total)
	}

	// Search with node-priv in allowedNodeIDs should find private video
	results, total, err = repo.Search(ctx, "Private", 10, 0, []string{"node-priv"})
	if err != nil {
		t.Fatalf("search error: %v", err)
	}
	if total != 1 || len(results) != 1 {
		t.Errorf("expected 1 result when allowedNodeIDs contains node-priv, got total=%d", total)
	}

	// Delete
	if err := repo.DeleteManifest(ctx, "node-1", "media-1"); err != nil {
		t.Fatalf("failed to delete manifest: %v", err)
	}

	// Confirm deleted
	_, err = repo.FindByMediaID(ctx, "media-1")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("expected ErrNotFound after deletion, got %v", err)
	}
}

func TestInMemoryRepositories_SimulatedErrors(t *testing.T) {
	ctx := context.Background()
	simErr := errors.New("simulated storage failure")

	userRepo := memory.NewInMemoryUserRepository()
	userRepo.SimulateErr = simErr

	u, _ := domain.NewUser("testuser", "hash")
	if err := userRepo.Create(ctx, u); !errors.Is(err, simErr) {
		t.Errorf("expected simulated error, got %v", err)
	}
	if _, err := userRepo.FindByUsername(ctx, "testuser"); !errors.Is(err, simErr) {
		t.Errorf("expected simulated error, got %v", err)
	}

	nodeRepo := memory.NewInMemoryNodeRepository()
	nodeRepo.SimulateErr = simErr
	n, _ := domain.NewNode("user-1", "node", "hash", "frp")
	if err := nodeRepo.Create(ctx, n); !errors.Is(err, simErr) {
		t.Errorf("expected simulated error, got %v", err)
	}

	catRepo := memory.NewInMemoryCatalogRepository()
	catRepo.SimulateErr = simErr
	if _, _, err := catRepo.Search(ctx, "", 10, 0, nil); !errors.Is(err, simErr) {
		t.Errorf("expected simulated error, got %v", err)
	}
	_ = time.Now()
}
