package host

// IsWindows11 returns always false on linux
func IsWindows11() bool {
	return false
}
