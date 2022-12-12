package ui

import "github.com/andlabs/ui"

// Font returns the font family and default size that should be
// used when rendering UI elements on Windows.
func Font() (ui.TextFamily, ui.TextSize) {
	return "Segoe UI", 10 // nolint:gomnd
}

// DefaultFontColor returns the color that should be used when rendering
// normal text in UI elements on Windows.
func DefaultFontColor() ui.TextColor {
	return ui.TextColor{R: 0, G: 0, B: 0, A: 1}
}
