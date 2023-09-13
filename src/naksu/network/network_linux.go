package network

import (
	"errors"
	"fmt"
	"math"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/google/gousb"
	"github.com/google/gousb/usbid"
	"github.com/jaypipes/pcidb"

	"naksu/config"
	"naksu/constants"
	"naksu/log"
	"naksu/mebroutines"
)

type nicType = int

const (
	nicRegexWireless = "^w"
	nicRegexEthernet = "^(en)|(em)|(eth)"
	nicTypeUnknown   = iota
	nicTypePCI
	nicTypeUSB
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
		log.Error("Network interface speed file '%s' does not exist", speedPath)

		return 0
	}

	/* #nosec */
	carrierFileContent, err := os.ReadFile(carrierPath)
	if err != nil {
		if strings.HasSuffix(err.Error(), "invalid argument") {
			// When the network interface is powered down, trying to read /sys/class/net/<device>/carrier
			// results in an "invalid argument" error. This is the kernel working as intended.

			return 0
		}

		log.Error("Unexpected error while trying to read link state from %s: %s", carrierPath, err.Error())
	} else if strings.TrimSpace(string(carrierFileContent)) == "0" {
		return 0 // The link is down e.g. because the device is disconnected from the network.
	}

	/* #nosec */
	speedFileContent, err := os.ReadFile(speedPath)
	if err != nil {
		log.Error("Could not read network interface speed from '%s': %v", speedPath, err)

		return 0
	}
	speedFileContentTrimmed := strings.TrimSpace(string(speedFileContent))
	speedInt, errConvert := strconv.ParseUint(speedFileContentTrimmed, 10, 64)
	if errConvert != nil {
		log.Error("Could not convert speed string '%s' to integer: %v", speedFileContentTrimmed, errConvert)

		return 0
	}

	return speedInt * 1000000 // nolint:gomnd
}

func getExtInterfaceType(extInterface string) nicType {
	modaliasPath := fmt.Sprintf("/sys/class/net/%s/device/modalias", extInterface)

	/* #nosec */
	modaliasContent, err := os.ReadFile(modaliasPath)
	if err != nil {
		log.Warning("Could not detect type of external network interface %s: %v", extInterface, err)

		return 0
	}

	// Sample modalias string:
	// pci:v00008086d000024FBsv00008086sd00002110bc02sc80i00
	modaliasStr := string(modaliasContent)
	const modaliasTypePrefixLength = 4
	if len(modaliasStr) >= modaliasTypePrefixLength {
		switch modaliasStr[:modaliasTypePrefixLength] {
		case "usb:":
			return nicTypeUSB
		case "pci:":
			return nicTypePCI
		}
	}

	log.Warning("Could not detect type of external network interface %s (%s): %s", extInterface, modaliasPath, modaliasStr)

	return nicTypeUnknown
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

func getPCIDeviceLegend(vendor string, device string) (string, error) {
	vendor = strings.ToLower(vendor)
	device = strings.ToLower(device)

	log.Debug("Getting device name from PCI database: %s:%s", vendor, device)

	searchKey := fmt.Sprintf("%s%s", vendor, device)

	if linuxPCIDatabase == nil {
		var err error
		// Do not try to load PCI database from the network
		linuxPCIDatabase, err = pcidb.New()
		if err != nil {
			log.Error("Could not initialise PCI database: %v", err)

			return "", err
		}
	}

	for key, devProduct := range linuxPCIDatabase.Products {
		if key == searchKey {
			log.Debug("Found device from PCI database: %s", devProduct.Name)

			return devProduct.Name, nil
		}
	}

	log.Debug("Did not find any matches from PCI database for key %s", searchKey)

	return "", errors.New("no matching legend found")
}

func getUSBDeviceLegend(vendor string, product string) (string, error) {
	var desc *gousb.DeviceDesc

	log.Debug("Getting device name from USB database: %s:%s", vendor, product)

	vendorInt, vendorErr := strconv.ParseUint(vendor, 16, 32)
	if vendorErr != nil {
		return "", fmt.Errorf("vendor id cannot be converted to hex: %w", vendorErr)
	}

	productInt, productErr := strconv.ParseUint(product, 16, 32)
	if productErr != nil {
		return "", fmt.Errorf("product id cannot be converted to hex: %w", productErr)
	}

	desc = &gousb.DeviceDesc{Vendor: gousb.ID(vendorInt), Product: gousb.ID(productInt)} // nolint:exhaustruct

	legend := usbid.Describe(desc)
	log.Debug("USB database gives following legend to device %s:%s: %s", vendor, product, legend)

	return legend, nil
}

func getPCIExtInterfaceLegend(extInterface string) (string, error) {
	vendorPath := fmt.Sprintf("/sys/class/net/%s/device/vendor", extInterface)
	devicePath := fmt.Sprintf("/sys/class/net/%s/device/device", extInterface)

	const pciDeviceCodeMinimumLength = 6

	prepareID := func(idByte []byte) (string, error) {
		idStr := string(idByte)
		if len(idStr) < pciDeviceCodeMinimumLength {
			return "", fmt.Errorf("malformatted pci device code: %s", idStr)
		}

		return idStr[2:6], nil
	}

	/* #nosec */
	vendor, errVendor := os.ReadFile(vendorPath)
	if errVendor != nil {
		log.Error("Trying to get vendor ID for %s but could not open %s for reading: %v", extInterface, vendorPath, errVendor)

		return "", fmt.Errorf("failed to get vendor id for external pci network interface %s", extInterface)
	}
	vendorID, vendorErr := prepareID(vendor)

	if vendorErr != nil {
		log.Error("Failed to get PCI vendor ID for %s: %v", extInterface, vendorErr)

		return "", vendorErr
	}

	/* #nosec */
	device, errDevice := os.ReadFile(devicePath)
	if errDevice != nil {
		log.Error("Trying to get device ID for %s but could not open %s for reading: %v", extInterface, devicePath, errDevice)

		return "", fmt.Errorf("failed to get device id for external pci network interface %s", extInterface)
	}
	deviceID, deviceErr := prepareID(device)

	if deviceErr != nil {
		log.Error("Failed to get PCI device ID for %s: %v", extInterface, deviceErr)

		return "", deviceErr
	}

	pciLegend, errSearch := getPCIDeviceLegend(vendorID, deviceID)
	if errSearch != nil {
		log.Error("Failed to get device legend for a PCI device %s:%s: %v", vendorID, deviceID, errSearch)
	}

	return pciLegend, errSearch
}

func getUSBExtInterfaceLegend(extInterface string) (string, error) {
	modaliasPath := fmt.Sprintf("/sys/class/net/%s/device/modalias", extInterface)

	/* #nosec */
	modaliasContent, err := os.ReadFile(modaliasPath)
	if err != nil {
		log.Error("Could not get vendor/product codes for external network interface %s: %v", extInterface, err)

		return "", fmt.Errorf("could not get vendor/product codes for external network interface %s: %w", extInterface, err)
	}

	modaliasStr := string(modaliasContent)

	// See http://people.skolelinux.org/pere/blog/Modalias_strings___a_practical_way_to_map__stuff__to_hardware.html
	const minimumModaliasRowLength = 14
	if len(modaliasStr) > minimumModaliasRowLength {
		vendorID := strings.ToLower(modaliasStr[5:9])
		deviceID := strings.ToLower(modaliasStr[10:14])

		usbLegend, errSearch := getUSBDeviceLegend(vendorID, deviceID)
		if errSearch != nil {
			log.Error("Failed to get device legend for a USB device %s:%s: %v", vendorID, deviceID, errSearch)
		}

		return usbLegend, errSearch
	}

	return "", fmt.Errorf("malformatted modalias string for external network interface %s (%s): %s", extInterface, modaliasPath, modaliasStr)
}

func getExtInterfaceLegend(extInterface string) string {
	// Defaults to generic device type derived from network name, see getExtInterfaceDefaultLegend()
	var legend string
	var err error

	extInterfaceType := getExtInterfaceType(extInterface)

	switch extInterfaceType {
	case nicTypePCI:
		legend, err = getPCIExtInterfaceLegend(extInterface)
	case nicTypeUSB:
		legend, err = getUSBExtInterfaceLegend(extInterface)
	default:
		err = fmt.Errorf("unknown network interface type: %d", extInterfaceType)
	}

	if err != nil {
		log.Error("Failed to get legend for external network interface %s: %v", extInterface, err)
		legend = getExtInterfaceDefaultLegend(extInterface)
	}

	return legend
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
			log.Debug("Could not check if the current interface is wireless")

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
	if err != nil {
		log.Debug("Query to external interfaces got an error: %v", err.Error())

		return result
	}

	for interfaceIndex := range interfaces {
		if isIgnoredExtInterfaceLinux(interfaces[interfaceIndex].Name) {
			log.Debug("Ignoring external network interface '%s'", interfaces[interfaceIndex].Name)
		} else {
			log.Debug("Adding external network interface '%s' to the list of available devices", interfaces[interfaceIndex].Name)
			var oneInterface constants.AvailableSelection
			oneInterface.ConfigValue = interfaces[interfaceIndex].Name

			legend := getExtInterfaceLegend(interfaces[interfaceIndex].Name)

			speed := getExtInterfaceSpeed(interfaces[interfaceIndex].Name)
			if speed > 0 {
				oneInterface.Legend = fmt.Sprintf("%s %s (%s)", legend, interfaces[interfaceIndex].Name, humanize.SI(float64(speed), "bit/s"))
			} else {
				oneInterface.Legend = fmt.Sprintf("%s %s", legend, interfaces[interfaceIndex].Name)
			}

			result = append(result, oneInterface)
		}
	}

	return result
}
