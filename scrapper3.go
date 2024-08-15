package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/gocolly/colly"
)

// product structure to keep the scraped data
type Product struct {
	Url, Image, Name, Price string
}

func main() {
	pagesToScrape := []string{
		"your-website/page/1/",
		"yourwebsite/page/2/",
		
	}

	// instantiate a new collector object
	c := colly.NewCollector(
		colly.AllowedDomains("your-website"),
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		// limit the parallel requests to 4 request at a time
		Parallelism: 4,
	})

	// initialize the slice of structs that will contain the scraped data
	var products []Product

	// set a global User Agent
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"

	// set up the proxy
	err := c.SetProxy("http://35.185.196.38:3128")
	if err != nil {
		log.Fatal(err)
	}

	// OnHTML callback for scraping product information
	c.OnHTML("li.product", func(e *colly.HTMLElement) {

		// initialize a new Product instance
		product := Product{}

		// scrape the target data
		product.Url = e.ChildAttr("a", "href")
		product.Image = e.ChildAttr("img", "src")
		product.Name = e.ChildText(".product-name")
		product.Price = e.ChildText(".price")

		// add the product instance with scraped data to the list of products
		products = append(products, product)
	})

	// register all pages to scrape
	for _, pageToScrape := range pagesToScrape {
		c.Visit(pageToScrape)

		// store the data to a CSV after extraction
		c.OnScraped(func(r *colly.Response) {

			// open the CSV file
			file, err := os.Create("products.csv")
			if err != nil {
				log.Fatalln("Failed to create output CSV file", err)
			}
			defer file.Close()

			// initialize a file writer
			writer := csv.NewWriter(file)

			// write the CSV headers
			headers := []string{
				"Url",
				"Image",
				"Name",
				"Price",
			}
			writer.Write(headers)

			// write each product as a CSV row
			for _, product := range products {
				// convert a Product to an array of strings
				record := []string{
					product.Url,
					product.Image,
					product.Name,
					product.Price,
				}

				// add a CSV record to the output file
				writer.Write(record)
			}
			writer.Flush()
		})
	}

	// wait for Colly to visit all pages
	c.Wait()
}
