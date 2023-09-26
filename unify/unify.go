package unify

import (
    "reflect"
    "strings"
    "net/url"
)

// The `void` type is defined as an empty struct.
// It is used as the value type for the map (`set`) to create
// a set-like data structure where only unique elements are stored.
type void struct{}

func toLower(domainsArray []string) []string {
    var result []string

    for _, element := range domainsArray {
        element = strings.ToLower(element)
        result = append(result, element)
    }

    return result
}

func removeHTTP(domainsArray []string) []string {
    var result []string

    for _, element := range domainsArray {
        if strings.Contains(element, "http") {
            u, _ := url.Parse(element)   
            result = append(result, u.Host)
        } else {
            result = append(result, element)
        }
    }

    return result
}

func deduplicate(domainsArray []string) []reflect.Value {
    set := make(map[string]void)

    for _, element := range domainsArray {
        set[element] = void{}
    }

    result := reflect.ValueOf(set).MapKeys()
    return result
}

// Bring domains to lowercase
// and remove duplicates + schemes
func Unify(domainsArray []string) []reflect.Value {
    domainsToLowercase := toLower(domainsArray)
    domainsRemoveHTTP := removeHTTP(domainsToLowercase)
    domainsDeduplicated := deduplicate(domainsRemoveHTTP)

    return domainsDeduplicated
}

// func Unify(domains []string) []reflect.Value {
//     set := make(map[string]void)

//     for _, element := range domains {
//         element = strings.ToLower(element)

//         if strings.Contains(element, "http") {
//             u, _ := url.Parse(element)
//             set[u.Host] = void{}

//         } else {
//             set[element] = void{}
//         }
//     }

//     result := reflect.ValueOf(set).MapKeys()
//     return result
// }
