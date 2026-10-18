package sqlite

import (
	"context"
	"testing"
	"time"

	"github.com/anosh26/custardfin/apps/edge/internal/domain"
)

func TestWatchProgressRepository_CRUD_AndCascade(t *testing.T) {
	ctx := context.Background()
	db, err := Open(ctx, ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	defer db.Close()

	mediaRepo := NewMediaRepository(db)
	watchRepo := NewWatchProgressRepository(db)

	// Create media record first (for foreign key constraint)
	m, err := domain.NewMedia(
		"/media/sintel.mp4",
		52428800,
		time.Now().UTC(),
		"hash_sintel_456",
		888,
		1920,
		1080,
		"h264",
		"aac",
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := mediaRepo.Save(ctx, m); err != nil {
		t.Fatal(err)
	}

	// 1. Save (Insert) watch progress
	wp, err := domain.NewWatchProgress(m.ID, 120, false)
	if err != nil {
		t.Fatal(err)
	}
	if err := watchRepo.SaveWatchProgress(ctx, wp); err != nil {
		t.Fatalf("failed saving watch progress: %v", err)
	}

	// 2. GetWatchProgress
	gotWP, err := watchRepo.GetWatchProgress(ctx, m.ID)
	if err != nil {
		t.Fatalf("GetWatchProgress failed: %v", err)
	}
	if gotWP == nil {
		t.Fatal("expected watch progress, got nil")
	}
	if gotWP.MediaID != m.ID || gotWP.PositionS != 120 || gotWP.Completed {
		t.Errorf("unexpected watch progress fields: %+v", gotWP)
	}

	// 3. Update watch progress (advance position and complete)
	gotWP.PositionS = 888
	gotWP.Completed = true
	gotWP.LastWatchedAt = time.Now().UTC()
	if err := watchRepo.SaveWatchProgress(ctx, gotWP); err != nil {
		t.Fatalf("failed updating watch progress: %v", err)
	}

	updatedWP, err := watchRepo.GetWatchProgress(ctx, m.ID)
	if err != nil || updatedWP == nil {
		t.Fatalf("GetWatchProgress failed after update: %v", err)
	}
	if updatedWP.PositionS != 888 || !updatedWP.Completed {
		t.Errorf("expected position 888 and completed true, got %+v", updatedWP)
	}

	// 4. Non-existent watch progress returns nil, nil
	missingWP, err := watchRepo.GetWatchProgress(ctx, "non-existent-media-id")
	if err != nil || missingWP != nil {
		t.Errorf("expected nil, nil for non-existent progress, got: %v, %v", missingWP, err)
	}

	// 5. DeleteWatchProgress directly
	if err := watchRepo.DeleteWatchProgress(ctx, m.ID); err != nil {
		t.Fatalf("failed deleting watch progress: %v", err)
	}
	deletedWP, err := watchRepo.GetWatchProgress(ctx, m.ID)
	if err != nil || deletedWP != nil {
		t.Errorf("expected nil after delete, got: %v, %v", deletedWP, err)
	}

	// 6. Test Foreign Key Cascade
	// Re-add watch progress
	if err := watchRepo.SaveWatchProgress(ctx, wp); err != nil {
		t.Fatal(err)
	}
	// Delete parent media
	if err := mediaRepo.Delete(ctx, m.ID); err != nil {
		t.Fatal(err)
	}
	// Verify watch_progress was automatically cascaded
	cascadedWP, err := watchRepo.GetWatchProgress(ctx, m.ID)
	if err != nil || cascadedWP != nil {
		t.Errorf("expected watch progress to cascade-delete when media was deleted, got: %v, %v", cascadedWP, err)
	}
}
