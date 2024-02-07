package generator

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	// charset types
	numbers         = "0123456789"
	alphabetic      = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphanumeric    = numbers + alphabetic
	length          = len(alphanumeric) - 1
	removeInvalidRe = regexp.MustCompile(`[^0-9A-Z]+`)

	// defaults
	minCodesCount uint32 = 1
	minLength     uint32 = 6
)

const patternChar = "#"
const divider = "-"

var (
	ErrNotFeasible       = errors.New("not feasible to generate requested number of codes")
	ErrPatternIsNotMatch = errors.New("pattern is not match with the length value")
)

// initialize random seed
func init() {
	//rand.Seed(time.Now().UnixNano())
}

type Generator struct {
	// MinimumLength of the code
	MinimumLength uint32 `json:"length"`

	// Pattern character
	PatternCharacter string `json:"pattern_character"`

	// Pattern Divider
	PatternDivider string `json:"pattern_divider"`

	// Count of the codes
	Count uint32 `json:"count"`

	// Charset to use
	Charset string `json:"charset"`

	// Prefix of the code
	Prefix string `json:"prefix"`

	// Suffix of the code
	Suffix string `json:"suffix"`

	// Pattern of the code
	Pattern string `json:"pattern"`
}

// Creates a new generator with options
func NewWithOptions(opts ...Option) (*Generator, error) {
	g := Default()
	if err := setOptions(opts...)(g); err != nil {
		return nil, err
	}

	return g, nil
}

// Creates a new generator with default values
func Default() *Generator {
	return &Generator{
		MinimumLength:    minLength,
		PatternCharacter: patternChar,
		PatternDivider:   divider,
		Count:            minCodesCount,
		Charset:          alphanumeric,
		Pattern:          "####-####-####",
	}
}

// Generates a list of codes
func (g *Generator) Run() ([]string, error) {
	if !isFeasible(g.Charset, g.Pattern, patternChar, g.Count) {
		return nil, ErrNotFeasible
	}

	result := make([]string, g.Count)

	var i uint32
	for i = 0; i < g.Count; i++ {
		code := g.one()
		if !hasBadWord(code) {
			result[i] = code
		} else {
			i-- // Re-generate the code if it contains a bad word
		}
	}

	return result, nil
}

func (g *Generator) CountUniqueCombinations(pattern, patternChar, alphanumeric string) int {
	// Count occurrences of the pattern character
	occurrences := strings.Count(g.Pattern, g.PatternCharacter)

	// Get the size of the alphanumeric set
	alphanumericSetSize := len(alphanumeric)

	// Calculate the number of combinations: alphanumericSetSize^occurrences
	combinations := 1
	for i := 0; i < occurrences; i++ {
		combinations *= alphanumericSetSize
	}

	return combinations
}

func (g *Generator) calculateLengths() []int {
	parts := strings.Split(g.Pattern, g.PatternDivider)
	lengths := make([]int, len(parts))

	for i, part := range parts {
		lengths[i] = strings.Count(part, g.PatternCharacter)
	}

	return lengths
}

// one generates one code
func (g *Generator) one() string {
	partLengths := g.calculateLengths()
	parts := make([]string, len(partLengths))
	for i, v := range partLengths {
		var code string
		for j := 0; j < v-1; j++ {
			code += randomChar([]byte(g.Charset))
		}
		check := checkCharacter(code, i+1)
		parts[i] = code + check
		if !hasBadWord(strings.Join(parts, "")) {
			i += 1
		}
	}

	suffix := g.Suffix

	return strings.Join(parts, g.PatternDivider) + suffix
}

func (g *Generator) Validate(code string) (string, error) {
	// Calculate the Lengths for each part based on the Pattern
	partLengths := g.calculateLengths()

	// make uppercase
	code = strings.ToUpper(code)

	// remove invalid characters
	code = removeInvalidRe.ReplaceAllLiteralString(code, "")

	// convert special letters to numbers
	//code = convertSpecialLetters(code)

	// split into parts
	parts := []string{}
	tmp := code
	for _, length := range partLengths {
		max := length
		if max > len(tmp) {
			max = len(tmp)
		}
		parts = append(parts, tmp[:max])
		tmp = tmp[max:]
	}

	// join with separator (shouldn't we test that)
	code = strings.Join(parts, g.PatternDivider)

	if len(parts) != len(partLengths) {
		return code, fmt.Errorf("wrong number of parts (%d)", len(parts))
	}
	for i, part := range parts {
		check := checkCharacter(part[:len(part)-1], i+1)
		if !strings.HasSuffix(part, check) {
			return code, fmt.Errorf("wrong part %d (%s) check character %s", i+1, part, check)
		}
	}

	return code, nil
}
