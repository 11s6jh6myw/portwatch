package portmeta

import (
	"testing"
)

func TestTagsFor_KnownPort(t *testing.T) {
	tags := TagsFor(80)
	if len(tags) == 0 {
		t.Fatal("expected tags for port 80")
	}
	if tags[0] != TagWeb {
		t.Errorf("expected TagWeb, got %s", tags[0])
	}
}

func TestTagsFor_UnknownPort(t *testing.T) {
	tags := TagsFor(9999)
	if len(tags) != 1 || tags[0] != TagUnknown {
		t.Errorf("expected [unknown], got %v", tags)
	}
}

func TestTagsFor_DatabasePort(t *testing.T) {
	for _, port := range []uint16{3306, 5432, 6379, 27017} {
		tags := TagsFor(port)
		found := false
		for _, tag := range tags {
			if tag == TagDatabase {
				found = true
			}
		}
		if !found {
			t.Errorf("port %d: expected TagDatabase in %v", port, tags)
		}
	}
}

func TestHasTag_True(t *testing.T) {
	if !HasTag(22, TagRemote) {
		t.Error("expected port 22 to have TagRemote")
	}
}

func TestHasTag_False(t *testing.T) {
	if HasTag(22, TagDatabase) {
		t.Error("port 22 should not have TagDatabase")
	}
}

func TestHasTag_UnknownPort(t *testing.T) {
	if !HasTag(9999, TagUnknown) {
		t.Error("expected unknown port to have TagUnknown")
	}
}
