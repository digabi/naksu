package mebroutines

// ExistsStdin returns always true on windows
func ExistsStdin() bool {
	LogDebug("Windows is always expected to have STDIN")
	return true
}
