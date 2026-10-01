package scraper

import (
	"fmt"
	"strings"
	"time"

	"github.com/AmineMabrouk17/gotn-shop/model"
	"github.com/gocolly/colly/v2"
)

type TunisiaNetScraper struct{}

func NewTunisiaNet() *TunisiaNetScraper {
	return &TunisiaNetScraper{}
}

func (s *TunisiaNetScraper) Name() string {
	return "TunisiaNet"
}

func (s *TunisiaNetScraper) Scrape(url string, maxPages int) ([]model.Product, error) {
	var products []model.Product
	itemCounter := 1

	c := colly.NewCollector(
		colly.UserAgent("gotn-shop-bot/1.0 (+https://github.com/AmineMabrouk17/gotn-shop)"),
		colly.AllowedDomains("www.tunisianet.com.tn", "tunisianet.com.tn"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*tunisianet.*",
		Parallelism: 1,
		Delay:       1 * time.Second,
	})

	currentPage := 1

	c.OnHTML(".item-product", func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.ChildText(".product-title a"))
		// Each card renders the price twice (desktop + mobile), so take only the
		// first match. ChildText concatenates ALL matches, which yields
		// "1 169,000 DT1 169,000 DT" and fails to parse as 0.
		priceText := strings.TrimSpace(e.DOM.Find(".price").First().Text())
		link := e.ChildAttr(".product-title a", "href")
		img := e.ChildAttr(".thumbnail-container img", "data-full-size-image-url")
		if img == "" {
			img = e.ChildAttr(".thumbnail-container img", "src")
		}
		// Stock lives in #stock_availability / .in-stock on this theme; there is
		// no .product-availability element, so matching it alone reports every
		// product as in stock.
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
				Store:     s.Name(),
				Title:     title,
				PriceText: priceText,
				PriceTND:  model.CleanPrice(priceText),
				Image:     img,
				Link:      link,
				InStock:   inStock,
			})
			itemCounter++
		}
	})

	c.OnHTML("a.next", func(e *colly.HTMLElement) {
		if currentPage < maxPages {
			nextURL := e.Attr("href")
			if nextURL != "" {
				currentPage++
				fmt.Printf("  -> page %d\n", currentPage)
				e.Request.Visit(nextURL)
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("  fetching: %s\n", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		status := 0
		if r != nil {
			status = r.StatusCode
		}
		fmt.Printf("  request failed (status %d): %v\n", status, err)
	})

	if maxPages < 1 {
		maxPages = 1
	}
	err := c.Visit(url)
	return products, err
}
