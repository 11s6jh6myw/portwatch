package portmeta

import (
	"testing"
)

func TestReachabilityLevel_String(t *testing.T) {
	cases := []struct {
		level ReachabilityLevel
		want  string
	}{
		{ReachabilityNone, "none"},
		{ReachabilityPrivate, "private"},
		{ReachabilityLimited, "limited"},
		{ReachabilityPublic, "public"},
		{ReachabilityLevel(99), "unknown"},
	}
	for _, tc := range cases {
		t.Run(tc.want, func(t *testing.T) {
			if got := tc.level.String(); got != tc.want {
				t.Errorf("String() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestReachabilityFor_PublicPort(t *testing.T) {
	for _, port := range []int{80, 443, 22, 8080} {
		if got := ReachabilityFor(port); got != ReachabilityPublic {
			t.Errorf("port %d: got %v, want public", port, got)
		}
	}
}

func TestReachabilityFor_PrivatePort(t *testing.T) {
	for _, port := range []int{5432, 3306, 6379, 27017} {
		if got := ReachabilityFor(port); got != ReachabilityPrivate {
			t.Errorf("port %d: got %v, want private", port, got)
		}
	}
}

func TestReachabilityFor_LimitedPort(t *testing.T) {
	for _, port := range []int{3389, 5900, 9090} {
		if got := ReachabilityFor(port); got != ReachabilityLimited {
			t.Errorf("port %d: got %v, want limited", port, got)
		}
	}
}

func TestReachabilityFor_UnknownPort(t *testing.T) {
	if got := ReachabilityFor(19999); got != ReachabilityNone {
		t.Errorf("got %v, want none", got)
	}
}

func TestIsReachable_True(t *testing.T) {
	if !IsReachable(80) {
		t.Error("expected port 80 to be reachable")
	}
}

func TestIsReachable_False(t *testing.T) {
	if IsReachable(19999) {
		t.Error("expected port 19999 to be not reachable")
	}
}
