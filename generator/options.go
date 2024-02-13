package generator

import "fmt"

// Define a base error type for validation errors
type ValidationError struct {
	Reason string
}

// Implement the Error() method for ValidationError
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s", e.Reason)
}

// Implement the Is method for comparing errors
func (e *ValidationError) Is(target error) bool {
	t, ok := target.(*ValidationError)
	return ok && t.Reason == e.Reason
}

// Define specific error types for each error variable
type ErrInvalidCountType struct{ ValidationError }
type ErrInvalidCharsetType struct{ ValidationError }
type ErrInvalidPatternType struct{ ValidationError }
type ErrPrefixBadWordType struct{ ValidationError }
type ErrSuffixBadWordType struct{ ValidationError }
type ErrInvalidSHA256IndexType struct{ ValidationError }

// Initialize the specific error variables with corresponding messages
var (
	ErrInvalidCount       = &ErrInvalidCountType{ValidationError{"invalid count, It should be greater than 0"}}
	ErrInvalidCharset     = &ErrInvalidCharsetType{ValidationError{"invalid charset, charset length should be greater than 0"}}
	ErrInvalidPattern     = &ErrInvalidPatternType{ValidationError{"invalid pattern, pattern cannot be empty"}}
	ErrPrefixBadWord      = &ErrPrefixBadWordType{ValidationError{"invalid prefix, prefix contains bad word"}}
	ErrSuffixBadWord      = &ErrSuffixBadWordType{ValidationError{"invalid suffix, suffix contains bad word"}}
	ErrInvalidSHA256Index = &ErrInvalidSHA256IndexType{ValidationError{"invalid SHA256 index, must be between 0 and 31"}}
)

type Option func(*Generator) error

// SetLength sets the length of the code
func SetMinimumLength(length uint32) Option {
	return func(g *Generator) error {
		if length == 0 {
			length = numberOfChar(g.Pattern, patternChar)
		}
		g.MinimumLength = length
		return nil
	}
}

// SetPatternCharacter sets the pattern character of the code
func SetPatternCharacter(patternCharacter string) Option {
	return func(g *Generator) error {
		g.PatternCharacter = patternCharacter
		return nil
	}
}

// SetPatternDivider sets the pattern divider of the code
func SetPatternDivider(patternDivider string) Option {
	return func(g *Generator) error {
		g.PatternDivider = patternDivider
		return nil
	}
}

// SetGenerateCount sets the count of the code
func SetGenerateCount(count uint32) Option {
	return func(g *Generator) error {
		if count == 0 {
			return ErrInvalidCount
		}
		g.Count = count
		return nil
	}
}

// SetCharset sets the charset of the code
func SetCharset(charset string) Option {
	return func(g *Generator) error {
		if len(charset) == 0 {
			return ErrInvalidCharset
		}
		g.Charset = charset
		return nil
	}
}

// SetPrefix sets the prefix of the code
func SetPrefix(prefix string) Option {
	if !hasBadWord(prefix) {
		return func(g *Generator) error {
			g.Prefix = prefix
			return nil
		}
	}

	return func(g *Generator) error {
		return ErrPrefixBadWord
	}
}

// SetSuffix sets the suffix of the code
func SetSuffix(suffix string) Option {
	if !hasBadWord(suffix) {
		return func(g *Generator) error {
			g.Suffix = suffix
			return nil
		}
	}
	return func(g *Generator) error {
		return ErrSuffixBadWord
	}
}

// SetPattern sets the pattern of the code
func SetPattern(pattern string) Option {
	return func(g *Generator) error {
		if pattern == "" {
			return ErrInvalidPattern
		}

		numPatternChar := numberOfChar(pattern, patternChar)
		if g.MinimumLength == 0 || g.MinimumLength != numPatternChar {
			g.MinimumLength = numPatternChar
		}

		g.Pattern = pattern
		return nil
	}
}

// SetSHA256 index
func SetCheckCharacterSHA256Index(sha256Index uint32) Option {
	return func(g *Generator) error {
		if sha256Index >= 32 {
			return ErrInvalidSHA256Index
		}
		g.CheckCharacterSHA256Index = sha256Index
		return nil
	}
}

func setOptions(opts ...Option) Option {
	return func(g *Generator) error {
		for _, opt := range opts {
			if err := opt(g); err != nil {
				return err
			}
		}
		return nil
	}
}
