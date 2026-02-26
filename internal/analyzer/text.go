package analyzer

import "unicode"

func beginsWithLowercase(s string) bool {
	if s == "" {
		return true
	}
	r := []rune(s)[0]
	return unicode.IsLower(r)
}

func containsNonASCII(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return true
		}
	}
	return false
}

func containsSymbols(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			continue
		}
		return true
	}
	return false
}