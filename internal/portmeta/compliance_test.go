package portmeta

import "testing"

func TestComplianceFlagsFor_KnownPort(t *testing.T) {
	flags := ComplianceFlagsFor(3306)
	if len(flags) == 0 {
		t.Fatal("expected compliance flags for port 3306")
	}
}

func TestComplianceFlagsFor_UnknownPort(t *testing.T) {
	flags := ComplianceFlagsFor(9999)
	if len(flags) != 1 || flags[0] != ComplianceNone {
		t.Errorf("expected [none], got %v", flags)
	}
}

func TestHasComplianceFlag_True(t *testing.T) {
	if !HasComplianceFlag(443, CompliancePCI) {
		t.Error("expected port 443 to have PCI flag")
	}
}

func TestHasComplianceFlag_False(t *testing.T) {
	if HasComplianceFlag(443, ComplianceHIPAA) {
		t.Error("expected port 443 not to have HIPAA flag")
	}
}

func TestHasComplianceFlag_TelnetAllFrameworks(t *testing.T) {
	for _, flag := range []ComplianceFlag{CompliancePCI, ComplianceHIPAA, ComplianceSOC2} {
		if !HasComplianceFlag(23, flag) {
			t.Errorf("expected port 23 to have flag %s", flag)
		}
	}
}

func TestComplianceFlag_Values(t *testing.T) {
	cases := []struct {
		flag ComplianceFlag
		want string
	}{
		{CompliancePCI, "pci"},
		{ComplianceHIPAA, "hipaa"},
		{ComplianceSOC2, "soc2"},
		{ComplianceNone, "none"},
	}
	for _, tc := range cases {
		if string(tc.flag) != tc.want {
			t.Errorf("flag = %q, want %q", tc.flag, tc.want)
		}
	}
}
