package model

import (
	"strconv"
	"strings"
	"unicode"
)

type Product struct {
	ID            int     `json:"id"`
	Store         string  `json:"store"`
	Title         string  `json:"title"`
	SKU           string  `json:"sku,omitempty"`
	PriceText     string  `json:"price_text"`
	PriceTND      float64 `json:"price_tnd"`
	OriginalPrice float64 `json:"original_price_tnd,omitempty"`
	Image         string  `json:"image,omitempty"`
	Link          string  `json:"link"`
	InStock       bool    `json:"in_stock"`
}

// CleanPrice parses Tunisian price strings like "1 299,000 DT" or "1\u202f169,000 DT"
func CleanPrice(raw string) float64 {
	// 1. Remove currency labels
	cleaned := strings.ReplaceAll(raw, "DT", "")
	cleaned = strings.ReplaceAll(cleaned, "TND", "")

	// 2. Strip ALL unicode whitespace (including \u00a0, \u202f, standard space, etc.)
	var b strings.Builder
	for _, r := range cleaned {
		if !unicode.IsSpace(r) {
			b.WriteRune(r)
		}
	}
	cleaned = b.String()

	// 3. Normalize comma to dot for decimal representation
	cleaned = strings.ReplaceAll(cleaned, ",", ".")

	// 4. Parse float
	price, err := strconv.ParseFloat(strings.TrimSpace(cleaned), 64)
	if err != nil {
		return 0.0
	}
	return price
}
