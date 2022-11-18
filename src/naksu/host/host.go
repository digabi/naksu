package host

import (
	"fmt"

	"naksu/box/vboxmanage"
	"naksu/constants"
	"naksu/log"
	"naksu/mebroutines"
	"naksu/xlate"

	"github.com/intel-go/cpuid"
	"github.com/mackerelio/go-osstat/memory"

	semver "github.com/blang/semver/v4"
	humanize "github.com/dustin/go-humanize"
)

// host can be used to get information of the host machine

// LowDiskSizeError is an error returned by CheckFreeDisk()
type LowDiskSizeError struct {
	Err     string
	LowPath string
	LowSize uint64
}

func (e *LowDiskSizeError) Error() string {
	return fmt.Sprintf("path %s has low disk size: %d (%s)", e.LowPath, e.LowSize, humanize.Bytes(e.LowSize))
}

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

// GetMemory returns system RAM (in megabytes)
func GetMemory() (uint64, error) {
	memory, err := memory.Get()
	if err != nil {
		return 0, err
	}

	const megabyteInBytes = 1024 * 1024

	return memory.Total / megabyteInBytes, nil
}

// CheckFreeDisk checks that all of the listed directories have more than
// given limit free disk space. If a directory has less than required disk space
// the returned error has prefix "low:" followed by a failed path. The uint64 returns free disk space
// of this location.
func CheckFreeDisk(limit uint64, directories []string) error {
	log.Debug("CheckFreeDisk: %v", directories)

	for _, thisDirectory := range directories {
		freeDisk, err := mebroutines.GetDiskFree(thisDirectory)

		if err != nil {
			log.Error("CheckFreeDisk could not get free disk for path '%s': %v", thisDirectory, err)
		} else {
			log.Debug("CheckFreeDisk: %s (%d bytes, %s)", thisDirectory, freeDisk, humanize.Bytes(freeDisk))

			if freeDisk < limit {
				return &LowDiskSizeError{"disk size is too low", thisDirectory, freeDisk}
			}
		}
	}

	return nil
}

// IsVirtualBoxVersionOK returns an user-formatted non-empty string if VirtualBox version
// is too low or too high (constants.VBoxMinVersion, constants.VBoxMaxVersion, respectively)
func IsVirtualBoxVersionOK() (string, error) {
	vBoxVersion, err := vboxmanage.GetVBoxManageVersion()
	if err != nil {
		log.Debug("Could not get VBoxManage version: %v", err)

		return "", fmt.Errorf("could not get vboxmanage version: %w", err)
	}

	if constants.VBoxMinVersion != "" {
		versionMin, err := semver.Parse(constants.VBoxMinVersion)
		if err != nil {
			return "", fmt.Errorf("could not parse minimum required version %s", constants.VBoxMinVersion)
		}

		if vBoxVersion.LT(versionMin) {
			return xlate.Get("Your VirtualBox version is old. Consider upgrading to %s or newer to avoid problems.", constants.VBoxMinVersion), nil
		}
	}

	if constants.VBoxMaxVersion != "" {
		versionMax, err := semver.Parse(constants.VBoxMaxVersion)
		if err != nil {
			return "", fmt.Errorf("could not parse maximum required version %s", constants.VBoxMaxVersion)
		}

		if vBoxVersion.GT(versionMax) {
			return xlate.Get("Your VirtualBox version is too new. Consider downgrading to %s to avoid problems.", constants.VBoxMaxVersion), nil
		}
	}

	return "", nil
}
