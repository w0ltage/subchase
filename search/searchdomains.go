package search

import (
    "time"
    "fmt"
    "strings"
    "crypto/tls"
    "net/http"

	"github.com/gocolly/colly"
	// "github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/extensions"
    "github.com/leaanthony/spinner"

)

func ChaseDomains(givenDomain string) []string {
    var domains []string

    loading_spinner := spinner.New("Collecting domains from Google and Yandex")
    loading_spinner.Start()

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
            message := fmt.Sprintf("Google got tired of requests and started replying %q.\nRestart %q after a couple of minutes.", err, "subchase")
            loading_spinner.Error(message)
        } else {
            message := fmt.Sprintln("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
            loading_spinner.Error(message)
        }
	})

    // Extract domains from Google search results
    collector.OnHTML("#center_col cite", func(e *colly.HTMLElement) {
        domSelection := e.DOM
        link := domSelection.Contents().First().Text()

        if strings.Contains(link, givenDomain) {
            message := fmt.Sprintf("Found %q", link)
            loading_spinner.UpdateMessage(message)

            domains = append(domains, link)
        }
    })

    // Find and visit next Google search results page
    collector.OnHTML("#pnnext[href]", func(e *colly.HTMLElement) {
        link := e.Attr("href")
        e.Request.Visit(link)
    })

    // Extract domains from Yandex search results
    collector.OnHTML("a.Link.Link_theme_outer", func(e *colly.HTMLElement) {
        link := e.Attr("href")
        message := fmt.Sprintf("Found %q", link)
        loading_spinner.UpdateMessage(message)

        domains = append(domains, link)
    })

    // Find and visit next Yandex search results page
    collector.OnHTML(".Pager-Item_type_next", func(e *colly.HTMLElement) {
        link := e.Attr("href")
        e.Request.Visit(link)
    })

    // Checks for YandexSmartCaptcha
    collector.OnHTML("#checkbox-captcha-form", func(e *colly.HTMLElement) {
        loading_spinner.UpdateMessage("Yandex captured us with SmartCaptcha :(")
    })

    googleQuery := "https://www.google.com/search?q=site:*." + givenDomain
    collector.Visit(googleQuery)

    // Yandex sucks at search by TLD
    if strings.ContainsAny(".", givenDomain) {
        yandexQuery := "https://yandex.com/search/?text=site:" + givenDomain + "&lr=100"

        collector.Visit(yandexQuery + "&lang=en")
        collector.Visit(yandexQuery + "&lang=ru")
    } else {
        loading_spinner.UpdateMessage("Search by TLD detected. Switching to Google only.")
    }

    collector.Wait()

    loading_spinner.UpdateMessage("Finished")
    loading_spinner.Success()

    return domains
}
