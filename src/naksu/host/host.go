package host

import (
	"runtime"

	"naksu/log"

	"github.com/intel-go/cpuid"
	"github.com/mackerelio/go-osstat/memory"
)

// host can be used to get information of the host machine

// IsHWVirtualisationCPU returns true if CPU supports hardware virtualisation
// This does not detect whether the support is turned in BIOS
// See IsHWVirtualisation()
func IsHWVirtualisationCPU() bool {
	if cpuid.HasFeature(cpuid.VMX) {
		log.Debug("Hardware virtualisation is supported by CPU (VT-x, CPU flag VMX)")
		return true
	}

	if cpuid.HasExtraFeature(cpuid.SVM) {
		log.Debug("Hardware virtualisation is supported by CPU (AMD-V, CPU flag SVM)")
		return true
	}

	log.Debug("Hardware virtualisation is not supported by CPU")
	return false
}

// GetCPUCoreCount returns number of CPU cores
func GetCPUCoreCount() int {
	return runtime.NumCPU()
}

// GetMemory returns system RAM (in megabytes)
func GetMemory() (uint64, error) {
	memory, err := memory.Get()
	if err != nil {
		return 0, err
	}

	return memory.Total / (1024 * 1024), nil
}
