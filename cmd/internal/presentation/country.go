package presentation

import (
	"net/http"
	"regexp"
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

	// Remove special characters and numbers, keep spaces
	regx := regexp.MustCompile(`[^\p{L}\s]`)
	normalizedString = regx.ReplaceAllString(normalizedString, "")

	return normalizedString
}
