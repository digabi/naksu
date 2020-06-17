package host

import (
	"fmt"
	"strings"
	"sort"

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

type Win32_ComputerSystem struct { //nolint
	// Win32_ComputerSystem is a struct used to query Windows WMI
	// (Windows Management Instrumentation)
	// The struct must be named with an underscore, otherwise it is not recognised
	// and results "Invalid class" exception.
	TotalPhysicalMemory uint64
}

// Win32_PnPEntity is a struct used to query Windows WMI
// (Windows Management Instrumentation)
// The struct must be named with an underscore, otherwise it is not recognised
// and results "Invalid class" exception.
type Win32_PnPEntity struct { //nolint
	PNPClass     string
	Manufacturer string
	Name         string
	DeviceID     string
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

func getPnpEntityData() []Win32_PnPEntity {
	var dst []Win32_PnPEntity
	query := wmi.CreateQuery(&dst, "")
	err := wmi.Query(query, &dst)
	if err != nil {
		log.Debug(fmt.Sprintf("getPnpEntityData() could not make WMI query: %v", err))
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

func getPnpEntityString() string {
	pnpEntityData := getPnpEntityData()

	var pnpEntities []string

	for thisEntity := range pnpEntityData {
		pnpEntities = append(pnpEntities,
			fmt.Sprintf(
				"%s %s %s [%s]",
				pnpEntityData[thisEntity].PNPClass,
				pnpEntityData[thisEntity].Manufacturer,
				pnpEntityData[thisEntity].Name,
				pnpEntityData[thisEntity].DeviceID,
			),
		)
	}

	// Sort in alphabetical order
	sort.Strings(pnpEntities)

	return strings.Join(pnpEntities[:], "\n")
}

// GetHwLog returns a single string containing various hardware
// information to be printed to a log file
func GetHwLog() string {
	cpuinfo := getProcessorString()
	memoryinfo := getMemoryString()
	pnpentities := getPnpEntityString()

	return fmt.Sprintf(`Processor Info
%s
Memory Info
%s
Plug-And-Play Devices
%s`, cpuinfo, memoryinfo, pnpentities)
}
