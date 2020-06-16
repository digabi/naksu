package host

import (
	"fmt"
	"strings"

	"naksu/log"

	"github.com/StackExchange/wmi"
)

// Win32_Processor is a struct used to query Windows WMI
// (Windows Management Instrumentation)
// The struct must be named with an underscore, otherwise it is not recognised
// and results "Invalid class" exception.
type Win32_Processor struct { //nolint
	Availability      uint16
	Caption           string
	CurrentClockSpeed uint32
	Description       string
	DeviceID          string
	LoadPercentage    uint16
	Manufacturer      string
	MaxClockSpeed     uint32
	Name              string
}

func getProcessorData() []Win32_Processor {
	var dst []Win32_Processor
	query := wmi.CreateQuery(&dst, "")
	err := wmi.Query(query, &dst)
	if err != nil {
		log.Debug(fmt.Sprintf("getProcessorData() could not make WMI query: %v", err))
	}

	return dst
}

func getProcessorString() string {
	processorData := getProcessorData()

	var processorInfo []string

	for thisProcessor := range processorData {
		processorInfo = append(processorInfo,
			fmt.Sprintf(
				"%s: %s, %s, %s Availability: %s, CurrentClockSpeed: %d, MaxClockSpeed: %d",
				processorData[thisProcessor].DeviceID,
				processorData[thisProcessor].Manufacturer,
				processorData[thisProcessor].Name,
				processorData[thisProcessor].Caption,
        GetWinProcessorAvailabilityLegend(processorData[thisProcessor].Availability),
				processorData[thisProcessor].CurrentClockSpeed,
				processorData[thisProcessor].MaxClockSpeed,
			),
		)
	}

	return strings.Join(processorInfo[:], "\n")
}

// GetHwLog returns a single string containing various hardware
// information to be printed to a log file
func GetHwLog() string {
	cpuinfo := getProcessorString()

	return fmt.Sprintf(`Processor Info
%s`, cpuinfo)
}
