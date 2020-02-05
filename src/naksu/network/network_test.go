package network

import (
	"testing"
)

type IgnoredExtInterfaceTestData = []struct {
	extNicSystemName string
	ignore           bool
}

func TestIgnoreExtInterfaceLinux(t *testing.T) {
	testIgnoreExtInterface(t, IgnoredExtInterfaceTestData{
		{"lo", true},
		{"vboxnet15", true},
		{"loremipsumlo", false},
		{"eth1", false},
	}, isIgnoredExtInterfaceLinux)
}

func TestIgnoreExtInterfaceWindows(t *testing.T) {
	testIgnoreExtInterface(t, IgnoredExtInterfaceTestData{
		{"lo", false},
		{"loremipsumlo", false},
		{"VirtualBox Host-Only Ethernet Adapter", true},
	}, isIgnoredExtInterfaceWindows)
}

func TestIgnoreExtInterfaceDarwin(t *testing.T) {
	testIgnoreExtInterface(t, IgnoredExtInterfaceTestData{
		{"lo", true},
		{"gif666", true},
		{"vboxnet124", true},
		{"bridge0", true},
		{"loremipsumlo", false},
		{"en0", false},
	}, isIgnoredExtInterfaceDarwin)
}

func testIgnoreExtInterface(t *testing.T, testData IgnoredExtInterfaceTestData, functionUnderTest func(string) bool) {
	for _, table := range testData {
		ignore := functionUnderTest(table.extNicSystemName)
		if ignore != table.ignore {
			t.Errorf("IsIgnoredExtInterface fails with parameter '%s'", table.extNicSystemName)
		}
	}
}
