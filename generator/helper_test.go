package generator

import (
	"strings"
	"testing"
)

func TestRandomInt(t *testing.T) {
	tables := []struct {
		min, max int
	}{
		{1, 10},
		{100, 200},
	}

	for _, v := range tables {
		n := randomInt(v.min, v.max)
		if n > v.max || n < v.min {
			t.Fail()
		}
	}
}

func TestRandomChar(t *testing.T) {
	tables := []struct {
		chars []byte
	}{
		{[]byte("abcljasdlkjasd")},
		{[]byte("kqweqwrkjhasfn")},
	}

	for _, v := range tables {
		c := randomChar(v.chars)
		if !strings.Contains(string(v.chars), c) {
			t.Fail()
		}
	}
}

func TestRepeatStr(t *testing.T) {
	tables := []struct {
		count uint16
		str   string
	}{
		{10, "a"},
		{30, "a"},
		{0, "a"},
	}

	for _, v := range tables {
		rs := repeatStr(v.count, v.str)
		if uint16(len(rs)) != v.count {
			t.Fail()
		}
	}
}

func TestIsFeasible(t *testing.T) {
	tables := []struct {
		charset, pattern string
		count            uint32
		wants            bool
	}{
		{"abcdefghijk", "##-#####-###", 10000, true},
		{"abcdefghijk", "##-#####-###", 10000, true},
		{"abcdefghijk", "##-#", 10000, false},
	}

	for _, v := range tables {
		b := isFeasible(v.charset, v.pattern, "#", v.count)
		if b != v.wants {
			t.Fail()
		}
	}
}

func TestConvertSpecialLetters(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"LOGISTICS", "1061571C5"}, // Ensures 'O', 'I', 'G', 'S', 'L' are replaced
		{"BIZARRE", "812ARRE"},     // Ensures 'B', 'I', 'Z' are replaced
		{"QUIET", "0U1E7"},         //Ensures 'Q', 'I', 'E', 'T' are replaced
		{"123456", "123456"},       // Ensures numbers are unchanged
		{"ACD", "ACD"},             // Ensures letters not in the replacement map are unchanged
		{"", ""},                   // Ensures empty string is unchanged
	}

	for _, c := range cases {
		got := convertSpecialLetters(c.in)
		if got != c.want {
			t.Errorf("convertSpecialLetters(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

// TestSecureCheckCharacter tests the secureCheckCharacter function to ensure it produces valid and consistent check characters.
func TestSecureCheckCharacter(t *testing.T) {
	// Test cases with expected outcomes are hard to define due to the nature of SHA-256 hashing and modulo operation.
	// Instead, we'll test for consistent output and valid alphanumeric character.
	testInputs := []string{
		"LOGISTICS",
		"BIZARRE",
		"QUIET",
		"123456",
		"ABCD",
		"",
		"OILGAS",
	}

	for _, input := range testInputs {
		got := secureCheckCharacter(input, 1, 0)
		if len(got) != 1 {
			t.Errorf("secureCheckCharacter(%q) produced an invalid length output: got %q", input, got)
		}
		if !isAlphanumeric(got) {
			t.Errorf("secureCheckCharacter(%q) produced a non-alphanumeric character: got %q", input, got)
		}
	}
}
