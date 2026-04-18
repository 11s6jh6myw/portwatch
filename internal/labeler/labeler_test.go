package labeler_test

import (
	"testing"

	"github.com/user/portwatch/internal/labeler"
)

func TestLabel_WellKnownPort(t *testing.T) {
	l := labeler.New(nil)
	if got := l.Label(22); got != "ssh" {
		t.Fatalf("expected ssh, got %s", got)
	}
}

func TestLabel_UnknownPort_ReturnsNumeric(t *testing.T) {
	l := labeler.New(nil)
	if got := l.Label(9999); got != "9999" {
		t.Fatalf("expected 9999, got %s", got)
	}
}

func TestLabel_OverrideTakesPrecedence(t *testing.T) {
	l := labeler.New(map[uint16]string{80: "my-app"})
	if got := l.Label(80); got != "my-app" {
		t.Fatalf("expected my-app, got %s", got)
	}
}

ideForUnknownPort(t *testing.T) {
	l := labeler.New(map[uint16]string{9000: "custom"})
	if got := l.Label(9000); got != "custom" {
		t.Fatalf("expected custom, got %s", got)
	}
}

func TestKnown_WellKnown(t *testing.T) {
	l := labeler.New(nil)
	if !l.Known(443) {
		t.Fatal("expected 443 to be known")
	}
}

func TestKnown_Unknown(t *testing.T) {
	l := labeler.New(nil)
	if l.Known(9999) {
		t.Fatal("expected 9999 to be unknown")
	}
}

func TestKnown_Override(t *testing.T) {
	l := labeler.New(map[uint16]string{1234: "svc"})
	if !l.Known(1234) {
		t.Fatal("expected 1234 to be known via override")
	}
}
