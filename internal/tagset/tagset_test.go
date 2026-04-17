package tagset_test

import (
	"testing"

	"github.com/example/portwatch/internal/tagset"
)

func TestNew_ValidTags(t *testing.T) {
	s := tagset.New("env=prod", "team=infra")
	v, ok := s.Get("env")
	if !ok || v != "prod" {
		t.Fatalf("expected env=prod, got %q ok=%v", v, ok)
	}
}

func TestNew_MalformedEntryIgnored(t *testing.T) {
	s := tagset.New("nodash", "=empty", "good=yes")
	if len(s.All()) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(s.All()))
	}
}

func TestAdd_OverwritesExisting(t *testing.T) {
	s := tagset.New("env=dev")
	s.Add("env", "prod")
	v, _ := s.Get("env")
	if v != "prod" {
		t.Fatalf("expected prod, got %q", v)
	}
}

func TestAll_SortedByKey(t *testing.T) {
	s := tagset.New("z=last", "a=first", "m=mid")
	tags := s.All()
	if tags[0].Key != "a" || tags[1].Key != "m" || tags[2].Key != "z" {
		t.Fatalf("unexpected order: %v", tags)
	}
}

func TestString_CommaJoined(t *testing.T) {
	s := tagset.New("b=2", "a=1")
	if s.String() != "a=1,b=2" {
		t.Fatalf("unexpected string: %q", s.String())
	}
}

func TestMerge_OtherWinsOnConflict(t *testing.T) {
	a := tagset.New("env=dev", "team=a")
	b := tagset.New("env=prod", "region=us")
	m := a.Merge(b)
	v, _ := m.Get("env")
	if v != "prod" {
		t.Fatalf("expected prod, got %q", v)
	}
	if _, ok := m.Get("team"); !ok {
		t.Fatal("expected team tag from a")
	}
	if _, ok := m.Get("region"); !ok {
		t.Fatal("expected region tag from b")
	}
}

func TestGet_MissingKey(t *testing.T) {
	s := tagset.New()
	_, ok := s.Get("missing")
	if ok {
		t.Fatal("expected not found")
	}
}
