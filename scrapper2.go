
package main

import (
    "fmt"
    "log"

    "github.com/gocolly/colly"
)

func main() {
    // instantiate a new collector object
    c := colly.NewCollector(
        colly.AllowedDomains("www.g2.com"),
    )

    // set a global User Agent
    c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"

    // set up the proxy
    err := c.SetProxy("http://35.185.196.38:3128")
    if err != nil {
        log.Fatal(err)
    }

    // OnError callback
    c.OnError(func(_ *colly.Response, err error) {
        log.Println("Something went wrong:", err)
    })

    // OnResponse callback to print the full HTML
    c.OnResponse(func(r *colly.Response) {
        fmt.Println("Page visited:", r.Request.URL.String())
        fmt.Println("Full HTML:\n", string(r.Body))
    })

    // OnScraped callback
    c.OnScraped(func(r *colly.Response) {
        fmt.Println("Finished scraping:", r.Request.URL.String())
    })

	    // open the target URL
	    c.Visit("https://www.g2.com/products/asana/reviews")
	}
