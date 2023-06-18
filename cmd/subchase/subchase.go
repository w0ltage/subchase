package main

import (
	// "encoding/json"
	"fmt"
	"log"
	// "os"
	"reflect"

	"github.com/gocolly/colly"
)

type void struct{}

func main() {
    givenDomain := "google.com"
    query := "https://www.google.com/search?q=site:" + givenDomain
    
    rawDomains := collectDomains(query)
    domains := removeDuplicates(rawDomains)

    for i := 0; i < len(domains); i++ {
        domain := domains[i]
        fmt.Println(domain.Interface())
    }

    // enc := json.NewEncoder(os.Stdout)
    // enc.SetIndent("", "  ")
    // enc.Encode(rawDomains)

}

func collectDomains(query string) []string {
    var domains []string

    googleCollector := colly.NewCollector(
        colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"),
        colly.CacheDir("./google_cache"),
        )

    googleCollector.OnHTML("#pnnext[href]", func(e *colly.HTMLElement) {
        link := e.Attr("href")

        err := e.Request.Visit(link)
        if err != nil {
            fmt.Println("Scraping error: ", err)
        }
    })

    googleCollector.OnRequest(func(r *colly.Request) {
        log.Println("visiting", r.URL.String())
    })

    googleCollector.OnHTML("#center_col cite.apx8Vc", func(e *colly.HTMLElement) {
        domSelection := e.DOM
        link := domSelection.Contents().First().Text()
        domains = append(domains, link)
    })

    err := googleCollector.Visit(query)
    if err != nil {
        fmt.Println("Scraping error: ", err)
    }

    return domains
}

func removeDuplicates(domains []string) []reflect.Value {
    set := make(map[string]void)

    for _, element := range domains {
        set[element] = void{}
    }

    result := reflect.ValueOf(set).MapKeys()
    return result
}
