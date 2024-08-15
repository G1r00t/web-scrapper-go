package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/gocolly/colly"
)

type product struct {
	Url, Image, Name, Price string
}

func main() {
	var products []product
	var visitedUrls sync.Map

	c := colly.NewCollector(
		colly.AllowedDomains("www.amazon.com"),
	)
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"

	// Adjust selectors according to Amazon's page structure
	c.OnHTML("div.s-main-slot div.s-result-item", func(e *colly.HTMLElement) {
		product := product{}
		product.Url = e.ChildAttr("a.a-link-normal", "href")
		product.Image = e.ChildAttr("img.s-image", "src")
		product.Name = e.ChildText("span.a-text-normal")
		product.Price = e.ChildText("span.a-price-whole")

		if product.Name != "" && product.Price != "" {
			// Append the product instance with scraped data to the list of products
			products = append(products, product)
			fmt.Printf("Scraped Product: %+v\n", product) // Debug output
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited:", r.Request.URL)
	})

	c.OnHTML("a.s-pagination-next", func(e *colly.HTMLElement) {
		// Extract the next page URL from the next button
		nextPage := e.Attr("href")

		// Check if the nextPage URL has been visited
		if _, found := visitedUrls.Load(nextPage); !found {
			fmt.Println("Scraping next page:", nextPage)
			// Mark the URL as visited
			visitedUrls.Store(nextPage, struct{}{})
			// Visit the next page
			e.Request.Visit(e.Request.AbsoluteURL(nextPage))
		}
	})

	c.OnScraped(func(r *colly.Response) {
		file, err := os.Create("amazonproducts.csv")
		if err != nil {
			log.Fatalln("Failed to create output CSV file", err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		headers := []string{"Url", "Image", "Name", "Price"}
		if err := writer.Write(headers); err != nil {
			log.Fatalln("Failed to write CSV headers", err)
		}

		for _, product := range products {
			record := []string{
				product.Url,
				product.Image,
				product.Name,
				product.Price,
			}
			if err := writer.Write(record); err != nil {
				log.Fatalln("Failed to write CSV record", err)
			}
		}
		writer.Flush()
		fmt.Println("Scraping completed. Data saved to amazonproducts.csv")
	})

	// Start scraping
	err := c.Visit("https://www.amazon.com/s?k=tv+stand&i=todays-deals&crid=I84UKWPDHOSE&sprefix=tv+stand%2Ctodays-deals%2C321&ref=nb_sb_noss_1")
	if err != nil {
		log.Fatalln("Failed to start scraping", err)
	}
}
