package model

import (
	"strconv"
	"strings"
	"unicode"
)

type Product struct {
	ID        int     `json:"id"`
	Store     string  `json:"store"`
	Title     string  `json:"title"`
	PriceText string  `json:"price_text"`
	PriceTND  float64 `json:"price_tnd"`
	Link      string  `json:"link"`
	InStock   bool    `json:"in_stock"`
}

// CleanPrice parses Tunisian price strings like "1 299,000 DT" into a float64 (1299.000)
func CleanPrice(raw string) float64 {
	// 1. Remove currency labels and every kind of space (regular, non-breaking
	//    U+00A0, and narrow non-breaking U+202F used as the thousands separator)
	cleaned := strings.ReplaceAll(raw, "DT", "")
	cleaned = strings.ReplaceAll(cleaned, "TND", "")
	cleaned = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, cleaned)

	// 2. Tunisian sites use comma (',') for decimals -> replace with dot ('.')
	cleaned = strings.ReplaceAll(cleaned, ",", ".")

	// 3. Convert to float
	price, err := strconv.ParseFloat(strings.TrimSpace(cleaned), 64)
	if err != nil {
		return 0.0
	}
	return price
}
