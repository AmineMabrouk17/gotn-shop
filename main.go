package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/AmineMabrouk17/gotn-shop/scraper"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("       gotn-shop: Tunisian Scraper      ")
	fmt.Println("========================================")

	// Example category: Laptops (PC Portables)
	targetURL := "https://www.tunisianet.com.tn/301-pc-portable-tunisie"
	maxPages := 2 // Scrape first 2 pages

	fmt.Printf("Scraping TunisiaNet (%s)...\n\n", targetURL)
	products, err := scraper.ScrapeTunisiaNet(targetURL, maxPages)
	if err != nil {
		fmt.Println("Error scraping:", err)
		return
	}

	fmt.Printf("\nDone! Collected %d products.\n", len(products))

	// Save to JSON
	fileData, err := json.MarshalIndent(products, "", "  ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	err = os.WriteFile("products.json", fileData, 0644)
	if err != nil {
		fmt.Println("Error saving file:", err)
		return
	}

	fmt.Println("Saved results to products.json!")
}
