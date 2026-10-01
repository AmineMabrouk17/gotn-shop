package scraper

import (
	"fmt"
	"strings"
	"time"

	"github.com/AmineMabrouk17/gotn-shop/model"
	"github.com/gocolly/colly/v2"
)

func ScrapeTunisiaNet(categoryURL string, maxPages int) ([]model.Product, error) {
	var products []model.Product
	itemCounter := 1

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
		colly.AllowedDomains("www.tunisianet.com.tn", "tunisianet.com.tn"),
	)

	// Rate limiting: 1 second delay between requests
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*tunisianet.*",
		Parallelism: 1,
		Delay:       1 * time.Second,
	})

	currentPage := 1

	// Scrape each product card
	c.OnHTML(".item-product", func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.ChildText(".product-title a"))
		// The card renders the price twice (desktop + mobile), so take only the
		// first match instead of concatenating every match like ChildText does.
		priceText := strings.TrimSpace(e.DOM.Find(".price").First().Text())
		link := e.ChildAttr(".product-title a", "href")
		// Stock lives in #stock_availability/.in-stock, not .product-availability.
		stockText := strings.ToLower(strings.TrimSpace(
			e.DOM.Find(".product-availability, .in-stock, .out-of-stock, #stock_availability").First().Text()))

		inStock := true
		if stockText != "" {
			for _, out := range []string{"épuisé", "epuise", "rupture", "hors stock", "indisponible", "non disponible", "out of stock"} {
				if strings.Contains(stockText, out) {
					inStock = false
					break
				}
			}
		}

		if title != "" && priceText != "" {
			products = append(products, model.Product{
				ID:        itemCounter,
				Store:     "TunisiaNet",
				Title:     title,
				PriceText: priceText,
				PriceTND:  model.CleanPrice(priceText),
				Link:      link,
				InStock:   inStock,
			})
			itemCounter++
		}
	})

	// Handle pagination up to maxPages
	c.OnHTML("a.next", func(e *colly.HTMLElement) {
		if currentPage < maxPages {
			nextURL := e.Attr("href")
			if nextURL != "" {
				currentPage++
				fmt.Printf("--> Moving to page %d...\n", currentPage)
				e.Request.Visit(nextURL)
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Fetching: %s\n", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Request URL: %s failed with response: %v\nError: %v\n", r.Request.URL, r, err)
	})

	err := c.Visit(categoryURL)
	return products, err
}
