package filter_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/filter"
)

func TestNew_ValidRules(t *testing.T) {
	f, err := filter.New([]string{"22", "80", "8000-9000"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Rules()) != 3 {
		t.Fatalf("expected 3 rules, got %d", len(f.Rules()))
	}
}

func TestNew_InvalidPort(t *testing.T) {
	_, err := filter.New([]string{"abc"})
	if err == nil {
		t.Fatal("expected error for non-numeric port")
	}
}

func TestNew_InvalidRange_Reversed(t *testing.T) {
	_, err := filter.New([]string{"9000-8000"})
	if err == nil {
		t.Fatal("expected error for reversed range")
	}
}

func TestNew_ZeroPort(t *testing.T) {
	_, err := filter.New([]string{"0"})
	if err == nil {
		t.Fatal("expected error for port 0")
	}
}

func TestExcluded_ExactPort(t *testing.T) {
	f, _ := filter.New([]string{"22", "443"})
	if !f.Excluded(22) {
		t.Error("expected port 22 to be excluded")
	}
	if !f.Excluded(443) {
		t.Error("expected port 443 to be excluded")
	}
	if f.Excluded(80) {
		t.Error("expected port 80 to not be excluded")
	}
}

func TestExcluded_Range(t *testing.T) {
	f, _ := filter.New([]string{"8000-8100"})
	if !f.Excluded(8000) {
		t.Error("expected lower bound to be excluded")
	}
	if !f.Excluded(8050) {
		t.Error("expected mid-range port to be excluded")
	}
	if !f.Excluded(8100) {
		t.Error("expected upper bound to be excluded")
	}
	if f.Excluded(7999) {
		t.Error("expected port below range to not be excluded")
	}
	if f.Excluded(8101) {
		t.Error("expected port above range to not be excluded")
	}
}

func TestExcluded_EmptyFilter(t *testing.T) {
	f, _ := filter.New(nil)
	if f.Excluded(80) {
		t.Error("empty filter should not exclude any port")
	}
}
