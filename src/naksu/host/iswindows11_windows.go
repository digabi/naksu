package host

import (
	"naksu/log"

	"golang.org/x/sys/windows"
)

func IsWindows11() bool {
	maj, min, build := windows.RtlGetNtVersionNumbers()

	log.Debug("Windows version: %d.%d (build %d)\n",
		maj,
		min,
		build)

	// Windows 10 and 11 both have major version as 10 but Windows 11 has build number >= 22000
	if maj == 10 && build >= 22000 {
		log.Debug("This is likely Windows 11.")
		return true
	} else {
		log.Debug("This is not Windows 11.")
		return false
	}
}
