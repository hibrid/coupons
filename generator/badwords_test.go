package generator

import (
	"testing"
)

func TestBadWordsNegative(t *testing.T) {
	if hasBadWord("LOVE") {
		t.Error("LOVE shouldn't be a bad word")
	}
}

func TestBadWordsPositive(t *testing.T) {
	if !hasBadWord("BOOBIES") {
		t.Error("BOOB should be a bad word")
	}
}

func TestBadWordsLowerCase(t *testing.T) {
	if !hasBadWord("Boobies") {
		t.Error("BOOB should be a bad word")
	}
}
