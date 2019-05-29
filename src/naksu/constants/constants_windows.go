package constants

// ExtNicsToIgnore is a list of devices to ignore. Each entry is a regular expression
// and the entry is dropped if any of these regexs match. See network.IgnoreExtInterface()
// for more.
var ExtNicsToIgnore = []string{}

// ExtNicNixLegendRules is a map between regular expressions matching *nix device names
// and user-friendly legends. This is necessary while we don't call lshw or similair
// to description of the network devices
var ExtNicNixLegendRules = []struct {
	RegExp string
	Legend string
}{}
