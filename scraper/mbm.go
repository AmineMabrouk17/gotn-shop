package scraper

import (
	"fmt"
	"strings"
	"time"

	"github.com/AmineMabrouk17/gotn-shop/model"
	"github.com/gocolly/colly/v2"
)

// MBMScraper targets mbm-tn.com, a PrestaShop site using the "electron" theme.
//
// Unlike MyTek and Scoop, this site serves plain server-rendered HTML with no
// Cloudflare challenge, and its robots.txt does not contain a blanket
// "Disallow: /*?" rule. Category paths and "?p=N" pagination are permitted.
type MBMScraper struct{}

func NewMBM() *MBMScraper {
	return &MBMScraper{}
}

func (s *MBMScraper) Name() string {
	return "MBM"
}

func (s *MBMScraper) Scrape(categoryURL string, maxPages int) ([]model.Product, error) {
	var products []model.Product
	itemCounter := 1
	if maxPages < 1 {
		maxPages = 1
	}

	c := colly.NewCollector(
		colly.UserAgent("gotn-shop-bot/1.0 (+https://github.com/AmineMabrouk17/gotn-shop)"),
		colly.AllowedDomains("mbm-tn.com", "www.mbm-tn.com"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*mbm-tn.com*",
		Parallelism: 1,
		Delay:       1 * time.Second,
	})

	visitedPages := 0
	seen := map[string]bool{categoryURL: true}

	c.OnHTML("article.product-miniature", func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.ChildText(".product-title a"))
		link := e.ChildAttr(".product-title a", "href")
		// This theme renders the card in a desktop grid and a mobile grid, so
		// .price matches several times per card. ChildText would concatenate
		// them and break parsing, so take the first match only.
		priceText := strings.TrimSpace(e.DOM.Find(".price").First().Text())
		image := e.ChildAttr(".tvproduct-image img", "src")

		// There is no .stock element on this theme: availability is encoded on
		// the add-to-cart button (class "tvproduct-out-of-stock" + disabled).
		btnClass := e.ChildAttr("button.add-to-cart", "class")
		btnTitle := strings.ToLower(e.ChildAttr("button.add-to-cart", "data-original-title"))
		inStock := !strings.Contains(btnClass, "tvproduct-out-of-stock") &&
			!strings.Contains(btnTitle, "out of stock")

		if title != "" && priceText != "" {
			products = append(products, model.Product{
				ID:        itemCounter,
				Store:     s.Name(),
				Title:     title,
				PriceText: priceText,
				PriceTND:  model.CleanPrice(priceText),
				Image:     image,
				Link:      link,
				InStock:   inStock,
			})
			itemCounter++
		}
	})

	// The theme loads more products with JS infinite scroll and renders its
	// "next" anchor only at runtime, so follow the <link rel="next"> the server
	// does emit instead.
	c.OnHTML(`link[rel="next"]`, func(e *colly.HTMLElement) {
		if visitedPages < maxPages-1 {
			next := e.Attr("href")
			// The page ships more than one rel="next" link pointing at the same
			// page (?page=2 and ?p=2), so de-duplicate before following.
			if next != "" && !seen[next] {
				seen[next] = true
				visitedPages++
				fmt.Printf("  -> page %d\n", visitedPages+1)
				e.Request.Visit(next)
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

	err := c.Visit(categoryURL)
	return products, err
}
