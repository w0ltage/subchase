package main

import (
	"flag"
	"fmt"
	"log"
	"os"

    "github.com/tokiakasu/subchase/version"
    "github.com/tokiakasu/subchase/unify"
    "github.com/tokiakasu/subchase/search"
    "github.com/tokiakasu/subchase/output"
)

func main() {
    const tool_version = "v0.3.0"

    var givenDomain string
    var quiet bool
    var jsonFlag bool

    flag.StringVar(&givenDomain, "d", "", "Specify the domain whose subdomains to look for (ex: -d google.com)")
    flag.BoolVar(&quiet, "silent", false, "Remove startup banner")
    flag.BoolVar(&jsonFlag, "json", false, "Output as JSON")
    flag.Parse()

    if !quiet {
        version.ShowBanner(tool_version)
    }

    if givenDomain == "" {
        log.Printf("No domain is passed to '-d' option\n\n")
        flag.Usage()
        os.Exit(1)
    }
    
    // Collect domains from search engines into []string
    rawDomains := search.ChaseDomains(givenDomain)

    if len(rawDomains) == 0 {
        log.Printf("No subdomains of %q was found", givenDomain)
    }

    // Bring elements in rawDomains slice to lower case 
    // + remove duplicates and schemes 
    domains := unify.Unify(rawDomains)

    if jsonFlag {
        data := output.SliceToJSON(domains)
        fmt.Println(string(data))
    } else {
        // Iterate through slice of unique domains
        for i := 0; i < len(domains); i++ {
            domain := domains[i]
            fmt.Println(domain.Interface())
        }
    }
}
