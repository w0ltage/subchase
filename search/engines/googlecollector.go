package engine

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
	"github.com/leaanthony/spinner"
	// "github.com/gocolly/colly/debug"
	// "github.com/gocolly/colly/extensions"
)

func GoogleEngine(collector *colly.Collector, targetDomain string, loadingSpinner *spinner.Spinner) []string {
    var foundDomains []string

    // Set error handler if "Too Many Requests" arise
	collector.OnError(func(r *colly.Response, err error) {
        if r.StatusCode == http.StatusTooManyRequests {
            message := fmt.Sprintf("Google got tired of requests and started replying %q.\nRestart %q after a couple of minutes.", err, "subchase")
            loadingSpinner.Error(message)
        } else {
            message := fmt.Sprintln("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
            loadingSpinner.Error(message)
        }
	})

    // Extract domains from Google search results
    collector.OnHTML("#center_col cite", func(e *colly.HTMLElement) {
        domSelection := e.DOM
        link := domSelection.Contents().First().Text()

        if strings.Contains(link, targetDomain) {
            message := fmt.Sprintf("Found %q", link)
            loadingSpinner.UpdateMessage(message)

            foundDomains = append(foundDomains, link)
        }
    })

    // Find and visit next Google search results page
    collector.OnHTML("#pnnext[href]", func(e *colly.HTMLElement) {
        link := e.Attr("href")
        e.Request.Visit(link)
    })

    googleQuery := "https://www.google.com/search?q=site:*." + targetDomain
    collector.Visit(googleQuery)

    return foundDomains
}
