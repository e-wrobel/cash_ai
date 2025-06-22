package repository

import (
	"sync"
	"testing"

	"github.com/e-wrobel/cash_ai/internal/domain"
	"github.com/e-wrobel/cash_ai/internal/event"
	"github.com/stretchr/testify/assert"
)

func TestEventStore_AppendAndAll(t *testing.T) {
	evt1 := event.New(event.Created, domain.Transaction{ID: "1"})
	evt2 := event.New(event.Updated, domain.Transaction{ID: "2"})

	tests := []struct {
		name   string
		events []event.Event
		want   []event.Event
	}{
		{
			name:   "no events",
			events: nil,
			want:   []event.Event{},
		},
		{
			name:   "one event",
			events: []event.Event{evt1},
			want:   []event.Event{evt1},
		},
		{
			name:   "multiple events",
			events: []event.Event{evt1, evt2},
			want:   []event.Event{evt1, evt2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := &EventStore{mu: sync.Mutex{}}
			for _, e := range tt.events {
				es.Append(e)
			}
			got := es.All()
			assert.Equal(t, len(tt.want), len(got), "All() got %d events, want %d", len(got), len(tt.want))
			for i := range got {
				assert.Equal(t, tt.want[i], got[i], "event at %d = %+v, want %+v", i, got[i], tt.want[i])
			}
			if len(got) > 0 {
				got[0] = event.New(event.Deleted, domain.Transaction{ID: "x"})
				assert.NotEqual(t, es.events[0], got[0], "All() should return a copy, not a reference")
			}
		})
	}
}

func TestNewEventStore(t *testing.T) {
	es := NewEventStore()
	if es == nil {
		t.Error("NewEventStore() returned nil")
	}
	assert.Equalf(t, len(es.All()), 0, "NewEventStore() should start with empty events slice")
	assert.NotEqualf(t, es, nil, "NewEventStore() should return a non-nil EventStore instance")
}
