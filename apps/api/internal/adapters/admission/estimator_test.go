package admission_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/anosh26/custardfin/apps/api/internal/adapters/admission"
	"github.com/anosh26/custardfin/apps/api/internal/adapters/memory"
	"github.com/anosh26/custardfin/apps/api/internal/domain"
	"github.com/anosh26/custardfin/libs/contracts"
)

func TestDefaultEstimator_EstimateLoad(t *testing.T) {
	ctx := context.Background()

	t.Run("nil heartbeat store returns LoadNormal", func(t *testing.T) {
		est := admission.NewDefaultEstimator(admission.Config{}, nil, nil)
		level, active, err := est.EstimateLoad(ctx, "node-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if level != domain.LoadNormal || active != 0 {
			t.Fatalf("expected LoadNormal with 0 active, got %s with %d", level, active)
		}
	})

	t.Run("offline node returns ErrNodeOffline", func(t *testing.T) {
		store := memory.NewInMemoryHeartbeatStore()
		est := admission.NewDefaultEstimator(admission.Config{}, store, nil)
		level, active, err := est.EstimateLoad(ctx, "offline-node")
		if !errors.Is(err, admission.ErrNodeOffline) {
			t.Fatalf("expected ErrNodeOffline, got %v", err)
		}
		if level != domain.LoadExhausted || active != 0 {
			t.Fatalf("expected LoadExhausted with 0, got %s with %d", level, active)
		}
	})

	t.Run("heartbeat store error propagates", func(t *testing.T) {
		store := memory.NewInMemoryHeartbeatStore()
		store.SimulateErr = errors.New("redis timeout")
		est := admission.NewDefaultEstimator(admission.Config{}, store, nil)
		_, _, err := est.EstimateLoad(ctx, "node-1")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("evaluates load tiers correctly", func(t *testing.T) {
		store := memory.NewInMemoryHeartbeatStore()
		est := admission.NewDefaultEstimator(admission.Config{
			ModerateThreshold:  4,
			SaturatedThreshold: 7,
			ExhaustedThreshold: 10,
		}, store, nil)

		tests := []struct {
			name          string
			activeStreams int
			expectedLevel domain.LoadLevel
		}{
			{"0 streams is Normal", 0, domain.LoadNormal},
			{"2 streams is Normal", 2, domain.LoadNormal},
			{"3 streams is Normal", 3, domain.LoadNormal},
			{"4 streams is Moderate", 4, domain.LoadModerate},
			{"6 streams is Moderate", 6, domain.LoadModerate},
			{"7 streams is Saturated", 7, domain.LoadSaturated},
			{"9 streams is Saturated", 9, domain.LoadSaturated},
			{"10 streams is Exhausted", 10, domain.LoadExhausted},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				hb := &contracts.Heartbeat{
					NodeID:        "node-active",
					SentAt:        time.Now(),
					Version:       "1.0.0",
					ActiveStreams: tt.activeStreams,
				}
				if err := store.RecordHeartbeat(ctx, hb, time.Minute); err != nil {
					t.Fatalf("failed recording heartbeat: %v", err)
				}

				level, active, err := est.EstimateLoad(ctx, "node-active")
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if level != tt.expectedLevel {
					t.Errorf("level = %s, want %s for active streams %d", level, tt.expectedLevel, tt.activeStreams)
				}
				if active != tt.activeStreams {
					t.Errorf("active = %d, want %d", active, tt.activeStreams)
				}
			})
		}
	})
}
