package network

import (
	"fmt"

	"github.com/StackExchange/wmi"
	humanize "github.com/dustin/go-humanize"

	"naksu/constants"
	"naksu/log"
)

// GetExtInterfaces returns map of network interfaces. It returns array of
// constants.AvailableSelection where the ConfigValue is system's internal value
// (e.g. "Intel(R) Ethernet Connection (4) I219-LM") while the value is human-readable legend.
// In Windows these two values are the same.
func GetExtInterfaces() []constants.AvailableSelection {
	// This struct must be named with an underscore, otherwise it is not recognised
	// and results "Invalid class" exception.
	type Win32_NetworkAdapter struct { //nolint
		AdapterTypeID   uint16
		Name            string
		Speed           uint64
		PhysicalAdapter bool
		NetEnabled      bool
	}

	result := constants.DefaultExtNicArray

	var dst []Win32_NetworkAdapter
	query := wmi.CreateQuery(&dst, "WHERE PhysicalAdapter=TRUE")
	err := wmi.Query(query, &dst)
	if err != nil {
		log.Debug("GetExtInterfaces() could query network adapters from WMI")
		log.Debug(fmt.Sprint(err))
		return result
	}

	formatLinkSpeed := func(networkInterface *Win32_NetworkAdapter) string {
		if networkInterface.NetEnabled {
			return humanize.SI(float64(networkInterface.Speed), "bit/s")
		}
		return "? bit/s"
	}

	for thisInterface := range dst {
		if isIgnoredExtInterfaceWindows(dst[thisInterface].Name) {
			log.Debug(fmt.Sprintf("Ignoring external network interface '%s'", dst[thisInterface].Name))
		} else if dst[thisInterface].PhysicalAdapter {
			linkSpeed := formatLinkSpeed(&dst[thisInterface])
			var oneInterface constants.AvailableSelection
			oneInterface.ConfigValue = dst[thisInterface].Name
			oneInterface.Legend = fmt.Sprintf("%s (%s)", dst[thisInterface].Name, linkSpeed)

			result = append(result, oneInterface)
		}
	}

	return result
}
