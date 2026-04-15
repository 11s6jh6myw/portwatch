// Package notify provides webhook and external notification delivery
// for portwatch port change events.
package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Event represents a port change notification payload.
type Event struct {
	Kind      string          `json:"kind"`       // "opened" or "closed"
	Port      snapshot.Port   `json:"port"`
	Host      string          `json:"host"`
	Timestamp time.Time       `json:"timestamp"`
}

// WebhookSender delivers port change events to an HTTP endpoint.
type WebhookSender struct {
	URL     string
	Timeout time.Duration
	client  *http.Client
}

// NewWebhookSender creates a WebhookSender targeting the given URL.
// If timeout is zero, a default of 10 seconds is used.
func NewWebhookSender(url string, timeout time.Duration) *WebhookSender {
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	return &WebhookSender{
		URL:     url,
		Timeout: timeout,
		client:  &http.Client{Timeout: timeout},
	}
}

// Send serialises the event as JSON and POSTs it to the configured URL.
// It returns an error if the request fails or the server responds with
// a non-2xx status code.
func (w *WebhookSender) Send(event Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("notify: marshal event: %w", err)
	}

	resp, err := w.client.Post(w.URL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("notify: post to %s: %w", w.URL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("notify: server returned %d for %s", resp.StatusCode, w.URL)
	}
	return nil
}
