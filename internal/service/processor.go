package service

import (
	"github.com/e-wrobel/cash_ai/internal/domain"
	"github.com/e-wrobel/cash_ai/internal/event"
	"github.com/e-wrobel/cash_ai/internal/pubsub"
	"github.com/e-wrobel/cash_ai/internal/repository"
)

type Processor struct {
	store *repository.EventStore
	bus   *pubsub.Bus
	// Fast projection of latest state.
	state map[string]domain.Transaction
}

func New(store *repository.EventStore, bus *pubsub.Bus) *Processor {
	return &Processor{store: store, bus: bus, state: map[string]domain.Transaction{}}
}

// Apply processes one transaction against current state and emits events.
func (p *Processor) Apply(tx domain.Transaction) {
	cur, ok := p.state[tx.ID]
	switch {
	case !ok:
		p.emit(event.Created, tx)
	case !cur.Equal(tx):
		p.emit(event.Updated, tx)
	default:
		// unchanged – nothing
	}
}

// Delete marks a missing transaction as deleted.
func (p *Processor) Delete(id string) {
	cur, ok := p.state[id]
	if !ok {
		return
	}
	p.emit(event.Deleted, cur)
}

func (p *Processor) emit(kind event.Type, tx domain.Transaction) {
	e := event.New(kind, tx)
	p.store.Append(e)
	p.bus.Publish(e)
	switch kind {
	case event.Deleted:
		delete(p.state, tx.ID)
	default:
		p.state[tx.ID] = tx
	}
}

// Replay rebuilds state from scratch.
func (p *Processor) Replay() {
	p.state = map[string]domain.Transaction{}
	for _, e := range p.store.All() {
		switch e.Kind {
		case event.Created, event.Updated:
			p.state[e.Transaction.ID] = e.Transaction
		case event.Deleted:
			delete(p.state, e.Transaction.ID)
		}
	}
}

func (p *Processor) StateIDs() map[string]struct{} {
	ids := make(map[string]struct{}, len(p.state))
	for id := range p.state {
		ids[id] = struct{}{}
	}
	return ids
}
