package utils

import (
	"fmt"
	"math"
	"strings"
)

// FormatNumber formats a number with thousand separators and 2 decimal places
// Example: 1234567.89 -> "1,234,567.89"
func FormatNumber(value float64) string {
	// Handle negative numbers
	negative := value < 0
	if negative {
		value = -value
	}

	// Split into integer and decimal parts
	integerPart := int64(math.Floor(value))
	decimalPart := int64(math.Round((value - math.Floor(value)) * 100))

	// Handle rounding edge case where decimal becomes 100
	if decimalPart >= 100 {
		integerPart++
		decimalPart = 0
	}

	// Format integer part with commas
	integerStr := formatIntegerWithCommas(integerPart)

	// Format decimal part (always 2 digits)
	decimalStr := fmt.Sprintf("%02d", decimalPart)

	// Combine parts
	result := fmt.Sprintf("%s.%s", integerStr, decimalStr)

	if negative {
		result = "-" + result
	}

	return result
}

// formatIntegerWithCommas adds thousand separators to an integer
func formatIntegerWithCommas(n int64) string {
	if n == 0 {
		return "0"
	}

	// Convert to string
	str := fmt.Sprintf("%d", n)

	// Add commas from right to left
	var result strings.Builder
	length := len(str)

	for i, digit := range str {
		if i > 0 && (length-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(digit)
	}

	return result.String()
}

