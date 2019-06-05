package network

// IsIgnoredExtInterface returns true if given system-level network device should
// be ignored (i.e. not to be shown to the user).
func IsIgnoredExtInterface(interfaceName string) bool {
	return isIgnoredExtInterfaceWindows(interfaceName)
}
