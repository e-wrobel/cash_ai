package api

import (
	"log"
	"net/http"
	"time"

	"github.com/e-wrobel/cash_ai/internal/adapter/mockprovider"
	"github.com/e-wrobel/cash_ai/internal/service"
)

type Server struct {
	proc   *service.Processor
	client *mockprovider.Client
}

func New(proc *service.Processor, client *mockprovider.Client) *Server {
	return &Server{proc: proc, client: client}
}

func (s *Server) routes() {
	http.HandleFunc("/trigger/fetch", s.handleFetch)
	http.HandleFunc("/replay", s.handleReplay)
}

func (s *Server) handleFetch(w http.ResponseWriter, _ *http.Request) {
	if err := s.fetchOnce(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) handleReplay(w http.ResponseWriter, _ *http.Request) {
	s.proc.Replay()
	w.WriteHeader(http.StatusOK)
}

func (s *Server) fetchOnce() error {
	txs, err := s.client.Fetch()
	if err != nil {
		return err
	}
	seen := map[string]struct{}{}
	for _, tx := range txs {
		seen[tx.ID] = struct{}{}
		s.proc.Apply(tx)
	}
	// delete missing
	for id := range s.proc.StateIDs() {
		if _, ok := seen[id]; !ok {
			s.proc.Delete(id)
		}
	}
	return nil
}

func (s *Server) Serve(addr string, interval time.Duration) error {
	s.routes()
	// background polling
	go func() {
		t := time.NewTicker(interval)
		for range t.C {
			if err := s.fetchOnce(); err != nil {
				log.Printf("poll error: %v", err)
			}
		}
	}()
	log.Printf("HTTP server on %s", addr)
	return http.ListenAndServe(addr, nil)
}
