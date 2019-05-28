package network

import (
  "fmt"
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
			if IgnoreExtInterface(interfaces[n].Name) {
				log.Debug(fmt.Sprintf("Ingnoring external network interface '%s'", interfaces[n].Name))
			} else {
				log.Debug(fmt.Sprintf("Adding external network interface '%s' to the list of available devices", interfaces[n].Name))
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
