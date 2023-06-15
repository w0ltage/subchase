package main

import (
	"fmt"
	"github.com/gocolly/colly"
)

func main() {
    domain := "google.com"
    query := "https://www.google.com/search?q=site:" + domain
    
    c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36" 

    c.OnHTML("#center_col cite.apx8Vc", func(e *colly.HTMLElement) {
        domSelection := e.DOM
        fmt.Println(domSelection.Contents().First().Text())
    })

    err := c.Visit(query)
    if err != nil {
        fmt.Println("Scraping error: ", err)
    }
}
