package host

import (
	"fmt"

	"naksu/mebroutines"
)

// getPowerplan returns a string describing current power plan
func getPowerplan() string {
	powercfgListOutput, _ := mebroutines.RunAndGetOutput([]string{"powercfg", "/list"}, false)
	powercfgQueryOutput, _ := mebroutines.RunAndGetOutput([]string{"powercfg", "/query"}, false)

	return fmt.Sprintf("%s\n\n%s", powercfgListOutput, powercfgQueryOutput)
}
