package engine

import (
	"fmt"
    "strings"

	"github.com/gocolly/colly"
	// "github.com/gocolly/colly/debug"
	// "github.com/gocolly/colly/extensions"
)

func YandexEngine(collector *colly.Collector, targetDomain string) []string {
    var foundDomains []string

    // Extract domains from Yandex search results
    collector.OnHTML("a.Link.Link_theme_outer", func(e *colly.HTMLElement) {
        link := e.Attr("href")
        // message := fmt.Sprintf("Found %q", link)

        foundDomains = append(foundDomains, link)
    })

    // Find and visit next Yandex search results page
    collector.OnHTML(".Pager-Item_type_next", func(e *colly.HTMLElement) {
        link := e.Attr("href")
        e.Request.Visit(link)
    })

    // Checks for YandexSmartCaptcha
    collector.OnHTML("#checkbox-captcha-form", func(e *colly.HTMLElement) {
        // loading_spinner.UpdateMessage("Yandex captured us with SmartCaptcha :(")
        fmt.Println("\nYandex captured us with SmartCaptcha")
    })

    if strings.ContainsAny(".", targetDomain) {
        yandexQuery := "https://yandex.com/search/?text=site:" + targetDomain + "&lr=100"

        collector.Visit(yandexQuery + "&lang=en")
        collector.Visit(yandexQuery + "&lang=ru")
    } else {
        // loading_spinner.UpdateMessage("Search by TLD detected. Switching to Google only.")
        fmt.Println("Search by TLD detected. Switching to Google only.")
    }

    return foundDomains
}

