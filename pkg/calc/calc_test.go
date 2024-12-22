package calc

import (
	"testing"
)

func TestCalc(t *testing.T) {
	tests := []struct {
		expression    string
		expected      float64
		expectErr     bool
		expectedError string
	}{

		{"3 + 5", 8, false, ""},
		{"10 - 2 * 3", 4, false, ""},
		{"(1 + 2) * 3", 9, false, ""},
		{"10 / 2 + 3 * (2 - 1)", 8, false, ""},
		{"(3 + 5) * (2 - 1)", 8, false, ""},
		{"(1 + (2 + 3))", 6, false, ""},
		{"2 * (3 + 5)", 16, false, ""},
		{"(2 + 3) * (4 - 1)", 15, false, ""},
		{"3 * (2 + 5) - 4", 17, false, ""},

		{"(5 + 3))", 0, true, "invalid expression"},
		{"(5 + 3", 0, true, "invalid expression"},
		{"3 / 0", 0, true, "division by zero"},
		{"", 0, true, "empty expression or invalid request"},
		{"2 * (3 + 5", 0, true, "invalid expression"},
		{"2 * (3 + 5))", 0, true, "invalid expression"},
		{"(1 + 2) * (3 + 4))", 0, true, "invalid expression"},
		{"(1 + 2) * (3 + 4", 0, true, "invalid expression"},
		{"2 + * 3", 0, true, "invalid expression"},
		{"2 + (3 * (4 - 2)", 0, true, "invalid expression"},
	}

	for _, test := range tests {
		result, err := Calc(test.expression)
		if (err != nil) != test.expectErr {
			t.Errorf("Calc(%q) = unexpected error: %v, expected error: %v", test.expression, err, test.expectErr)
		}
		if test.expectErr {
			if err != nil && err.Error() != test.expectedError {
				t.Errorf("Calc(%q) = unexpected error message: %v; want %v", test.expression, err.Error(), test.expectedError)
			}
		} else {
			if result != test.expected {
				t.Errorf("Calc(%q) = %v; want %v", test.expression, result, test.expected)
			}
		}
	}
}
