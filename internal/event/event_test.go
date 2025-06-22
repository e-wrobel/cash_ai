package event

import (
	"testing"
	"time"

	"github.com/e-wrobel/cash_ai/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tx := domain.Transaction{
		ID:        "tx1",
		UserID:    "user1",
		Amount:    100.0,
		Type:      "credit",
		Timestamp: 1234567890,
	}

	tests := []struct {
		name string
		kind Type
		tx   domain.Transaction
	}{
		{
			name: "basic created event",
			kind: Created,
			tx:   tx,
		},
		{
			name: "basic updated event",
			kind: Updated,
			tx:   tx,
		},
		{
			name: "basic deleted event",
			kind: Deleted,
			tx:   tx,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evt := New(tt.kind, tt.tx)
			assert.Equal(t, tt.kind, evt.Kind)
			assert.Equal(t, tt.tx, evt.Transaction)
			assert.NotEqual(t, uuid.Nil, evt.ID)
			assert.WithinDuration(t, time.Now().UTC(), evt.At, time.Second)
		})
	}
}
