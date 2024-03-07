package common

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
)

func TestMin(t *testing.T) {
	tests := []struct {
		a, b, expected float64
	}{
		{5.6, 7.8, 5.6},
		{10.2, 3.4, 3.4},
		{-3.5, -2.9, -3.5},
		{0, 0, 0},
	}

	for _, test := range tests {
		result := min(test.a, test.b)
		if result != test.expected {
			t.Errorf("min(%f, %f) = %f; want %f", test.a, test.b, result, test.expected)
		}
	}
}

func TestDecimalToString(t *testing.T) {
	tests := []struct {
		d        decimal.Decimal
		expected string
	}{
		{decimal.NewFromFloat(5.6789), "5.68"},
		{decimal.NewFromFloat(10.0), "10"},
		{decimal.NewFromFloat(-3.14159), "-3.14"},
		{decimal.NewFromFloat(0.0), "0"},
	}

	for _, test := range tests {
		result := decimalToString(test.d)
		if result != test.expected {
			t.Errorf("decimalToString(%s) = %s; want %s", test.d.String(), result, test.expected)
		}
	}
}

func TestPrintLnDecimalToString(t *testing.T) {
	// As this function involves printing to stdout, testing becomes tricky.
	// You could refactor the code to make it more testable by separating concerns.
	// For now, we'll just ensure it doesn't panic.
	// Mocking stdout is also an option but would be more complex.
	d := decimal.NewFromFloat(5.6789)
	varName := "testVar"
	printLnDecimalToString(d, varName, RuntimeCallerInfoProvider{})
	// If no panic occurs during the function execution, it's considered successful.
}

func TestPrintLnDecimalToString3(t *testing.T) {
	// Backup stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	varName := "testDecimal"
	testDecimal := decimal.NewFromFloat(123.456)

	// Expected patterns in the output
	expectedDecimalPattern := decimalToString(testDecimal) // Use the same formatting function to match the output

	// Execute the function
	printLnDecimalToString(testDecimal, varName, RuntimeCallerInfoProvider{})

	// Close pipe and restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read the output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Compile regex to match the expected output format loosely
	expectedOutputRegex := regexp.MustCompile(fmt.Sprintf("Variable '%s' at .*:\\d+: %s\n", varName, expectedDecimalPattern))

	// Test assertions
	if !expectedOutputRegex.MatchString(output) {
		t.Errorf("Output does not match the expected format.\nOutput: %s\nExpected to match pattern: %s", output, expectedOutputRegex.String())
	}
}

type MockFailCallerInfoProvider struct{}

func (MockFailCallerInfoProvider) Caller(skip int) (uintptr, string, int, bool) {
	return 0, "", 0, false // Simulate failure
}

func TestPrintLnDecimalToStringFailCase(t *testing.T) {
	// Backup stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Use the mock provider that simulates failure
	provider := MockFailCallerInfoProvider{}

	// Execute the function with the mock provider
	printLnDecimalToString(decimal.NewFromFloat(123.456), "testDecimal", provider)

	// Close pipe and restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read the output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	expectedOutput := "Error retrieving caller info"

	// Test assertion for the failure case
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Output does not contain the expected error message.\nOutput: %s\nExpected to contain: %s", output, expectedOutput)
	}
}
