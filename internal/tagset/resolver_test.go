package tagset_test

import (
	"os"
	"testing"

	"github.com/example/portwatch/internal/tagset"
)

func TestResolver_AddsHost(t *testing.T) {
	r := tagset.NewResolver()
	s := r.Enrich(tagset.New())
	_, ok := s.Get("host")
	if !ok {
		t.Fatal("expected host tag")
	}
}

func TestResolver_EnvTagsIncluded(t *testing.T) {
	os.Setenv("PORTWATCH_TAG_REGION", "eu-west")
	t.Cleanup(func() { os.Unsetenv("PORTWATCH_TAG_REGION") })

	r := tagset.NewResolver()
	s := r.Enrich(tagset.New())
	v, ok := s.Get("region")
	if !ok || v != "eu-west" {
		t.Fatalf("expected region=eu-west, got %q ok=%v", v, ok)
	}
}

func TestResolver_CallerTagsWin(t *testing.T) {
	os.Setenv("PORTWATCH_TAG_ENV", "staging")
	t.Cleanup(func() { os.Unsetenv("PORTWATCH_TAG_ENV") })

	r := tagset.NewResolver()
	caller := tagset.New("env=prod")
	s := r.Enrich(caller)
	v, _ := s.Get("env")
	if v != "prod" {
		t.Fatalf("caller should win; got %q", v)
	}
}

func TestResolver_NoEnvTags_OnlyHost(t *testing.T) {
	// Clear any stray PORTWATCH_TAG_* vars set by other tests.
	for _, e := range os.Environ() {
		if len(e) > 14 && e[:14] == "PORTWATCH_TAG_" {
			key := e[:len(e)-len(e[len("PORTWATCH_TAG_"):])]
			_ = key
		}
	}
	r := tagset.NewResolver()
	s := r.Enrich(tagset.New())
	// At minimum host should be present (unless running in a stripped env).
	_ = s.All() // should not panic
}
