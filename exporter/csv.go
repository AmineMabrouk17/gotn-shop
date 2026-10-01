package exporter

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/AmineMabrouk17/gotn-shop/model"
)

func ExportToCSV(filename string, products []model.Product) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header row
	headers := []string{"ID", "Store", "Title", "Price (DT)", "Raw Price", "In Stock", "Link"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	for _, p := range products {
		row := []string{
			strconv.Itoa(p.ID),
			p.Store,
			p.Title,
			fmt.Sprintf("%.3f", p.PriceTND),
			p.PriceText,
			strconv.FormatBool(p.InStock),
			p.Link,
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	// csv.Writer buffers internally: real I/O errors only surface here, so this
	// check is required for ExportToCSV to report a failed write honestly.
	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}

	return file.Sync()
}
