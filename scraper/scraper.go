package scraper

import "github.com/AmineMabrouk17/gotn-shop/model"

// StoreScraper represents any store source (web scraper, API, or local export)
type StoreScraper interface {
	Name() string
	Scrape(categoryOrQuery string, maxPages int) ([]model.Product, error)
}
