package host

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"naksu/log"
)

const linuxCPUSettingsPath = "/sys/devices/system/cpu/"

func getPowerplanFilenames() []string {
	var powerplanFilenames []string

	dir, err := os.Open(linuxCPUSettingsPath)
	if err != nil {
		log.Error("Could not open path %s to read power plan filenames: %v", linuxCPUSettingsPath, err)

		return powerplanFilenames
	}

	defer dir.Close()

	files, err := dir.ReadDir(-1)
	if err != nil {
		log.Error("Could not read filenames in directory %s to read power plan filenames: %v", linuxCPUSettingsPath, err)

		return powerplanFilenames
	}

	re := regexp.MustCompile(`^cpu\d+`)

	for _, file := range files {
		if re.MatchString(file.Name()) {
			powerplanFilenames = append(powerplanFilenames, fmt.Sprintf("%s%s/cpufreq/scaling_governor", linuxCPUSettingsPath, file.Name()))
		}
	}

	return powerplanFilenames
}

func getPowerplanString(filename string) string {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		log.Error("Could not open power plan file %s for reading: %v", filename, err)

		return "error"
	}

	reOnlyletters := regexp.MustCompile(`\W`)

	return string(reOnlyletters.ReplaceAll(fileContent, []byte("")))
}

// getPowerplan returns a string describing current power plan
func getPowerplan() string {
	powerplanFilenames := getPowerplanFilenames()

	var powerplans []string

	for _, powerplanFile := range powerplanFilenames {
		powerplans = append(powerplans, getPowerplanString(powerplanFile))
	}

	return strings.Join(powerplans, "+")
}
