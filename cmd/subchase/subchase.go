package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
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

    // googleQuery := "https://www.google.com/search?q=site:" + givenDomain
    yandexQuery := "https://yandex.com/search/?text=site:" + givenDomain + "&lr=100&p=0"

    googleCollector := colly.NewCollector(
        // colly.CacheDir("./sites_cache"),
        colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/114.0"),
        colly.Async(true),
        colly.Debugger(&debug.LogDebugger{}),
        )
    
    // Referer sets valid Referer HTTP header to requests
    extensions.Referer(googleCollector)

    // Extract domains from Google search results
    googleCollector.OnHTML("#center_col cite.apx8Vc", func(e *colly.HTMLElement) {
        domSelection := e.DOM
        link := domSelection.Contents().First().Text()
        domains = append(domains, link)
    })

    // Find and visit next Google search results page
    googleCollector.OnHTML("#pnnext[href]", func(e *colly.HTMLElement) {
        link := e.Attr("href")

        err := e.Request.Visit(link)
        if err != nil {
            log.Println("Google scraping error: ", err)
        }
    })

    // Set error handler
	googleCollector.OnError(func(r *colly.Response, err error) {
        log.Println("Error:", err)
	})

    // Instantiate Yandex collector
    yandexCollector := googleCollector.Clone()
    extensions.Referer(yandexCollector)
    yandexCollector.Limit(&colly.LimitRule{
        DomainGlob: "*yandex*",
        Parallelism: 4,
        Delay: 4 * time.Second,
    })

    // Disable TLS security check for a client
    yandexCollector.WithTransport(&http.Transport{
        TLSClientConfig:&tls.Config{InsecureSkipVerify: true},
    })

    // For some unknown reason, requests get "Forbidden"
    // without using a localhost proxy (mitmproxy)
    err := yandexCollector.SetProxy("localhost:80")
    if err != nil {
        log.Fatalln("Error happened with proxy:", err)
    }

    // Add headers to requests to Yandex
    yandexCollector.OnRequest(func(r *colly.Request) {
        log.Println("visiting", r.URL.String())
        r.Headers.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
        r.Headers.Add("Accept-Language", "en-US,en;q=0.5")
        r.Headers.Add("Accept-Encoding", "gzip")
    })

    // Extract domains from Yandex search results
    yandexCollector.OnHTML("a.Link.Link_theme_outer.Path-Item.link.path__item.link.organic__greenurl", func(e *colly.HTMLElement) {
        link := e.ChildText("b")
        domains = append(domains, link)
    })

    yandexCollector.OnHTML("a.link.serp-url__link.serp-url__link_bold", func(e *colly.HTMLElement) {
        link := e.Text
        domains = append(domains, link)
    })

    // Find and visit next Yandex search results page
    yandexCollector.OnHTML(".Pager-Item_type_next", func(e *colly.HTMLElement) {
        link := e.Attr("href")

        err := e.Request.Visit(link)
        if err != nil {
             log.Println("Yandex scraping error: ", err)
        }
    })

    yandexCollector.OnHTML("#checkbox-captcha-form", func(e *colly.HTMLElement) {
        log.Println("Captcha found! Aborting parsing.")
    })

    // For debug
    // collector.OnResponse(func(r *colly.Response) {
    //     log.Printf("%s\n", r.Body)
    // })

	yandexCollector.OnError(func(r *colly.Response, err error) {
        log.Println("Error:", err)
	})

    // googleCollector.Visit(googleQuery)
    // googleCollector.Wait()

    yandexCollector.Visit(yandexQuery + "&lang=en")
    // yandexCollector.Visit(yandexQuery + "&lang=ru")
    yandexCollector.Wait()

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
