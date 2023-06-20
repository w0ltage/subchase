package main

import (
	"fmt"
	"log"
	"reflect"
    "strings"
    "net/http"
    "net/url"
    "crypto/tls"
    "flag"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

// The `void` type is defined as an empty struct.
// It is used as the value type for the map (`set`) to create
// a set-like data structure where only unique elements are stored.
type void struct{}

func main() {
    var givenDomain string

    flag.StringVar(&givenDomain, "d", "", "Pass hostname (ex: google.com)")
    flag.Parse()

    if givenDomain == "" {
        log.Fatalln("No hostname (domain) is passed. Exiting.")
    }
    
    rawDomains := findDomains(givenDomain)
    uniqDomains := processFoundDomains(rawDomains)

    // Iterate through slice of unique domains
    for i := 0; i < len(uniqDomains); i++ {
        domain := uniqDomains[i]
        fmt.Println(domain.Interface())
    }
}

func findDomains(givenDomain string) []string {
    var domains []string

    googleQuery := "https://www.google.com/search?q=site:" + givenDomain
    yandexQuery := "https://yandex.com/search/?text=site:" + givenDomain + "&lr=100&p=0"

    // Instantiate default collector
    collector := colly.NewCollector(
        // colly.CacheDir("./sites_cache"),
        colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/114.0"),
        )

    // Referer sets valid Referer HTTP header to requests
    extensions.Referer(collector)

    // Disable TLS security check for a client
    collector.WithTransport(&http.Transport{
        TLSClientConfig:&tls.Config{InsecureSkipVerify: true},
    })

    // For some unknown reason, requests get "Forbidden"
    // without using a localhost proxy (mitmproxy)
    err := collector.SetProxy("http://localhost:8080")
    if err != nil {
        log.Fatalln("Error happened with proxy:", err)
    }

    // Add headers to all requests
    collector.OnRequest(func(r *colly.Request) {
        // log.Println("visiting", r.URL.String())
        r.Headers.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
        r.Headers.Add("Accept-Language", "en-US,en;q=0.5")
        r.Headers.Add("Accept-Encoding", "gzip")
    })

    // Extract domains from Google search results
    collector.OnHTML("#center_col cite.apx8Vc", func(e *colly.HTMLElement) {
        domSelection := e.DOM
        link := domSelection.Contents().First().Text()
        domains = append(domains, link)
    })

    // Find and visit next Google search results page
    collector.OnHTML("#pnnext[href]", func(e *colly.HTMLElement) {
        link := e.Attr("href")

        err := e.Request.Visit(link)
        if err != nil {
            log.Println("Google scraping error: ", err)
        }
    })

    // Extract domains from Yandex search results
    collector.OnHTML("a.Link.Link_theme_outer.Path-Item.link.path__item.link.organic__greenurl", func(e *colly.HTMLElement) {
        link := e.ChildText("b")
        domains = append(domains, link)
    })

    collector.OnHTML("a.link.serp-url__link.serp-url__link_bold", func(e *colly.HTMLElement) {
        link := e.Text
        domains = append(domains, link)
    })

    // Find and visit next Yandex search results page
    collector.OnHTML(".Pager-Item_type_next", func(e *colly.HTMLElement) {
        link := e.Attr("href")

        err := e.Request.Visit(link)
        if err != nil {
             log.Println("Yandex scraping error: ", err)
        }
    })

    // For debug
    // collector.OnResponse(func(r *colly.Response) {
    //     log.Printf("%s\n", r.Body)
    // })

    // Set error handler
	collector.OnError(func(r *colly.Response, err error) {
        log.Println("Error:", err)
	})

    collector.Visit(googleQuery)
    collector.Visit(yandexQuery)

    return domains
}

// remove subdomain duplicates and schemes
func processFoundDomains(domains []string) []reflect.Value {
    set := make(map[string]void)

    for _, element := range domains {
        element = strings.ToLower(element)

        if strings.Contains(element, "http") {
            u, _ := url.Parse(element)

            set[u.Host] = void{}
        } else {
            set[element] = void{}
        }
    }

    result := reflect.ValueOf(set).MapKeys()
    return result
}
