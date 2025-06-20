package repository

import (
	"sync"

	"github.com/e-wrobel/cash_ai/internal/event"
)

// EventStore persists every change forever and can replay history.
type EventStore struct {
	mu     sync.Mutex
	events []event.Event
}

func NewEventStore() *EventStore { return &EventStore{} }

func (s *EventStore) Append(e event.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, e)
}

// All returns a *copy* of the ledger.
func (s *EventStore) All() []event.Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	dup := make([]event.Event, len(s.events))
	copy(dup, s.events)
	return dup
}
