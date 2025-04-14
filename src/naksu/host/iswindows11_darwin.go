package host

// IsWindows11 returns always false on Darwin
func IsWindows11() bool {
	return false
}
