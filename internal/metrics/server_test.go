package metrics_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

func TestHandleMetrics_ReturnsJSON(t *testing.T) {
	c := metrics.New()
	c.RecordScan(15*time.Millisecond, 4)
	c.RecordAlert()

	srv := metrics.NewServer(":0", c)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)

	// Exercise handler directly via exported helper.
	srv.ServeHTrics(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var snap metrics.Snapshot
	if err := json.NewDecoder(rr.Body).Decode(&snap); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if snap.ScansTotal != 1 {
		t.Errorf("expected 1 scan, got %d", snap.ScansTotal)
	}
	if snap.AlertsTotal != 1 {
		t.Errorf("expected 1 alert, got %d", snap.AlertsTotal)
	}
	if snap.OpenPorts != 4 {
		t.Errorf("expected 4 open ports, got %d", snap.OpenPorts)
	}
}

func TestHandleHealth_Returns200(t *testing.T) {
	c := metrics.New()
	srv := metrics.NewServer(":0", c)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	srv.ServeHealth(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}
