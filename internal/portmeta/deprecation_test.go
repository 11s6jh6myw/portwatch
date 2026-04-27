package portmeta

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func TestDeprecationLevel_String(t *testing.T) {
	cases := []struct {
		level DeprecationLevel
		want  string
	}{
		{DeprecationNone, "none"},
		{DeprecationMinor, "minor"},
		{DeprecationModerate, "moderate"},
		{DeprecationHigh, "high"},
		{DeprecationObsolete, "obsolete"},
		{DeprecationLevel(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("DeprecationLevel(%d).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}

func TestDeprecationFor_ObsoletePort(t *testing.T) {
	// Telnet (port 23) deprecated 1995 — well over 25 years ago
	p := scanner.PortInfo{Port: 23, Proto: "tcp"}
	got := DeprecationFor(p)
	if got < DeprecationHigh {
		t.Errorf("DeprecationFor(telnet) = %v, want >= high", got)
	}
}

func TestDeprecationFor_ActivePort(t *testing.T) {
	// Port 443 (HTTPS) is not deprecated
	p := scanner.PortInfo{Port: 443, Proto: "tcp"}
	got := DeprecationFor(p)
	if got != DeprecationNone {
		t.Errorf("DeprecationFor(https) = %v, want none", got)
	}
}

func TestDeprecationFor_UnknownPort(t *testing.T) {
	p := scanner.PortInfo{Port: 59999, Proto: "tcp"}
	if got := DeprecationFor(p); got != DeprecationNone {
		t.Errorf("DeprecationFor(unknown) = %v, want none", got)
	}
}

func TestIsDeprecated_True(t *testing.T) {
	p := scanner.PortInfo{Port: 21, Proto: "tcp"} // FTP
	if !IsDeprecated(p) {
		t.Error("expected FTP (21) to be deprecated")
	}
}

func TestIsDeprecated_False(t *testing.T) {
	p := scanner.PortInfo{Port: 22, Proto: "tcp"} // SSH
	if IsDeprecated(p) {
		t.Error("expected SSH (22) not to be deprecated")
	}
}

func TestDeprecationFor_FTP_IsHighOrObsolete(t *testing.T) {
	// FTP deprecated ~1990, 30+ years ago
	p := scanner.PortInfo{Port: 21, Proto: "tcp"}
	got := DeprecationFor(p)
	if got < DeprecationHigh {
		t.Errorf("DeprecationFor(ftp) = %v, want high or obsolete", got)
	}
}
