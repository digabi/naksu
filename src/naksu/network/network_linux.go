package network

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"regexp"
	"strconv"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/jaypipes/pcidb"

	"naksu/config"
	"naksu/constants"
	"naksu/log"
	"naksu/mebroutines"
)

const (
	nicRegexWireless = "^w"
	nicRegexEthernet = "^(en)|(em)|(eth)"
)

// extNicNixDefaultLegendRules is a map between regular expressions matching *nix device names
// and user-friendly legends. This map will be used if we cannot detect the real product name
// from /sys/class/net/*/device/{vendor,device}
var extNicNixDefaultLegendRules = []struct {
	RegExp string
	Legend string
}{
	{nicRegexWireless, "Wireless"},
	{nicRegexEthernet, "Ethernet"},
}

// linuxPCIDatabase is a map reflecting the PCI device database. It will be initialised
// by getPCIDeviceLegend()
var linuxPCIDatabase *pcidb.PCIDB

func getExtInterfaceSpeed(extInterface string) uint64 {
	carrierPath := fmt.Sprintf("/sys/class/net/%s/carrier", extInterface)
	speedPath := fmt.Sprintf("/sys/class/net/%s/speed", extInterface)

	if !mebroutines.ExistsFile(speedPath) {
		log.Debug(fmt.Sprintf("Network interface speed file '%s' does not exist", speedPath))
		return 0
	}

	/* #nosec */
	carrierFileContent, err := ioutil.ReadFile(carrierPath)
	if err != nil {
		log.Debug(fmt.Sprintf("Could not read network interface carrier from '%s': %v", carrierPath, err))
		return 0
	}
	if strings.TrimSpace(string(carrierFileContent)) == "0" {
		return 0 // Link not up
	}

	/* #nosec */
	speedFileContent, err := ioutil.ReadFile(speedPath)
	if err != nil {
		log.Debug(fmt.Sprintf("Could not read network interface speed from '%s': %v", speedPath, err))
		return 0
	}
	speedFileContentTrimmed := strings.TrimSpace(string(speedFileContent))
	speedInt, errConvert := strconv.ParseUint(speedFileContentTrimmed, 10, 64)
	if errConvert != nil {
		log.Debug(fmt.Sprintf("Could not convert speed string '%s' to integer: %v", speedFileContentTrimmed, errConvert))
		return 0
	}

	return speedInt * 1000000
}

func getPCIDeviceLegend(vendor string, device string) (string, error) {
	searchKey := fmt.Sprintf("%s%s", strings.ToLower(vendor), strings.ToLower(device))
	log.Debug(fmt.Sprintf("Getting device name from PCI database: %s", searchKey))

	if linuxPCIDatabase == nil {
		var err error
		// Do not try to load PCI database from the network
		pcidb.WithDisableNetworkFetch()
		linuxPCIDatabase, err = pcidb.New()
		if err != nil {
			log.Debug(fmt.Sprintf("Could not initialise PCI database: %v", err))
			return "", err
		}
	}

	for key, devProduct := range linuxPCIDatabase.Products {
		if key == searchKey {
			log.Debug(fmt.Sprintf("Found device from PCI database: %s", devProduct.Name))
			return devProduct.Name, nil
		}
	}

	log.Debug(fmt.Sprintf("Did not find any matches from PCI database for key %s", searchKey))
	return "", errors.New("no matching legend found")
}

func getExtInterfaceDefaultLegend(extInterface string) string {
	for _, table := range extNicNixDefaultLegendRules {
		matched, err := regexp.MatchString(table.RegExp, extInterface)
		if err == nil && matched {
			return table.Legend
		}
	}

	return "Unknown device"
}

func getExtInterfaceLegend(extInterface string) string {
	vendorPath := fmt.Sprintf("/sys/class/net/%s/device/vendor", extInterface)
	devicePath := fmt.Sprintf("/sys/class/net/%s/device/device", extInterface)

	defaultLegend := getExtInterfaceDefaultLegend(extInterface)

	prepareID := func(idByte []byte) string {
		idStr := string(idByte)
		if len(idStr) < 6 {
			return ""
		}

		return idStr[2:6]
	}

	/* #nosec */
	vendor, errVendor := ioutil.ReadFile(vendorPath)
	if errVendor != nil {
		log.Debug(fmt.Sprintf("Trying to get vendor ID for %s but could not open %s for reading: %v", extInterface, vendorPath, errVendor))
		return defaultLegend
	}
	vendorID := prepareID(vendor)

	/* #nosec */
	device, errDevice := ioutil.ReadFile(devicePath)
	if errDevice != nil {
		log.Debug(fmt.Sprintf("Trying to get device ID for %s but could not open %s for reading: %v", extInterface, devicePath, errDevice))
		return defaultLegend
	}
	deviceID := prepareID(device)

	pciLegend, errSearch := getPCIDeviceLegend(vendorID, deviceID)
	if errSearch != nil {
		log.Debug(fmt.Sprintf("Failed to get device legend for a PCI device %s:%s: %v", vendorID, deviceID, errSearch))
		return defaultLegend
	}

	return pciLegend
}

func interfaceNames() []string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return []string{}
	}

	physicalInterfaces := []string{}
	for n := range interfaces {
		if !isIgnoredExtInterfaceLinux(interfaces[n].Name) {
			physicalInterfaces = append(physicalInterfaces, interfaces[n].Name)
		}
	}

	return physicalInterfaces
}

// UsingWirelessInterface returns true if the user has selected a wireless
// interface in the naksu UI. If the user has made no selection or the selected
// interface is not wireless, returns false.
func UsingWirelessInterface() bool {
	if selectedInterface := config.GetExtNic(); selectedInterface != "" {
		isWireless, err := regexp.MatchString(nicRegexWireless, selectedInterface)
		if err != nil {
			log.Debug(fmt.Sprintf("Could not check if the current interface is wireless"))
			return false
		}
		return isWireless
	}
	return false
}

// CurrentLinkSpeed returns the link speed of the selected ext interface
// or, if no interface selection has been made in the naksu UI, the interface
// that currently has the lowest link speed. The unit is megabits per second.
//
// Only works for wired network interfaces.
func CurrentLinkSpeed() uint64 {
	var interfaces []string
	if selectedInterface := config.GetExtNic(); selectedInterface != "" {
		interfaces = []string{selectedInterface}
	} else {
		interfaces = interfaceNames()
	}

	var minLinkSpeed uint64 = math.MaxUint64
	for _, interfaceName := range interfaces {
		linkSpeed := getExtInterfaceSpeed(interfaceName)
		if linkSpeed < minLinkSpeed && linkSpeed > 0 {
			minLinkSpeed = linkSpeed
		}
	}

	if minLinkSpeed == math.MaxUint64 {
		return 0
	}
	return bpsToMbps(minLinkSpeed)
}

// GetExtInterfaces returns a slice of network interfaces, represented as values
// of type constants.AvailableSelection where ConfigValue is the system's internal
// name for the interface (e.g. "eno1") and Legend is a human-readable legend
// (e.g. "Wireless Network 802.11abc")
//
// Currently the Linux implementation returns interface name and speed as a legend
func GetExtInterfaces() []constants.AvailableSelection {
	result := constants.DefaultExtNicArray

	interfaces, err := net.Interfaces()
	if err == nil {
		for n := range interfaces {
			if isIgnoredExtInterfaceLinux(interfaces[n].Name) {
				log.Debug(fmt.Sprintf("Ignoring external network interface '%s'", interfaces[n].Name))
			} else {
				log.Debug(fmt.Sprintf("Adding external network interface '%s' to the list of available devices", interfaces[n].Name))
				var oneInterface constants.AvailableSelection
				oneInterface.ConfigValue = interfaces[n].Name

				legend := getExtInterfaceLegend(interfaces[n].Name)

				speed := getExtInterfaceSpeed(interfaces[n].Name)
				if speed > 0 {
					oneInterface.Legend = fmt.Sprintf("%s %s (%s)", legend, interfaces[n].Name, humanize.SI(float64(speed), "bit/s"))
				} else {
					oneInterface.Legend = fmt.Sprintf("%s %s", legend, interfaces[n].Name)
				}

				result = append(result, oneInterface)
			}
		}
	}

	return result
}
