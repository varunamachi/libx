package str

import (
	"strings"
	"unicode"
)

func EqFold(str string, targets ...string) bool {
	for _, s := range targets {
		if strings.EqualFold(str, s) {
			return true
		}
	}
	return false
}

func RemoveSpaces(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}
