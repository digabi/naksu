package host

// IsHyperV returns always false on Linux
func IsHyperV() bool {
	return false
}
