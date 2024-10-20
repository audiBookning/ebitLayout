package widgets

import (
	"strings"
)

// Helper function to check if a rune is a word separator
func isWordSeparator(r rune) bool {
	separators := " \n\t.,;:!?'\"()-"
	return strings.ContainsRune(separators, r)
}

// Helper function to clamp a value within a range
func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func clampFloat(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Helper function to get the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
