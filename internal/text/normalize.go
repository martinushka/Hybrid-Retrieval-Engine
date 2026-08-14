package text

import (
	"strings"
	"unicode"
)

func Normalize(s string) string {
	s = strings.ToLower(s)

	var b strings.Builder
	b.Grow(len(s))

	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			b.WriteRune(r)
			continue
		}

		b.WriteRune(' ')
	}

	return strings.Join(strings.Fields(b.String()), " ")
}

func Tokens(s string) []string {
	normalized := Normalize(s)

	if normalized == "" {
		return nil
	}

	return strings.Fields(normalized)
}
