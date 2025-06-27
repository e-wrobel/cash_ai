package pubsub

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/e-wrobel/cash_ai/internal/event"
)

func TestSubscribe(t *testing.T) {
	bus := New()
	ch := bus.Subscribe()
	assert.NotNil(t, ch)
}

func TestPublishToSingleSubscriber(t *testing.T) {
	bus := New()
	ch := bus.Subscribe()
	ev := event.Event{Kind: "test"}
	bus.Publish(ev)

	select {
	case received := <-ch:
		assert.Equal(t, ev, received)
	case <-time.After(time.Second):
		t.Error("Didn't receive event on channel")
	}
}

func TestPublishToMultipleSubscribers(t *testing.T) {
	bus := New()
	ch1 := bus.Subscribe()
	ch2 := bus.Subscribe()
	ev := event.Event{Kind: "multi"}

	bus.Publish(ev)

	for _, ch := range []<-chan event.Event{ch1, ch2} {
		select {
		case received := <-ch:
			if received != ev {
				t.Errorf("Expected %v, received %v", ev, received)
			}
		case <-time.After(time.Second):
			t.Error("Timeout while waiting for event on channel")
		}
	}
}
