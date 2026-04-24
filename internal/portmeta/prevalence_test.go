package portmeta

import (
	"testing"
	"time"
)

func TestPrevalenceLevel_String(t *testing.T) {
	cases := []struct {
		level PrevalenceLevel
		want  string
	}{
		{PrevalenceUnknown, "unknown"},
		{PrevalenceRare, "rare"},
		{PrevalenceUncommon, "uncommon"},
		{PrevalenceCommon, "common"},
		{PrevalenceUbiquitous, "ubiquitous"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("String() = %q, want %q", got, tc.want)
		}
	}
}

func TestPrevalenceFor_UbiquitousPort(t *testing.T) {
	for _, port := range []int{22, 80, 443} {
		if got := PrevalenceFor(port, time.Time{}); got != PrevalenceUbiquitous {
			t.Errorf("port %d: got %v, want ubiquitous", port, got)
		}
	}
}

func TestPrevalenceFor_CommonPort(t *testing.T) {
	if got := PrevalenceFor(3306, time.Time{}); got != PrevalenceCommon {
		t.Errorf("got %v, want common", got)
	}
}

func TestPrevalenceFor_UncommonPort(t *testing.T) {
	old := time.Now().Add(-48 * time.Hour)
	if got := PrevalenceFor(9200, old); got != PrevalenceUncommon {
		t.Errorf("got %v, want uncommon", got)
	}
}

func TestPrevalenceFor_UncommonPort_RecentlySeenIsRare(t *testing.T) {
	recent := time.Now().Add(-1 * time.Hour)
	if got := PrevalenceFor(9200, recent); got != PrevalenceRare {
		t.Errorf("got %v, want rare for recently-seen uncommon port", got)
	}
}

func TestPrevalenceFor_UnknownPort(t *testing.T) {
	if got := PrevalenceFor(19999, time.Time{}); got != PrevalenceRare {
		t.Errorf("got %v, want rare", got)
	}
}

func TestIsPrevalent_True(t *testing.T) {
	if !IsPrevalent(PrevalenceCommon) {
		t.Error("expected common to be prevalent")
	}
	if !IsPrevalent(PrevalenceUbiquitous) {
		t.Error("expected ubiquitous to be prevalent")
	}
}

func TestIsPrevalent_False(t *testing.T) {
	if IsPrevalent(PrevalenceRare) {
		t.Error("expected rare to not be prevalent")
	}
	if IsPrevalent(PrevalenceUncommon) {
		t.Error("expected uncommon to not be prevalent")
	}
}
