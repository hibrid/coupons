package common

import (
	"fmt"
	"runtime"

	"github.com/shopspring/decimal"
)

// min returns the minimum of two float64 values.
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func decimalToString(d decimal.Decimal) string {
	return d.Round(2).String()
}

func printLnDecimalToString(d decimal.Decimal, varName string) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		fmt.Println("Error retrieving caller info")
		return
	}
	fmt.Printf("Variable '%s' at %s:%d: %s\n", varName, file, line, decimalToString(d))
}
