package contracts

import (
	"errors"
	"testing"
	"time"
)

func TestHeartbeatValidation_HappyPath(t *testing.T) {
	hb := &Heartbeat{
		NodeID:  "550e8400-e29b-41d4-a716-446655440001",
		SentAt:  time.Now().UTC(),
		Version: "0.1.0",
		GPUCaps: &GPUCapabilities{
			NVENC: true,
			QSV:   false,
			VAAPI: false,
		},
		ActiveStreams: 5,
	}

	if err := hb.Validate(); err != nil {
		t.Fatalf("expected valid heartbeat, got: %v", err)
	}

	// Boundary cap check: exactly 10 streams is allowed
	hb.ActiveStreams = MaxConcurrentStreamsPerNode
	if err := hb.Validate(); err != nil {
		t.Errorf("expected 10 active streams to be valid, got: %v", err)
	}
}

func TestHeartbeatValidation_InvalidFields(t *testing.T) {
	tests := []struct {
		name        string
		hb          Heartbeat
		expectedErr error
	}{
		{
			name: "empty node_id",
			hb: Heartbeat{
				NodeID:        "",
				ActiveStreams: 0,
			},
			expectedErr: ErrInvalidNodeID,
		},
		{
			name: "negative active streams",
			hb: Heartbeat{
				NodeID:        "node-1",
				ActiveStreams: -1,
			},
			expectedErr: ErrNegativeStreamCount,
		},
		{
			name: "exceeded stream cap (11)",
			hb: Heartbeat{
				NodeID:        "node-1",
				ActiveStreams: 11,
			},
			expectedErr: ErrExceededStreamCap,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.hb.Validate()
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}
