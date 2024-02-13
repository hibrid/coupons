package generator

import (
	"strings"
	"testing"
)

func TestDefault(t *testing.T) {
	defaults := Default()
	if defaults.MinimumLength != 6 && defaults.Count != 1 {
		t.Error("wrong defaults")
	}
}

func TestCountUniqueCombinations(t *testing.T) {
	// Initialize common variables for tests
	numbers := "0123456789"
	alphabetic := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphanumeric := numbers + alphabetic
	patternChar := "#"

	// Define test cases
	tests := []struct {
		name     string
		pattern  string
		expected int
	}{
		{"No Pattern Char", "ABC-DEF", 1},
		{"Single Pattern Char", "A#C", 36},
		{"Two Pattern Chars", "A#B#", 1296},       // 36^2
		{"Pattern With Divider", "#-#", 1296},     // 36^2
		{"Complex Pattern", "A#B#-#C#D", 1679616}, // 36^4
	}

	// Execute each test
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gen, err := NewWithOptions(
				SetPattern(tc.pattern),
			)
			if err != nil {
				t.Errorf("Error creating generator: %s", err)
			}
			result := gen.CountUniqueCombinations(tc.pattern, patternChar, alphanumeric)
			if result != tc.expected {
				t.Errorf("For pattern '%s', expected %d combinations, got %d", tc.pattern, tc.expected, result)
			}
		})
	}
}

func TestNew(t *testing.T) {
	g, _ := NewWithOptions(
		SetMinimumLength(5),
		SetGenerateCount(5),
		SetPattern("#####-#####"),
	)
	if g.MinimumLength != 5 && g.Pattern != "#####-#####" {
		t.Error("wonky constructor")
	}
}

func TestCheckDigit(t *testing.T) {
	var runs = []struct {
		code string
		part int
	}{
		{"ASD7", 1},
		{"ZPMR", 2}, // same as 2PMR
		{"2PMR", 2},
		{"4GGL", 3},
		{"GUX0", 4},
		{"0LJ2", 1},
		{"OLJ2", 1}, // same as 0LJW
		{"W8MY", 2},
		{"XXMT", 3},
		{"7TVQA", 1},
		{"PHN56", 2},
		{"5BD70", 3},

		{"23ERBS", 1},
		{"1R1N4Y", 2},
		{"F11MHM", 3},
		{"3U00FO", 4},
	}
	for _, run := range runs {
		check := checkCharacter(run.code[:len(run.code)-1], run.part)
		if !strings.HasSuffix(run.code, strings.ToUpper(check)) {
			t.Errorf("check digit failed for %s got %s", run.code, check)
		}
	}
}

func TestSecureCheckDigit(t *testing.T) {
	var runs = []struct {
		code string
		part int
	}{
		{"ASDJ", 1},
		{"ZPMS", 2}, // same as 2PMR
		{"2PMS", 2},
		{"4GG7", 3},
		{"GUXK", 4},
		{"0LJ1", 1},
		{"OLJ1", 1}, // same as 0LJW
		{"W8M5", 2},
		{"XXMU", 3},
		{"7TVQE", 1},
		{"PHN5B", 2},
		{"5BD7Y", 3},

		{"23ERBT", 1},
		{"1R1N4R", 2},
		{"F11MH1", 3},
		{"3U00F9", 4},
	}
	for _, run := range runs {
		check := secureCheckCharacter(run.code[:len(run.code)-1], run.part, 0)
		if !strings.HasSuffix(run.code, strings.ToUpper(check)) {
			t.Errorf("check digit failed for %s got %s", run.code, check)
		}
	}
}

func TestCreateCode(t *testing.T) {
	g, _ := NewWithOptions(
		SetMinimumLength(4),
		SetGenerateCount(100000),
		SetPattern("######-######-######-######"),
	)
	code, err := g.Run()
	if err != nil {
		t.Errorf("code should be valid got %s", err)
	}
	if len(code[0]) != 27 {
		t.Errorf("code should be 27 characters long got %d", len(code[0]))
	}
	for i := 0; i < 1000; i++ {

		validatedCode, err := g.Validate(code[i])
		if err != nil {
			t.Errorf("code should be valid got %s", err)
		}
		if code[i] != validatedCode {
			t.Errorf("code should be valid got %s", validatedCode)
		}
	}
}

func TestCreateCodeWithPrefix(t *testing.T) {
	g, _ := NewWithOptions(
		SetMinimumLength(4),
		SetGenerateCount(100000),
		SetPattern("######-######-######-######"),
		SetPrefix("TEST"),
	)
	code, err := g.Run()
	if err != nil {
		t.Errorf("code should be valid got %s", err)
	}
	if len(code[0]) != 32 {
		t.Errorf("code should be 32 characters long got %d", len(code[0]))
	}
	for i := 0; i < 1000; i++ {

		validatedCode, err := g.Validate(code[i])
		if err != nil {
			t.Errorf("code should be valid got %s", err)
		}
		if code[i] != validatedCode {
			t.Errorf("code should be valid got %s", validatedCode)
		}
	}
}

func TestCreateCodeWithSuffix(t *testing.T) {
	g, _ := NewWithOptions(
		SetMinimumLength(4),
		SetGenerateCount(100000),
		SetPattern("######-######-######-######"),
		SetSuffix("TEST"),
	)
	code, err := g.Run()
	if err != nil {
		t.Errorf("code should be valid got %s", err)
	}
	if len(code[0]) != 32 {
		t.Errorf("code should be 32 characters long got %d", len(code[0]))
	}
	for i := 0; i < 1000; i++ {

		validatedCode, err := g.Validate(code[i])
		if err != nil {
			t.Errorf("code should be valid got %s", err)
		}
		if code[i] != validatedCode {
			t.Errorf("code should be valid got %s", validatedCode)
		}
	}
}

func TestCreateCodeWithPrefixAndSuffix(t *testing.T) {
	g, _ := NewWithOptions(
		SetMinimumLength(4),
		SetGenerateCount(100000),
		SetPattern("######-######-######-######"),
		SetSuffix("TEST"),
		SetPrefix("TEST"),
	)
	code, err := g.Run()
	if err != nil {
		t.Errorf("code should be valid got %s", err)
	}
	if len(code[0]) != 37 {
		t.Errorf("code should be 37 characters long got %d", len(code[0]))
	}
	for i := 0; i < 1000; i++ {

		validatedCode, err := g.Validate(code[i])
		if err != nil {
			t.Errorf("code should be valid got %s", err)
		}
		if code[i] != validatedCode {
			t.Errorf("code should be valid got %s", validatedCode)
		}
	}
}

func TestValidCodes(t *testing.T) {
	var runs = []struct {
		options []Option
		code    string
	}{
		{[]Option{}, "55GR-DHME-50NT"},
		{[]Option{SetGenerateCount(1), SetPattern("####-####-####-####")}, "U5HD-HKD8-8RNL-1EXI"},
		{[]Option{SetGenerateCount(1), SetPattern("######-######-######")}, "WYLKQ9-U35V4O-9N84DY"},
	}
	for _, run := range runs {
		generator, err := NewWithOptions(run.options...) //run.g.Validate(run.code)
		if err != nil {
			t.Errorf("code %s should be valid got %s", run.code, err)
		}
		code, err := generator.Validate(run.code)
		if err != nil {
			t.Errorf("code %s should be valid got %s %s", run.code, code, err)
		}
	}
}

func TestCodesWithPrefix(t *testing.T) {
	var runs = []struct {
		options     []Option
		code        string
		expectError bool
	}{
		{[]Option{SetPrefix("test")}, "test-55GR-DHME-50NT", true},
		{[]Option{SetGenerateCount(1), SetPattern("####-####-####-####")}, "test-U5HL-HKDI-8RNQ-1EXQ", true},
		{[]Option{SetGenerateCount(1), SetPattern("######-######-######")}, "WYLKQ9-U35V4O-9N84DY", false},
		{[]Option{SetGenerateCount(1), SetPattern("######-######-######")}, "WYLKQ9-U35V4O-9N84DY-test", true},
		{[]Option{SetGenerateCount(1), SetPattern("######-######-######")}, "test-WYLKQ9-U35V4O-9N84DY-test", true},

		{[]Option{SetPrefix("test"), SetSuffix("test2"), SetGenerateCount(1), SetPattern("######-######-######")}, "test-WYLKQ9-U35V4O-9N84DY-test2", false},
	}
	for _, run := range runs {
		generator, err := NewWithOptions(run.options...) //run.g.Validate(run.code)
		if err != nil {
			t.Errorf("code %s should be valid got %s", run.code, err)
		}
		code, err := generator.Validate(run.code)
		if err != nil && !run.expectError {
			t.Errorf("code %s should be valid got %s %s", run.code, code, err)
		}
	}
}
