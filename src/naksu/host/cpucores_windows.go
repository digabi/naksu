package host

// GetCPUCoreCount returns number of CPU cores
func GetCPUCoreCount() (int, error) {
	processorData := getProcessorData()

	var coreCount uint32

	for thisProcessor := range processorData {
		coreCount += *processorData[thisProcessor].NumberOfCores
	}

	return int(coreCount), nil
}
