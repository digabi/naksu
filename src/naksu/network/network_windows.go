package network

import (
	"fmt"
	"math"
	"regexp"

	"github.com/StackExchange/wmi"
	humanize "github.com/dustin/go-humanize"

	"naksu/config"
	"naksu/constants"
	"naksu/log"
)

// Win32_NetworkAdapter must be named with an underscore.
// Otherwise it is not recognised, which results in an "Invalid class" exception.
type Win32_NetworkAdapter struct { //nolint
	Name            *string
	Speed           uint64
	PhysicalAdapter bool
	NetEnabled      bool
}

func queryInterfaces(filter string) []Win32_NetworkAdapter {
	var dst []Win32_NetworkAdapter
	query := wmi.CreateQuery(&dst, filter)
	err := wmi.Query(query, &dst)
	if err != nil {
		log.Debug("Could not query network adapters from WMI")
		log.Debug(fmt.Sprint(err))
		return []Win32_NetworkAdapter{}
	}
	return dst
}

// GetExtInterfaces returns a slice of network interfaces, represented as
// values of type constants.AvailableSelection where ConfigValue is the
// system's internal name for the interface  (e.g. "Intel(R) Ethernet Connection
// (4) I219-LM") and Legend is a human-readable legend.
//
// In Windows these two values are the same.
func GetExtInterfaces() []constants.AvailableSelection {
	result := constants.DefaultExtNicArray
	interfaces := queryInterfaces("WHERE PhysicalAdapter=TRUE")

	formatLinkSpeed := func(networkInterface *Win32_NetworkAdapter) string {
		if networkInterface.NetEnabled {
			return humanize.SI(float64(networkInterface.Speed), "bit/s")
		}
		return "? bit/s"
	}

	for thisInterface := range interfaces {
		if isIgnoredExtInterfaceWindows(*interfaces[thisInterface].Name) {
			log.Debug(fmt.Sprintf("Ignoring external network interface '%s'", *interfaces[thisInterface].Name))
		} else if interfaces[thisInterface].PhysicalAdapter {
			linkSpeed := formatLinkSpeed(&interfaces[thisInterface])
			physicalInterface := constants.AvailableSelection{
				ConfigValue: *interfaces[thisInterface].Name,
				Legend:      fmt.Sprintf("%s (%s)", *interfaces[thisInterface].Name, linkSpeed),
			}

			result = append(result, physicalInterface)
		}
	}

	return result
}

// UsingWirelessInterface returns true if the user has selected a wireless
// interface in the naksu UI. If the user has made no selection or the selected
// interface is not wireless, returns false.
func UsingWirelessInterface() bool {
	if selectedInterfaceName := config.GetExtNic(); selectedInterfaceName != "" {
		nameContainsWireless, err := regexp.MatchString("Wireless", selectedInterfaceName)
		if err != nil {
			log.Debug("Could not check if the current interface is wireless")
			return false
		}
		return nameContainsWireless
	}
	return false
}

func selectedInterfaceOrAll() []Win32_NetworkAdapter {
	if selectedInterface := config.GetExtNic(); selectedInterface != "" {
		// #nosec (SQL query formatting warning)
		interfaces := queryInterfaces(fmt.Sprintf("WHERE Name='%s'AND  PhysicalAdapter=TRUE AND NetEnabled=TRUE", selectedInterface))
		if len(interfaces) != 1 {
			log.Debug(fmt.Sprintf("Found %d (not 1!) adapters with name %s", len(interfaces), selectedInterface))
		}
		return interfaces
	}
	return queryInterfaces("WHERE PhysicalAdapter=TRUE AND NetEnabled=TRUE")
}

// CurrentLinkSpeed returns the link speed of the selected ext interface
// or, if no interface selection has been made in the naksu UI, the interface
// that currently has the lowest link speed. The unit is megabits per second.
func CurrentLinkSpeed() uint64 {
	interfaces := selectedInterfaceOrAll()

	var minLinkSpeed uint64 = math.MaxUint64
	for _, thisInterface := range interfaces {
		if !isIgnoredExtInterfaceWindows(*thisInterface.Name) && thisInterface.Speed < minLinkSpeed {
			minLinkSpeed = thisInterface.Speed
		}
	}

	if minLinkSpeed == math.MaxUint64 {
		return 0
	}
	return bpsToMbps(minLinkSpeed)
}
