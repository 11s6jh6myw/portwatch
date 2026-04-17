package schedule

import (
	"testing"
	"time"
)

func cfg() Config {
	return Config{
		MinInterval: 5 * time.Second,
		MaxInterval: 60 * time.Second,
		StepDown:    0.5,
		StepUp:      2.0,
	}
}

func TestNew_StartsAtMax(t *testing.T) {
	s := New(cfg())
	if s.Current() != 60*time.Second {
		t.Fatalf("expected 60s, got %v", s.Current())
	}
}

func TestRecordActivity_ShrinksInterval(t *testing.T) {
	s := New(cfg())
	s.RecordActivity()
	if s.Current() != 30*time.Second {
		t.Fatalf("expected 30s, got %v", s.Current())
	}
}

func TestRecordActivity_ClampsToMin(t *testing.T) {
	s := New(cfg())
	for i := 0; i < 20; i++ {
		s.RecordActivity()
	}
	if s.Current() != 5*time.Second {
		t.Fatalf("expected min 5s, got %v", s.Current())
	}
}

func TestRecordQuiet_GrowsInterval(t *testing.T) {
	s := New(cfg())
	s.RecordActivity() // 30s
	s.RecordQuiet()    // 60s
	if s.Current() != 60*time.Second {
		t.Fatalf("expected 60s, got %v", s.Current())
	}
}

func TestRecordQuiet_ClampsToMax(t *testing.T) {
	s := New(cfg())
	for i := 0; i < 10; i++ {
		s.RecordQuiet()
	}
	if s.Current() != 60*time.Second {
		t.Fatalf("expected max 60s, got %v", s.Current())
	}
}

func TestDefaultConfig_Valid(t *testing.T) {
	c := DefaultConfig()
	if c.MinInterval >= c.MaxInterval {
		t.Fatal("min must be less than max")
	}
	if c.StepDown >= 1.0 {
		t.Fatal("StepDown must be < 1")
	}
	if c.StepUp <= 1.0 {
		t.Fatal("StepUp must be > 1")
	}
}
