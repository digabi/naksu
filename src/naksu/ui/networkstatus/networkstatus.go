package networkstatus

import (
	"fmt"

	"github.com/andlabs/ui"

	"naksu/network"
	naksuUi "naksu/ui"
	"naksu/xlate"
)

var networkStatusString *ui.AttributedString
var networkStatusArea *ui.Area

type networkStatusAreaHandler struct {
}

func (networkStatusAreaHandler) Draw(a *ui.Area, p *ui.AreaDrawParams) {
	fontFamily, size := naksuUi.Font()
	tl := ui.DrawNewTextLayout(&ui.DrawTextLayoutParams{
		String:      networkStatusString,
		Width:       p.AreaWidth,
		DefaultFont: &ui.FontDescriptor{Size: size, Family: fontFamily, Weight: ui.TextWeightNormal},
		Align:       ui.DrawTextAlignLeft,
	})
	defer tl.Free()
	p.Context.Text(tl, 0, 0)
}

func (networkStatusAreaHandler) MouseEvent(a *ui.Area, me *ui.AreaMouseEvent) {
}

func (networkStatusAreaHandler) MouseCrossed(a *ui.Area, left bool) {
}

func (networkStatusAreaHandler) DragBroken(a *ui.Area) {
}

func (networkStatusAreaHandler) KeyEvent(a *ui.Area, ke *ui.AreaKeyEvent) (handled bool) {
	// reject all keys
	return false
}

func appendWithAttributes(attributedString *ui.AttributedString, text string, attrs ...ui.Attribute) {
	start := len(attributedString.String())
	end := start + len(text)
	attributedString.AppendUnattributed(text)
	for _, a := range attrs {
		attributedString.SetAttribute(a, start, end)
	}
}

func ensureUIComponentsInitialized() {
	if networkStatusArea == nil {
		networkStatusArea = ui.NewArea(networkStatusAreaHandler{})
	}
	if networkStatusString == nil {
		networkStatusString = ui.NewAttributedString("")
	}
}

// Update network status area
func Update() {
	if network.UsingWirelessInterface() {
		showNetworkStatus(xlate.Get("Wireless connection"), true)
	} else {
		linkSpeedMbit := network.CurrentLinkSpeed()

		switch {
		case linkSpeedMbit == 0:
			showNetworkStatus(xlate.Get("No network connection"), true)
		case linkSpeedMbit < 1000:
			statusText := fmt.Sprintf(xlate.Get("Network speed is too low (%d Mbit/s)"), linkSpeedMbit)
			showNetworkStatus(statusText, true)
		default:
			showNetworkStatus(xlate.Get("OK"), false)
		}
	}
}

func showNetworkStatus(text string, warning bool) {
	ensureUIComponentsInitialized()
	networkStatusString = ui.NewAttributedString("")
	normalTextColor := naksuUi.DefaultFontColor()
	appendWithAttributes(networkStatusString, xlate.Get("Network status: "), normalTextColor)
	if warning {
		appendWithAttributes(networkStatusString, text, ui.TextColor{R: 1, G: 0, B: 0, A: 1})
	} else {
		appendWithAttributes(networkStatusString, text, normalTextColor)
	}
	networkStatusArea.QueueRedrawAll()
}

// Area returns the Area UI component singleton that shows the network status
func Area() *ui.Area {
	ensureUIComponentsInitialized()
	return networkStatusArea
}
