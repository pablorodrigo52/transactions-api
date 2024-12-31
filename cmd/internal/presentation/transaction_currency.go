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

type TransactionCurrencyDTO struct {
	TransactionID           int64   `json:"transaction_id"`
	Description             string  `json:"description"`
	TransactionDate         string  `json:"transaction_date"`
	PurchaseAmount          float32 `json:"purchase_amount"`
	ExchangeRate            float32 `json:"exchange_rate"`
	ConvertedPurchaseAmount float32 `json:"converted_purchase_amount"`
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
