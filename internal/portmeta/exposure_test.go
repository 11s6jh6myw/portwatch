package portmeta

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func TestExposureLevel_String(t *testing.T) {
	cases := []struct {
		level ExposureLevel
		want  string
	}{
		{ExposureNone, "none"},
		{ExposureLow, "low"},
		{ExposureMedium, "medium"},
		{ExposureHigh, "high"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("ExposureLevel(%d).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}

func TestExposureFor_HighExposure(t *testing.T) {
	for _, port := range []int{22, 80, 443, 3306, 6379} {
		if got := ExposureFor(port); got != ExposureHigh {
			t.Errorf("ExposureFor(%d) = %v, want high", port, got)
		}
	}
}

func TestExposureFor_MediumExposure(t *testing.T) {
	for _, port := range []int{445, 5900, 9200} {
		if got := ExposureFor(port); got != ExposureMedium {
			t.Errorf("ExposureFor(%d) = %v, want medium", port, got)
		}
	}
}

func TestExposureFor_LowExposure(t *testing.T) {
	// port 512 is in the privileged range but not in high/medium sets
	if got := ExposureFor(512); got != ExposureLow {
		t.Errorf("ExposureFor(512) = %v, want low", got)
	}
}

func TestExposureFor_NoExposure(t *testing.T) {
	if got := ExposureFor(54321); got != ExposureNone {
		t.Errorf("ExposureFor(54321) = %v, want none", got)
	}
}

func TestFilterByMinExposure(t *testing.T) {
	annotator := NewExposureAnnotator()
	ports := []scanner.PortInfo{
		{Port: 22},
		{Port: 445},
		{Port: 54321},
	}
	annotated := annotator.Annotate(ports)
	result := FilterByMinExposure(annotated, ExposureMedium)
	if len(result) != 2 {
		t.Fatalf("expected 2 ports with exposure >= medium, got %d", len(result))
	}
}
