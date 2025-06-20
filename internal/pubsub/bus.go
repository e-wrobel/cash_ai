package pubsub

import (
	"sync"

	"github.com/e-wrobel/cash_ai/internal/event"
)

// Bus is a minimal in‑memory fan‑out message bus (replace with Kafka, NATS …).
type Bus struct {
	mu          sync.Mutex
	subscribers []chan event.Event
}

func New() *Bus { return &Bus{} }

// Subscribe returns a read‑only channel receiving every event.
func (b *Bus) Subscribe() <-chan event.Event {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan event.Event, 256)
	b.subscribers = append(b.subscribers, ch)
	return ch
}

// Publish fan‑outs to all subscribers (non‑blocking best‑effort).
func (b *Bus) Publish(e event.Event) {
	b.mu.Lock()
	subs := append([]chan event.Event(nil), b.subscribers...)
	b.mu.Unlock()
	for _, ch := range subs {
		select {
		case ch <- e:
		default:
		}
	}
}
