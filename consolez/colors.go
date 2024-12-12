package consolez

import (
	"github.com/fatih/color"
)

var (
	colorDefault            = color.New(color.Reset)
	colorHighlight          = color.New(color.Bold)
	colorSecondaryHighlight = color.New(color.Bold, color.Faint)
	colorSecondary          = color.New(color.Faint)
	colorInfo               = color.New(color.FgCyan)
	colorSuccess            = color.New(color.FgGreen)
	colorWarning            = color.New(color.FgYellow)
	colorError              = color.New(color.FgHiRed)
)

// GetColorDefault returns a color.
func GetColorDefault() *color.Color {
	return colorDefault
}

// GetColorHighlight returns a color.
func GetColorHighlight() *color.Color {
	return colorHighlight
}

// GetColorSecondaryHighlight returns a color.
func GetColorSecondaryHighlight() *color.Color {
	return colorSecondaryHighlight
}

// GetColorSecondary returns a color.
func GetColorSecondary() *color.Color {
	return colorSecondary
}

// GetColorInfo returns a color.
func GetColorInfo() *color.Color {
	return colorInfo
}

// GetColorSuccess returns a color.
func GetColorSuccess() *color.Color {
	return colorSuccess
}

// GetColorWarning returns a color.
func GetColorWarning() *color.Color {
	return colorWarning
}

// GetColorError returns a color.
func GetColorError() *color.Color {
	return colorError
}
