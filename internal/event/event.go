package event

import (
	"time"

	"github.com/e-wrobel/cash_ai/internal/domain"
	"github.com/google/uuid"
)

// Type enumerates domain event kinds.
type Type string

const (
	Created Type = "created"
	Updated Type = "updated"
	Deleted Type = "deleted"
)

// Event is an immutable ledger entry.
type Event struct {
	ID          uuid.UUID          `json:"id"`
	Kind        Type               `json:"kind"`
	Transaction domain.Transaction `json:"transaction"`
	At          time.Time          `json:"at"`
}

func New(kind Type, tx domain.Transaction) Event {
	return Event{
		ID:          uuid.New(),
		Kind:        kind,
		Transaction: tx,
		At:          time.Now().UTC(),
	}
}
