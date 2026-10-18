package sqlite

import (
	"context"
	"testing"
	"time"

	"github.com/anosh26/custardfin/apps/edge/internal/domain"
)

func TestMediaRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	db, err := Open(ctx, ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	defer db.Close()

	repo := NewMediaRepository(db)

	now := time.Now().UTC().Truncate(time.Millisecond)
	m1, err := domain.NewMedia(
		"/media/tears_of_steel.mp4",
		104857600,
		now,
		"hash_tears_of_steel_123",
		734,
		1920,
		1080,
		"h264",
		"aac",
	)
	if err != nil {
		t.Fatalf("failed constructing media entity: %v", err)
	}

	// 1. Save (Insert)
	if err := repo.Save(ctx, m1); err != nil {
		t.Fatalf("failed to insert media: %v", err)
	}

	// 2. FindByPath
	gotByPath, err := repo.FindByPath(ctx, "/media/tears_of_steel.mp4")
	if err != nil {
		t.Fatalf("FindByPath returned error: %v", err)
	}
	if gotByPath == nil {
		t.Fatal("expected media to be found by path, got nil")
	}
	if gotByPath.ID != m1.ID {
		t.Errorf("expected ID %s, got %s", m1.ID, gotByPath.ID)
	}
	if gotByPath.Title != "tears_of_steel" {
		t.Errorf("expected title 'tears_of_steel', got %q", gotByPath.Title)
	}
	if gotByPath.ContentHash != "hash_tears_of_steel_123" {
		t.Errorf("expected content hash 'hash_tears_of_steel_123', got %q", gotByPath.ContentHash)
	}

	// 3. FindByID
	gotByID, err := repo.FindByID(ctx, m1.ID)
	if err != nil || gotByID == nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if gotByID.Path != m1.Path {
		t.Errorf("expected path %s, got %s", m1.Path, gotByID.Path)
	}

	// 4. FindByHash
	gotByHash, err := repo.FindByHash(ctx, "hash_tears_of_steel_123")
	if err != nil || gotByHash == nil {
		t.Fatalf("FindByHash failed: %v", err)
	}
	if gotByHash.ID != m1.ID {
		t.Errorf("expected ID %s, got %s", m1.ID, gotByHash.ID)
	}

	// 5. Update via Save
	gotByID.Title = "Tears of Steel Remastered"
	gotByID.FileSize = 209715200
	if err := repo.Save(ctx, gotByID); err != nil {
		t.Fatalf("failed updating media: %v", err)
	}

	updated, err := repo.FindByID(ctx, m1.ID)
	if err != nil || updated == nil {
		t.Fatalf("FindByID failed after update: %v", err)
	}
	if updated.Title != "Tears of Steel Remastered" || updated.FileSize != 209715200 {
		t.Errorf("expected updated title and size, got %s and %d", updated.Title, updated.FileSize)
	}

	// 6. List
	list, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 item in list, got %d", len(list))
	}

	// 7. Non-existent lookups return nil, nil
	notFoundPath, err := repo.FindByPath(ctx, "/does/not/exist.mp4")
	if err != nil || notFoundPath != nil {
		t.Errorf("expected nil, nil for missing path, got: %v, %v", notFoundPath, err)
	}
	notFoundID, err := repo.FindByID(ctx, "non-existent-uuid")
	if err != nil || notFoundID != nil {
		t.Errorf("expected nil, nil for missing id, got: %v, %v", notFoundID, err)
	}
	notFoundHash, err := repo.FindByHash(ctx, "non-existent-hash")
	if err != nil || notFoundHash != nil {
		t.Errorf("expected nil, nil for missing hash, got: %v, %v", notFoundHash, err)
	}

	// 8. Delete
	if err := repo.Delete(ctx, m1.ID); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	deletedCheck, err := repo.FindByID(ctx, m1.ID)
	if err != nil || deletedCheck != nil {
		t.Errorf("expected nil after deletion, got: %v, %v", deletedCheck, err)
	}
}
