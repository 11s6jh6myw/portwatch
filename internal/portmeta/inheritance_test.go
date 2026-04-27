package portmeta_test

import (
	"testing"

	"github.com/user/portwatch/internal/portmeta"
	"github.com/user/portwatch/internal/scanner"
)

func makeInhPort(port int) scanner.PortInfo {
	return scanner.PortInfo{Port: port}
}

func TestInheritanceLevel_String(t *testing.T) {
	cases := []struct {
		lvl  portmeta.InheritanceLevel
		want string
	}{
		{portmeta.InheritanceNone, "none"},
		{portmeta.InheritanceDerived, "derived"},
		{portmeta.InheritanceExtended, "extended"},
		{portmeta.InheritanceCore, "core"},
	}
	for _, tc := range cases {
		if got := tc.lvl.String(); got != tc.want {
			t.Errorf("String() = %q, want %q", got, tc.want)
		}
	}
}

func TestInheritanceFor_CorePort(t *testing.T) {
	p := makeInhPort(80)
	if got := portmeta.InheritanceFor(p); got != portmeta.InheritanceCore {
		t.Errorf("port 80: got %v, want core", got)
	}
}

func TestInheritanceFor_ExtendedPort(t *testing.T) {
	p := makeInhPort(8080)
	if got := portmeta.InheritanceFor(p); got != portmeta.InheritanceExtended {
		t.Errorf("port 8080: got %v, want extended", got)
	}
}

func TestInheritanceFor_DerivedPort(t *testing.T) {
	p := makeInhPort(3306)
	if got := portmeta.InheritanceFor(p); got != portmeta.InheritanceDerived {
		t.Errorf("port 3306: got %v, want derived", got)
	}
}

func TestInheritanceFor_UnknownPort(t *testing.T) {
	p := makeInhPort(19999)
	if got := portmeta.InheritanceFor(p); got != portmeta.InheritanceNone {
		t.Errorf("port 19999: got %v, want none", got)
	}
}

func TestIsInherited_True(t *testing.T) {
	if !portmeta.IsInherited(makeInhPort(443)) {
		t.Error("expected port 443 to be inherited")
	}
}

func TestIsInherited_False(t *testing.T) {
	if portmeta.IsInherited(makeInhPort(54321)) {
		t.Error("expected port 54321 to not be inherited")
	}
}
