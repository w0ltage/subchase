package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/gocolly/colly"
	// "github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/extensions"
)

// The `void` type is defined as an empty struct.
// It is used as the value type for the map (`set`) to create
// a set-like data structure where only unique elements are stored.
type void struct{}

type OutputJSON struct {
    Domains []string `json:"domains"`
}

const codename = "subchase"
const version = "v0.2.0"

func main() {
    var givenDomain string
    var quiet bool
    var jsonFlag bool

    flag.StringVar(&givenDomain, "d", "", "Specify the domain whose subdomains to look for (ex: -d google.com)")
    flag.BoolVar(&quiet, "silent", false, "Remove startup banner")
    flag.BoolVar(&jsonFlag, "json", false, "Output as JSON")
    flag.Parse()

    if !quiet {
        showBanner()
    }

    if givenDomain == "" {
        log.Printf("No domain is passed to '-d' option\n\n")
        flag.Usage()
        os.Exit(1)
    }
    
    // Collect domains from search engines into []string
    rawDomains := findDomains(givenDomain)

    if len(rawDomains) == 0 {
        log.Printf("No subdomains of %q was found", givenDomain)
    }

    // Bring elements in rawDomains slice to lower case 
    // + remove duplicates and schemes 
    domains := processFoundDomains(rawDomains)

    if jsonFlag {
        data := sliceToJSON(domains)
        fmt.Println(string(data))
    } else {
        // Iterate through slice of unique domains
        for i := 0; i < len(domains); i++ {
            domain := domains[i]
            fmt.Println(domain.Interface())
        }
    }
}

func findDomains(givenDomain string) []string {
    var domains []string


    // Instantiate default collector
    collector := colly.NewCollector(
        colly.Async(true),
        // colly.CacheDir("./sites_cache"),
        // colly.Debugger(&debug.LogDebugger{}),
        colly.DetectCharset(),
        colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/98.0"),
        )

    collector.Limit(&colly.LimitRule{
        Parallelism: 2,
        RandomDelay: 5 * time.Second,
    })

    // Setting the max TLS version to 1.2
    // Without specifying the maximum version of TLS 1.2, 
    // requests get a response "403 Forbidden".
    collector.WithTransport(&http.Transport{
        TLSClientConfig: &tls.Config{
            MaxVersion: tls.VersionTLS12,
        },
    })

    // Referer sets valid Referer HTTP header to requests
    extensions.Referer(collector)
    // extensions.RandomUserAgent(collector)

    // Add headers to requests to imitate Firefox
    collector.OnRequest(func(r *colly.Request) {
        r.Headers.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
        r.Headers.Add("Accept-Language", "en-US,en;q=0.5")
        r.Headers.Add("Accept-Encoding", "gzip")
        r.Headers.Add("DNT", "1")
        r.Headers.Add("Connection", "keep-alive")
        r.Headers.Add("Upgrade-Insecure-Requests", "1")
        r.Headers.Add("Sec-Fetch-Dest", "document")
        r.Headers.Add("Sec-Fetch-Mode", "navigate")
        r.Headers.Add("Sec-Fetch-Site", "same-origin")
        r.Headers.Add("Sec-Fetch-User", "?1")
    })

    // Set error handler
	collector.OnError(func(r *colly.Response, err error) {
        if r.StatusCode == http.StatusTooManyRequests {
            log.Printf("\nGoogle got tired of requests and started replying %q. Restart %q after a couple of minutes.", err, codename)
        } else {
            log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
        }
	})

    // Extract domains from Google search results
    collector.OnHTML("#center_col cite", func(e *colly.HTMLElement) {
        domSelection := e.DOM
        link := domSelection.Contents().First().Text()

        if strings.Contains(link, givenDomain) {
            domains = append(domains, link)
        }
    })

    // Find and visit next Google search results page
    collector.OnHTML("#pnnext[href]", func(e *colly.HTMLElement) {
        link := e.Attr("href")
        e.Request.Visit(link)
    })

    // Extract domains from Yandex search results
    collector.OnHTML("a.Link.Link_theme_outer.Path-Item.link.path__item.link.organic__greenurl", func(e *colly.HTMLElement) {
        link := e.ChildText("b")
        domains = append(domains, link)
    })

    // Find and visit next Yandex search results page
    collector.OnHTML(".Pager-Item_type_next", func(e *colly.HTMLElement) {
        link := e.Attr("href")
        e.Request.Visit(link)
    })

    // Checks for YandexSmartCaptcha
    // collector.OnHTML("#checkbox-captcha-form", func(e *colly.HTMLElement) {
    //     log.Println("Captcha found! Aborting operation.")
    // })


    googleQuery := "https://www.google.com/search?q=site:*." + givenDomain
    collector.Visit(googleQuery)

    // Yandex sucks at search by TLD
    if strings.Contains(".", givenDomain) {
        yandexQuery := "https://yandex.com/search/?text=site:" + givenDomain + "&lr=100"

        collector.Visit(yandexQuery + "&lang=en")
        collector.Visit(yandexQuery + "&lang=ru")
    }

    collector.Wait()

    return domains
}

// Bring domains to lowercase
// and remove duplicates + schemes
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

func showBanner() {
    fmt.Printf(
`               __         __                  
   _______  __/ /_  _____/ /_  ____ _________ 
  / ___/ / / / __ \/ ___/ __ \/ __ %c/ ___/ _ \
 (__  ) /_/ / /_/ / /__/ / / / /_/ (__  )  __/
/____/\__,_/_.___/\___/_/ /_/\__,_/____/\___/  %v

`, '`', version)
}

func sliceToJSON(values []reflect.Value) []byte {
    var strings []string

    for _, value := range values {
        strings = append(strings, value.String())
    }

    outputJSON := OutputJSON{
        Domains: strings,
    }

    jsonData, err := json.Marshal(outputJSON)
    if err != nil {
        log.Println("Error marhaling to JSON:", err)
        os.Exit(1)
    }

    return jsonData
}
