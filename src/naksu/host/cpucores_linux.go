package host

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"naksu/log"
)

// GetCPUCoreCount returns number of CPU cores
func GetCPUCoreCount() (int, error) {
	cpuinfoFileContentBytearr, err := os.ReadFile("/proc/cpuinfo")

	if err != nil {
		log.Error("Could not open /proc/cpuinfo to detect number of CPU cores: %v", err)

		return 0, err
	}

	cpuinfoFileContent := string(cpuinfoFileContentBytearr)

	re := regexp.MustCompile(`cpu cores\s+: (\d+)`)
	result := re.FindStringSubmatch(cpuinfoFileContent)
	if result != nil {
		coresStr := result[1]
		cores, err := strconv.ParseUint(coresStr, 10, 64)

		return int(cores), err
	}

	log.Debug(`Could not detect number of CPU cores from /proc/cpuinfo. The file appears to miss lines with "cpu cores" strings. Complete dump of the file follows:`)
	log.Debug(cpuinfoFileContent)

	return 0, fmt.Errorf("could not detect number of cpu cores from /proc/cpuinfo")
}
