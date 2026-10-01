package main

import (
	"fmt"

	"github.com/AmineMabrouk17/gotn-shop/exporter"
	"github.com/AmineMabrouk17/gotn-shop/model"
	"github.com/AmineMabrouk17/gotn-shop/scraper"
)

type target struct {
	scraper scraper.StoreScraper
	url     string
	pages   int
}

func main() {
	fmt.Println("========================================")
	fmt.Println("         gotn-shop: Store Scraper       ")
	fmt.Println("========================================")

	// Each store is a StoreScraper, so adding one is a single line here.
	targets := []target{
		{scraper.NewTunisiaNet(), "https://www.tunisianet.com.tn/301-pc-portable-tunisie", 2},
		{scraper.NewMBM(), "https://mbm-tn.com/145-pc-portable", 2},
	}

	var all []model.Product

	for _, t := range targets {
		fmt.Printf("\nFetching from %s (Pages: %d)...\n", t.scraper.Name(), t.pages)
		products, err := t.scraper.Scrape(t.url, t.pages)
		if err != nil {
			fmt.Printf("  %s failed: %v\n", t.scraper.Name(), err)
			continue
		}
		fmt.Printf("  %s: scraped %d products\n", t.scraper.Name(), len(products))
		all = append(all, products...)
	}

	fmt.Printf("\nTotal: %d products across %d stores.\n", len(all), len(targets))

	csvFile := "products.csv"
	if err := exporter.ExportToCSV(csvFile, all); err != nil {
		fmt.Printf("Failed to export CSV: %v\n", err)
		return
	}
	fmt.Printf("Exported results to %s\n", csvFile)
}
