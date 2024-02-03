package coupon

import (
	"math"
	"math/rand"
	"strings"
)

// return random int in the range min...max
func randomInt(min, max int) int {
	return min + rand.Intn(1+max-min)
}

func randString(n int, cs []byte) string {
	var bytes = make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = randomChar(cs)[0]
	}
	return string(bytes)
}

// return random char string from charset
func randomChar(cs []byte) string {
	return string(cs[randomInt(0, len(cs)-1)])
}

// repeat string with one str (#)
func repeatStr(count uint16, str string) string {
	return strings.Repeat(str, int(count))
}

func numberOfChar(str, char string) uint16 {
	return uint16(strings.Count(str, char))
}

func isFeasible(charset, pattern, char string, count uint16) bool {
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

func convertSpecialLetters(code string) string {
	replacements := map[string]string{
		"O": "0",
		"I": "1",
		"Z": "2",
		"S": "5",
	}

	for from, to := range replacements {
		code = strings.Replace(code, from, to, -1)
	}

	return code
}
