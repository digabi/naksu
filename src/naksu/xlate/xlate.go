package xlate

import (
	"fmt"

	"github.com/leonelquinteros/gotext"
)

var currentPoIsSet bool
var currentPo gotext.Po

// SetLanguage sets active language
func SetLanguage(newLanguage string) {
	currentPo = *gotext.NewPo()

	switch newLanguage {
	case "fi":
		currentPo.Parse([]byte(getPoStrFi()))
		currentPoIsSet = true
	case "sv":
		currentPo.Parse([]byte(getPoStrSv()))
		currentPoIsSet = true
	default:
		currentPoIsSet = false
	}
}

// Get returns translated string for given key with sprintf-like syntax
func Get(str string, vars ...interface{}) string {
	if currentPoIsSet {
		return currentPo.Get(str, vars...)
	}

	return fmt.Sprintf(str, vars...)
}

// GetRaw returns translated string for given key without processing sprintf-like syntax
func GetRaw(str string) string {
	if currentPoIsSet {
		return currentPo.Get(str)
	}

	return str
}
