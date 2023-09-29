package search

import (
	"crypto/tls"
	// "fmt"
	"net/http"
	// "strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/tokiakasu/subchase/search/engines"
	// "github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/extensions"
	"github.com/leaanthony/spinner"
	// "github.com/charmbracelet/bubbles/spinner"
)

func ChaseDomains(targetDomain string) []string {
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

    // yandexDomains := engine.YandexEngine(collector, targetDomain)
    googleDomains := engine.GoogleEngine(collector, targetDomain)

    // domains = append(domains, yandexDomains...)
    domains = append(domains, googleDomains...)

    collector.Wait()

    loading_spinner.UpdateMessage("Finished")
    loading_spinner.Success()

    return domains
}
