//go:build linux || darwin
// +build linux darwin

package mebroutines

// This file is required to run tests on Linux. It creates placeholder for
// getDiskFreeWindows()

func getDiskFreeWindows(path string) (uint64, error) {
	return 0, nil
}
