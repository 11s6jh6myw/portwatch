package metrics

import (
	"encoding/json"
	"net/http"
)

// Server exposes metrics over HTTP.
type Server struct {
	collector *Collector
	addr      string
}

// NewServer returns a metrics HTTP server bound to addr.
func NewServer(addr string, c *Collector) *Server {
	return &Server{addr: addr, collector: c}
}

// ListenAndServe starts the HTTP server. It blocks until the server stops.
func (s *Server) ListenAndServe() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", s.handleMetrics)
	mux.HandleFunc("/healthz", s.handleHealth)
	return http.ListenAndServe(s.addr, mux)
}

func (s *Server) handleMetrics(w http.ResponseWriter, _ *http.Request) {
	snap := s.collector.Snapshot()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(snap)
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
