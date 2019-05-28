package network_test

import (
	"naksu/network"
	"testing"
)

func TestIgnoreExtInterface(t *testing.T) {
	tables := []struct {
		extNicSystemName string
		ignore           bool
	}{
		{"lo", true},
		{"loremipsumlo", false},
		{"eth1", false},
	}

	for _, table := range tables {
		ignore := network.IgnoreExtInterface(table.extNicSystemName)
		if ignore != table.ignore {
			t.Errorf("IgnoreExtInterface fails with parameter '%s'", table.extNicSystemName)
		}
	}
}
