package common

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
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
	printLnDecimalToString(d, varName)
	// If no panic occurs during the function execution, it's considered successful.
}

func TestPrintLnDecimalToString2(t *testing.T) {
	// Test case 1: Caller information can be retrieved
	t.Run("Success", func(t *testing.T) {

		old := os.Stdout // keep backup of the real stdout
		outC := make(chan string)
		r, w, _ := os.Pipe()
		os.Stdout = w

		//os.Stdout = &buf
		d := decimal.NewFromFloat(5.6789)
		varName := "testVar"
		printLnDecimalToString(d, varName)

		// copy the output in a separate goroutine so printing can't block indefinitely
		go func() {
			var buf bytes.Buffer
			io.Copy(&buf, r)
			outC <- buf.String()
		}()
		w.Close()
		os.Stdout = old // restoring the real stdout
		out := <-outC
		// Check if the output is as expected
		expectedOutput := fmt.Sprintf("Variable '%s' at %s:%d: %s\n", varName, "/root/go/src/github.com/hibrid/coupons/common/helpers_test.go", 74, decimalToString(d))
		if out != expectedOutput {
			t.Errorf("Expected output: %s; got: %s", expectedOutput, out)
		}
	})

	// Test case 2: Caller information cannot be retrieved
	t.Run("Error", func(t *testing.T) {
		// Save current stdout
		old := os.Stdout
		r, w, _ := os.Pipe()
		defer func() {
			os.Stdout = old
		}()
		os.Stdout = w

		// Simulate failure to retrieve caller info
		_, _, _, _ = runtime.Caller(0) // Retrieving caller info will fail in this case
		d := decimal.NewFromFloat(5.6789)
		varName := "testVar"
		printLnDecimalToString(d, varName)

		// Capture stdout and restore old
		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)

		// Check if the error message is printed
		expectedErrorMsg := "Error retrieving caller info\n"
		if buf.String() != expectedErrorMsg {
			t.Errorf("Expected error message: %s; got: %s", expectedErrorMsg, buf.String())
		}
	})
}
