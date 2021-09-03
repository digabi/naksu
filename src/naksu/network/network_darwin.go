package network

import (
	"net"

	"naksu/constants"
	"naksu/log"
)

// GetExtInterfaces returns map of network interfaces. It returns array of
// constants.AvailableSelection where the ConfigValue is system's internal value
// (e.g. "eno1") and Legend human-readable legend (e.g. "Wireless Network 802.11abc")
func GetExtInterfaces() []constants.AvailableSelection {
	result := constants.DefaultExtNicArray

	interfaces, err := net.Interfaces()
	if err == nil {
		for n := range interfaces {
			if isIgnoredExtInterfaceDarwin(interfaces[n].Name) {
				log.Debug("Ingnoring external network interface '%s'", interfaces[n].Name)
			} else {
				log.Debug("Adding external network interface '%s' to the list of available devices", interfaces[n].Name)
				var oneInterface constants.AvailableSelection
				oneInterface.ConfigValue = interfaces[n].Name
				// We want to change this to more user-friendly later
				oneInterface.Legend = interfaces[n].Name

				result = append(result, oneInterface)
			}
		}
	}

	return result
}

// UsingWirelessInterface returns true if the user has selected a wireless
// interface in the naksu UI. If the user has made no selection or the selected
// interface is not wireless, returns false.
//
// Dummy implementation for MacOS.
func UsingWirelessInterface() bool {
	return false
}

// CurrentLinkSpeed returns the link speed of the selected ext interface
// or, if no interface selection has been made in the naksu UI, the interface
// that currently has the lowest link speed. The unit is megabits per second.
//
// Dummy implementation for MacOS.
func CurrentLinkSpeed() uint64 {
	return 0
}
