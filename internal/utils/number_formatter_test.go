package utils

import (
	"testing"
)

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected string
	}{
		{
			name:     "simple number",
			input:    1234567.89,
			expected: "1,234,567.89",
		},
		{
			name:     "number with zero decimal",
			input:    1000,
			expected: "1,000.00",
		},
		{
			name:     "small number",
			input:    500.5,
			expected: "500.50",
		},
		{
			name:     "zero",
			input:    0,
			expected: "0.00",
		},
		{
			name:     "negative number",
			input:    -1234.56,
			expected: "-1,234.56",
		},
		{
			name:     "large number",
			input:    123456789.99,
			expected: "123,456,789.99",
		},
		{
			name:     "single digit",
			input:    5.25,
			expected: "5.25",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatNumber(tt.input)
			if result != tt.expected {
				t.Errorf("FormatNumber(%f) = %s; expected %s", tt.input, result, tt.expected)
			}
		})
	}
}

