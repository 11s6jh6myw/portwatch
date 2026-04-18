package portmeta_test

import (
	"testing"

	"github.com/example/portwatch/internal/portmeta"
)

func TestCategorize_KnownPorts(t *testing.T) {
	cases := []struct {
		port uint16
		want portmeta.Category
	}{
		{80, portmeta.CategoryWeb},
		{443, portmeta.CategoryWeb},
		{8080, portmeta.CategoryWeb},
		{22, portmeta.CategoryRemote},
		{3389, portmeta.CategoryRemote},
		{3306, portmeta.CategoryDatabase},
		{5432, portmeta.CategoryDatabase},
		{6379, portmeta.CategoryDatabase},
		{27017, portmeta.CategoryDatabase},
		{25, portmeta.CategoryMail},
		{53, portmeta.CategoryDNS},
		{9092, portmeta.CategoryMessaging},
	}
	for _, tc := range cases {
		t.Run(tc.want.String(), func(t *testing.T) {
			got := portmeta.Categorize(tc.port)
			if got != tc.want {
				t.Errorf("port %d: got %s, want %s", tc.port, got, tc.want)
			}
		})
	}
}

func TestCategorize_UnknownPort(t *testing.T) {
	got := portmeta.Categorize(12345)
	if got != portmeta.CategoryUnknown {
		t.Errorf("expected unknown, got %s", got)
	}
}

func TestCategory_String(t *testing.T) {
	if portmeta.CategoryWeb.String() != "web" {
		t.Errorf("unexpected string for CategoryWeb")
	}
}

func TestCategorizeAll_ReturnsMappedCategories(t *testing.T) {
	ports := []uint16{80, 443, 9999}
	result := portmeta.CategorizeAll(ports)
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	if result[80] != portmeta.CategoryWeb {
		t.Errorf("port 80: expected web, got %s", result[80])
	}
	if result[9999] != portmeta.CategoryUnknown {
		t.Errorf("port 9999: expected unknown, got %s", result[9999])
	}
}
