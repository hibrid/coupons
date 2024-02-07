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
	if g.MinimumLength != 5 && g.Pattern != "####-####" {
		t.Error("wonky constructor")
	}
}

func TestCheckDigit(t *testing.T) {
	var runs = []struct {
		code string
		part int
	}{
		{"ASD7", 1},
		{"ZPMR", 2},
		{"4GGB", 3},
		{"GUX5", 4},
		{"OLJW", 1},
		{"W8MY", 2},
		{"XXMT", 3},
		{"7TVQX", 1},
		{"PHN56", 2},
		{"5BD7X", 3},

		{"23ERBV", 1},
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
		t.Errorf("code should be 19 characters long got %d", len(code[0]))
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
		{[]Option{}, "55GP-DHMV-50N5"},
		{[]Option{SetGenerateCount(1), SetPattern("####-####-####-####")}, "U5HL-HKDI-8RNQ-1EXQ"},
		{[]Option{SetGenerateCount(1), SetPattern("######-######-######-######")}, "WYLKQJ-U35V4I-9N84DK"},
		/*
			{Default, "55g2-dhm0-50nn"},
			{Default, "SSGZ-DHMO-SONN"},
			{New(7, 12), "QBXA5CV4Q85E-HNYV4U3UD69M-B7XU1BHF3FYE-HXT9LD4Q0DAH-U6WMKC1WNF4N-5PCG5C4JF0GL-5DTUNJ40LRB5"},
			{New(1, 4), "1K7Q"},
			{New(2, 4), "1K7Q-CTFM"},
			{New(3, 4), "1K7Q-CTFM-LMTC"},
			{New(4, 4), "7YQH-1FU7-E1HX-0BG9"},
			{New(5, 4), "YENH-UPJK-PTE0-20U6-QYME"},
			{New(6, 4), "YENH-UPJK-PTE0-20U6-QYME-RBK1"},
		*/
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

/*
func TestInvalidCodes(t *testing.T) {
	var runs = []struct {
		g    *generator
		code string
	}{
		{Default, "55G2-DHM0-50NK"}, // wrong check
		{Default, "55G2-DHM-50NN"},  // not enough characters
		{New(3, 4), "1K7Q-CTFM"},    // too short
		{New(1, 4), "1K7C"},
		{New(2, 4), "1K7Q-CTFW"},
		{New(3, 4), "1K7Q-CTFM-LMT1"},
		{New(4, 4), "7YQH-1FU7-E1HX-0BGP"},
		{New(5, 4), "YENH-UPJK-PTE0-20U6-QYMT"},
		{New(6, 4), "YENH-UPJK-PTE0-20U6-QYME-RBK2"},
	}
	for _, run := range runs {
		validated, err := run.g.Validate(run.code)
		if err == nil {
			t.Errorf("code %s should be invalid got %s", run.code, validated)
		}
	}
}

func TestCodesNormalized(t *testing.T) {
	var runs = []struct {
		g    *generator
		code string
		exp  string
	}{
		{Default, "i9oD/V467/8Dsz", "190D-V467-8D52"},   // alternate separator
		{Default, " i9oD V467 8Dsz ", "190D-V467-8D52"}, // whitespace accepted
		{Default, " i9oD_V467_8Dsz ", "190D-V467-8D52"}, // underscores accepted
		{Default, "i9oDV4678Dsz", "190D-V467-8D52"},     // no separator required
	}
	for _, run := range runs {
		validated, err := run.g.Validate(run.code)
		if err != nil || validated != run.exp {
			t.Errorf("code %s should be %s got %s %s", run.code, run.exp, validated, err)
		}
	}
}

func TestPattern(t *testing.T) {
	code := Generate()
	matched, _ := regexp.MatchString(`^[0-9A-Z-]+$`, code)
	if !matched {
		t.Error("should only contain uppercase letters, digits, and dashes")
	}
	matched, _ = regexp.MatchString(`^\w{4}-\w{4}-\w{4}$`, code)
	if !matched {
		t.Error("should look like XXXX-XXXX-XXXX")
	}
	g := New(2, 5)
	code = g.Generate()
	matched, _ = regexp.MatchString(`^\w{5}-\w{5}$`, code)
	if !matched {
		t.Error("should generate an arbitrary number of parts")
	}
}

func TestDefaultSelfContained(t *testing.T) {
	for i := 0; i < 10; i++ {
		code := Generate()
		validated, err := Validate(code)
		if err != nil {
			t.Errorf("generated %s got %s %s", code, validated, err)
		}
	}
}

func TestCustomSelfContained(t *testing.T) {
	g := New(4, 6)
	for i := 0; i < 10; i++ {
		code := g.Generate()
		validated, err := g.Validate(code)
		if err != nil {
			t.Errorf("generated %s got %s %s", code, validated, err)
		}
	}
}
*/
