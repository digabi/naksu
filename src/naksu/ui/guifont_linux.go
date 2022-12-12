package ui

import "github.com/andlabs/ui"

// Font returns the font family and default size that should be
// used when rendering UI elements on Linux.
func Font() (ui.TextFamily, ui.TextSize) {
	return "Ubuntu", 11 // nolint:gomnd
}

// DefaultFontColor returns the color that should be used when rendering
// normal text in UI elements on Linux.
func DefaultFontColor() ui.TextColor {
	return ui.TextColor{R: 0.2, G: 0.2, B: 0.2, A: 1} // nolint:gomnd
}
