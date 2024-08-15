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
		colly.AllowedDomains("www.scrapingcourse.com"),
	)
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"
	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		product := product{}
		product.Url = e.ChildAttr("a", "href")
		product.Image = e.ChildAttr("img", "src")
		product.Name = e.ChildText(".product-name")
		product.Price = e.ChildText(".price")

		// add the product instance with scraped data to the list of products
		products = append(products, product)
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting:", r.URL)
	})
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong: ", err)
	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited: ", r.Request.URL)
	})
	c.OnHTML("a.next", func(e *colly.HTMLElement) {

		// extract the next page URL from the next button
		nextPage := e.Attr("href")

		// check if the nextPage URL has been visited
		if _, found := visitedUrls.Load(nextPage); !found {
			fmt.Println("scraping:", nextPage)
			// mark the URL as visited
			visitedUrls.Store(nextPage, struct{}{})
			// visit the next page
			e.Request.Visit(nextPage)
		}
	})
	// c.OnHTML("a", func(e, *colly.HTMLElement) {
	// 	fmt.Println("%v", e.Attr("href"))
	// })
	c.OnScraped(func(r *colly.Response) {
		file, err := os.Create("products.csv")
		if err != nil {
			log.Fatalln("failed to create output CSV file", err)
		}
		defer file.Close()
		writer := csv.NewWriter(file)
		headers := []string{
			"Url",
			"Image",
			"Name",
			"Price",
		}
		writer.Write(headers)
		for _, product := range products {
			record := []string{
				product.Url,
				product.Image,
				product.Name,
				product.Price,
			}
			writer.Write(record)

		}
		defer writer.Flush()
		fmt.Println(r.Request.URL, "Scrapped!")
	})
	c.Visit("https://www.scrapingcourse.com/ecommerce")

}