package host

import (
	"naksu/log"

	"github.com/intel-go/cpuid"
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

// GetWinProcessorAvailabilityLegend returns legend for Win32_Processor Availability
// code. It is used by getProcessorString (gethwlog_windows) and is implemented
// here in order to get tested
// See https://docs.microsoft.com/en-us/windows/win32/cimwin32prov/win32-processor
func GetWinProcessorAvailabilityLegend(legendCode uint16) string {
  legends := [22]string{
    "N/A",
    "Other",
    "Unknown",
    "Running/Full Power",
    "Warning",
    "In Test",
    "Not Applicable",
    "Power Off",
    "Off Line",
    "Off Duty",
    "Degraded",
    "Not Installed",
    "Install Error",
    "Power Save - Unknown",
    "Power Save - Low Power Mode",
    "Power Save - Standby",
    "Power Cycle",
    "Power Save - Warning",
    "Paused",
    "Not Ready",
    "Not Configured",
    "Quiesced",
  }

  if int(legendCode) > len(legends) {
    legendCode = 0
  }

  return legends[legendCode]
}
