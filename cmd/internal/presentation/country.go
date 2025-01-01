package presentation

import (
	"net/http"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Country string

func (c *Country) Validate() {
	if c == nil || *c == "" {
		panic(NewApiError(http.StatusBadRequest, "country name is required"))
	}
}

func (c *Country) Normalize() string {

	// Remove accents [ñ -> n], [ç -> c], [á -> a]
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	normalizedString, _, err := transform.String(t, string(*c))
	if err != nil {
		panic(NewApiError(http.StatusBadRequest, "country name not in pattern"))
	}
	
	return c.toTitleCase(normalizedString)
}

func (c *Country) toTitleCase(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			firstLetter := word[0]
			if firstLetter >= 65 && firstLetter < 122 {
				word := string(unicode.ToTitle(rune(firstLetter))) + word[1:]
				words[i] = word
			} else {
				break
			}
		}
	}
	return strings.Join(words, " ")
}
