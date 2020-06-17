package host

import (
	"fmt"
	"strings"

	"naksu/log"

	"github.com/StackExchange/wmi"
	humanize "github.com/dustin/go-humanize"
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

// Win32_ComputerSystem is a struct used to query Windows WMI
// (Windows Management Instrumentation)
// The struct must be named with an underscore, otherwise it is not recognised
// and results "Invalid class" exception.
type Win32_ComputerSystem struct { //nolint
	TotalPhysicalMemory uint64
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

func getMemoryData() []Win32_ComputerSystem {
	var dst []Win32_ComputerSystem
	query := wmi.CreateQuery(&dst, "")
	err := wmi.Query(query, &dst)
	if err != nil {
		log.Debug(fmt.Sprintf("getMemoryData() could not make WMI query: %v", err))
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

func getMemoryString() string {
	memoryData := getMemoryData()

	var totalPhysicalMemory uint64
	totalPhysicalMemory = 0

	for thisMemoryRecord := range memoryData {
		if memoryData[thisMemoryRecord].TotalPhysicalMemory > totalPhysicalMemory {
			totalPhysicalMemory = memoryData[thisMemoryRecord].TotalPhysicalMemory
		}
	}

	return fmt.Sprintf("Total Memory: %d (%s)", totalPhysicalMemory, humanize.Bytes(totalPhysicalMemory))
}

// GetHwLog returns a single string containing various hardware
// information to be printed to a log file
func GetHwLog() string {
	cpuinfo := getProcessorString()
	memoryinfo := getMemoryString()

	return fmt.Sprintf(`Processor Info
%s
Memory Info
%s`, cpuinfo, memoryinfo)
}
