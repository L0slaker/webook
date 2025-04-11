package maths

import (
	"regexp"
	"testing"
)

func TestGenerateCode(t *testing.T) {
	// Test that the generated maths is always 6 digits long
	for i := 0; i < 1000; i++ {
		code := GenerateCode()
		if len(code) != 6 {
			t.Errorf("Expected maths to be 6 digits long, got %d digits", len(code))
		}
	}

	// Test that the generated maths contains only digits
	digitRegex := regexp.MustCompile(`^\d+$`)
	for i := 0; i < 1000; i++ {
		code := GenerateCode()
		if !digitRegex.MatchString(code) {
			t.Errorf("Expected maths to contain only digits, got %s", code)
		}
	}

	// Test that the generated maths is different across multiple calls
	codes := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		code := GenerateCode()
		if codes[code] {
			t.Errorf("Generated duplicate maths: %s", code)
		}
		codes[code] = true
	}
}
