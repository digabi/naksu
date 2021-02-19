package host

import (
	"runtime"
	"naksu/log"
)

// GetCPUCoreCount returns number of CPU cores
func GetCPUCoreCount() (int, error) {
	log.Debug("At the moment, CPU core count in Darwin is calcuated using runtime.NumCPU() which returns number of CPU threads")
	return runtime.NumCPU(), nil
}
