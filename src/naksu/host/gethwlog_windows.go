package host

import (
	"fmt"
	"sort"
	"strings"

	"naksu/log"

	humanize "github.com/dustin/go-humanize"
	"github.com/yusufpapurcu/wmi"
)

// Win32_Processor is a struct used to query Windows WMI
// (Windows Management Instrumentation)
// The struct must be named with an underscore, otherwise it is not recognised
// and results "Invalid class" exception.
type Win32_Processor struct { //nolint
	Availability      *uint16
	Caption           *string
	CurrentClockSpeed *uint32
	Description       *string
	DeviceID          *string
	LoadPercentage    *uint16
	Manufacturer      *string
	MaxClockSpeed     *uint32
	Name              *string
	NumberOfCores     *uint32
}

type Win32_ComputerSystem struct { //nolint
	// Win32_ComputerSystem is a struct used to query Windows WMI
	// (Windows Management Instrumentation)
	// The struct must be named with an underscore, otherwise it is not recognised
	// and results "Invalid class" exception.
	TotalPhysicalMemory *uint64
}

// Win32_PnPEntity is a struct used to query Windows WMI
// (Windows Management Instrumentation)
// The struct must be named with an underscore, otherwise it is not recognised
// and results "Invalid class" exception.
type Win32_PnPEntity struct { //nolint
	PNPClass     *string
	Manufacturer *string
	Name         *string
	DeviceID     *string
}

func getProcessorData() []Win32_Processor {
	result := make(chan []Win32_Processor)

	// Do this in Goroutine to avoid "cannot change thread mode" in Windows WMI call
	go func() {
		var dst []Win32_Processor
		query := wmi.CreateQuery(&dst, "")
		err := wmi.Query(query, &dst)
		if err != nil {
			log.Error("getProcessorData() could not make WMI query: %v", err)
			result <- []Win32_Processor{}
		} else {
			result <- dst
		}
	}()

	return <-result
}

func getMemoryData() []Win32_ComputerSystem {
	result := make(chan []Win32_ComputerSystem)

	// Do this in Goroutine to avoid "cannot change thread mode" in Windows WMI call
	go func() {
		var dst []Win32_ComputerSystem
		query := wmi.CreateQuery(&dst, "")
		err := wmi.Query(query, &dst)
		if err != nil {
			log.Error("getMemoryData() could not make WMI query: %v", err)
			result <- []Win32_ComputerSystem{}
		} else {
			result <- dst
		}
	}()

	return <-result
}

func getPnpEntityData() []Win32_PnPEntity {
	result := make(chan []Win32_PnPEntity)

	// Do this in Goroutine to avoid "cannot change thread mode" in Windows WMI call
	go func() {
		var dst []Win32_PnPEntity
		query := wmi.CreateQuery(&dst, "")
		err := wmi.Query(query, &dst)
		if err != nil {
			log.Error("getPnpEntityData() could not make WMI query: %v", err)
			result <- []Win32_PnPEntity{}
		} else {
			result <- dst
		}
	}()

	return <-result
}

func getProcessorString() string {
	processorData := getProcessorData()

	var processorInfo []string

	for thisProcessor := range processorData {
		processorInfo = append(processorInfo,
			fmt.Sprintf(
				"%s: %s, %s, %s Availability: %s, CurrentClockSpeed: %d, MaxClockSpeed: %d, NumberOfCores: %d",
				*processorData[thisProcessor].DeviceID,
				*processorData[thisProcessor].Manufacturer,
				*processorData[thisProcessor].Name,
				*processorData[thisProcessor].Caption,
				getWinProcessorAvailabilityLegend(*processorData[thisProcessor].Availability),
				*processorData[thisProcessor].CurrentClockSpeed,
				*processorData[thisProcessor].MaxClockSpeed,
				*processorData[thisProcessor].NumberOfCores,
			),
		)
	}

	return strings.Join(processorInfo, "\n")
}

func getMemoryString() string {
	memoryData := getMemoryData()

	var totalPhysicalMemory uint64

	for thisMemoryRecord := range memoryData {
		if *memoryData[thisMemoryRecord].TotalPhysicalMemory > totalPhysicalMemory {
			totalPhysicalMemory = *memoryData[thisMemoryRecord].TotalPhysicalMemory
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
				*pnpEntityData[thisEntity].PNPClass,
				*pnpEntityData[thisEntity].Manufacturer,
				*pnpEntityData[thisEntity].Name,
				*pnpEntityData[thisEntity].DeviceID,
			),
		)
	}

	// Sort in alphabetical order
	sort.Strings(pnpEntities)

	return strings.Join(pnpEntities, "\n")
}

// getWinProcessorAvailabilityLegend returns legend for Win32_Processor Availability
// code. It is used by getProcessorString (gethwlog_windows) and is implemented
// here in order to get tested
// See https://docs.microsoft.com/en-us/windows/win32/cimwin32prov/win32-processor
func getWinProcessorAvailabilityLegend(legendCode uint16) string {
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

// GetHwLog returns a single string containing various hardware
// information to be printed to a log file
func GetHwLog() string {
	cpuInfo := getProcessorString()
	powerplan := getPowerplan()
	memoryInfo := getMemoryString()
	pnpEntities := getPnpEntityString()

	return fmt.Sprintf(`
===== Processor Info
%s

===== Power configuration
%s

===== Memory Info
%s

===== Plug-And-Play Devices
%s

===== End of Hardware Log
`, cpuInfo, powerplan, memoryInfo, pnpEntities)
}
