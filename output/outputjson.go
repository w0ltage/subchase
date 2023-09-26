package output

import (
    "os"
    "log"
    "reflect"
    "encoding/json"
)

type OutputJSON struct {
    Domains []string `json:"domains"`
}

func SliceToJSON(values []reflect.Value) []byte {
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
