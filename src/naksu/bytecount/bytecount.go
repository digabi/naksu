// Convert byte size to bytes, kilobytes, megabytes, gigabytes, ...
// in either SI (decimal) or IEC (binary) format.
// Source: https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/

package bytecount

import (
	"fmt"
)

// BytesToHumanSI converts given bytes to human-readable format using SI megs (1000 instead of 1024)
func BytesToHumanSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

// BytesToHumanIEC converts given bytes to human-readable format using IEC megs (1024 instead of 1000)
func BytesToHumanIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}
