package oracledb

import "testing"

var versionMatchTests = []struct {
	oVersion    string
	lowVersion  string
	highVersion string
	expected    int
}{
	{oVersion: "11.2", lowVersion: "0", highVersion: "9999", expected: 1},
	{oVersion: "11.2", lowVersion: "0", highVersion: "11.2", expected: 1},
	{oVersion: "11.2", lowVersion: "11.2", highVersion: "9999", expected: 1},
	{oVersion: "11.2", lowVersion: "0", highVersion: "11.1", expected: 0},
	{oVersion: "11.1", lowVersion: "11.2", highVersion: "9999", expected: 0},
	{oVersion: "11.1", lowVersion: "0", highVersion: "10.9", expected: 0},
	{oVersion: "11.2", lowVersion: "10.1", highVersion: "9999", expected: 1},
}

func TestVersionMatch(t *testing.T) {
	for _, tt := range versionMatchTests {
		actual := VersionMatch(tt.oVersion, tt.lowVersion, tt.highVersion)
		if actual != tt.expected {
			t.Errorf(`VersionMatch(...) compare [oVersion] [H/L] [lowVersion] "%s" "%s" "%s"): expected "%d", actual "%d"`,
				tt.oVersion, tt.highVersion, tt.lowVersion, tt.expected, actual)
		}
	}
}
