package pubsub

import (
	"context"
	"encoding/json"
	"time"

	"github.com/e-wrobel/cash_ai/internal/event"
	"github.com/segmentio/kafka-go"
)

type KafkaBus struct {
	writer       *kafka.Writer
	readerConfig kafka.ReaderConfig
}

func NewKafkaBus(brokers []string, topic, groupID string) *KafkaBus {
	return &KafkaBus{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireOne,
		},
		readerConfig: kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 1,
			MaxBytes: 10e6,
		},
	}
}

func (k *KafkaBus) Subscribe() <-chan event.Event {
	ch := make(chan event.Event, 256)
	go func() {
		reader := kafka.NewReader(k.readerConfig)
		defer reader.Close()
		for {
			m, err := reader.ReadMessage(context.Background())
			if err != nil {
				time.Sleep(time.Second)
				continue
			}
			var e event.Event
			if err := json.Unmarshal(m.Value, &e); err == nil {
				ch <- e
			}
		}
	}()
	return ch
}

func (k *KafkaBus) Publish(e event.Event) {
	data, err := json.Marshal(e)
	if err != nil {
		return
	}
	_ = k.writer.WriteMessages(context.Background(), kafka.Message{
		Value: data,
	})
}
