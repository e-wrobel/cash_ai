package main

import (
	"log"
	"os"
	"time"

	"github.com/e-wrobel/cash_ai/internal/adapter/http/api"
	"github.com/e-wrobel/cash_ai/internal/adapter/mockprovider"
	"github.com/e-wrobel/cash_ai/internal/pubsub"
	"github.com/e-wrobel/cash_ai/internal/repository"
	"github.com/e-wrobel/cash_ai/internal/service"
)

func main() {
	// Mock Provider URL and pull interval.
	providerURL := getenv("PROVIDER_URL", "http://localhost:9000")
	interval, err := time.ParseDuration(getenv("PULL_INTERVAL", "30s"))
	if err != nil {
		log.Fatalf("bad PULL_INTERVAL: %v", err)
	}

	// Initialize components.
	// Event store, service and mock provider.
	store := repository.NewEventStore()
	bus := pubsub.New()
	proc := service.New(store, bus)
	prov := mockprovider.New(providerURL)

	// Subscriber.
	go func() {
		for e := range bus.Subscribe() {
			log.Printf("event %s – %s", e.Kind, e.Transaction.ID)
		}
	}()

	srv := api.New(proc, prov)
	if err := srv.Serve(":8080", interval); err != nil {
		log.Fatal(err)
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
