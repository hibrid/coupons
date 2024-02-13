package generator

import (
	"crypto/sha256"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

// return random int in the range min...max
func randomInt(min, max int) int {
	return min + rand.Intn(1+max-min)
}

// return random char string from charset
func randomChar(cs []byte) string {
	return string(cs[randomInt(0, len(cs)-1)])
}

// repeat string with one str (#)
func repeatStr(count uint16, str string) string {
	return strings.Repeat(str, int(count))
}

func numberOfChar(str, char string) uint32 {
	return uint32(strings.Count(str, char))
}

func isFeasible(charset, pattern, char string, count uint32) bool {
	ls := numberOfChar(pattern, char)
	return math.Pow(float64(len(charset)), float64(ls)) >= float64(count)
}

func checkCharacter(code string, check int) string {
	code = convertSpecialLetters(code)
	for _, r := range code {
		k := strings.IndexRune(alphanumeric, r)
		check = check*19 + k
	}
	return string(alphanumeric[check%int(length)])
}

func secureCheckCharacter(input string, check int, sha256Index int) string {
	// First, convert special letters in the input to standardize it
	standardInput := convertSpecialLetters(input)

	// Incorporate the check integer into the input string before hashing
	inputWithCheck := standardInput + strconv.Itoa(check)

	// Then, compute the SHA-256 hash of the standardized input
	hash := sha256.Sum256([]byte(inputWithCheck))

	// For simplicity, use the first byte of the hash to determine the character
	// This byte is effectively a number between 0 and 255
	indexByte := hash[sha256Index]

	// Use modulo operation with the size of the alphanumeric set to get a valid index
	// Assuming the alphanumeric set is 0-9 + A-Z, giving 36 characters
	index := indexByte % 36

	// Map the index to a character in the alphanumeric set
	// This assumes '0'-'9' (10 chars) followed by 'A'-'Z' (26 chars) in the alphanumeric set
	var checkChar rune
	if index < 10 {
		// Map to '0'-'9'
		checkChar = rune('0' + index)
	} else {
		// Map to 'A'-'Z'
		checkChar = rune('A' + index - 10)
	}

	return string(checkChar)
}

func convertSpecialLetters(code string) string {
	replacements := map[string]string{
		"O": "0",
		"I": "1",
		"Z": "2",
		"S": "5",
		"L": "1",
		"G": "6",
		"B": "8",
		"Q": "0",
		"T": "7",
	}

	for from, to := range replacements {
		code = strings.Replace(code, from, to, -1)
	}

	return code
}

// Helper function to check if the string is alphanumeric
func isAlphanumeric(s string) bool {
	for _, r := range s {
		if !('0' <= r && r <= '9') && !('A' <= r && r <= 'Z') {
			return false
		}
	}
	return true
}
