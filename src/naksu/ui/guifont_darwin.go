package ui

import "github.com/andlabs/ui"

// Font returns the font family and default size that should be
// used when rendering UI elements on Linux.
func Font() (ui.TextFamily, ui.TextSize) {
	return "Helvetica", 11
}

// DefaultFontColor returns the color that should be used when rendering
// normal text in UI elements on Linux.
func DefaultFontColor() ui.TextColor {
	return ui.TextColor{R: 0.5, G: 0.5, B: 0.5, A: 1}
}
