package util

import (
	"strings"
	"unicode"
)

// NormalizeString TBD
func NormalizeString(s string) string {
	return strings.ToLower(s)
}

// NormalizeTag TBD
func NormalizeTag(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return unicode.ToLower(r)
	}, s)
}
