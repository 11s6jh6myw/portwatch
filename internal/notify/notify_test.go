package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/snapshot"
)

func makeEvent(kind string, port uint16) notify.Event {
	return notify.Event{
		Kind: kind,
		Port: snapshot.Port{
			Port:     port,
			Protocol: "tcp",
			Process:  "test",
		},
		Host:      "localhost",
		Timestamp: time.Now(),
	}
}

func TestWebhookSender_SendOpened(t *testing.T) {
	var received notify.Event

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("unexpected content-type: %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	sender := notify.NewWebhookSender(srv.URL, 5*time.Second)
	event := makeEvent("opened", 8080)

	if err := sender.Send(event); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}
	if received.Kind != "opened" {
		t.Errorf("expected kind=opened, got %q", received.Kind)
	}
	if received.Port.Port != 8080 {
		t.Errorf("expected port 8080, got %d", received.Port.Port)
	}
}

func TestWebhookSender_NonSuccessStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	sender := notify.NewWebhookSender(srv.URL, 5*time.Second)
	err := sender.Send(makeEvent("closed", 443))
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}

func TestWebhookSender_DefaultTimeout(t *testing.T) {
	sender := notify.NewWebhookSender("http://example.com", 0)
	if sender.Timeout != 10*time.Second {
		t.Errorf("expected default timeout 10s, got %v", sender.Timeout)
	}
}

func TestWebhookSender_UnreachableHost(t *testing.T) {
	sender := notify.NewWebhookSender("http://127.0.0.1:1", 500*time.Millisecond)
	err := sender.Send(makeEvent("opened", 22))
	if err == nil {
		t.Fatal("expected error for unreachable host, got nil")
	}
}
