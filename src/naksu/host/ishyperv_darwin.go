package host

// IsHyperV returns always false on Darwin
func IsHyperV() bool {
	return false
}
