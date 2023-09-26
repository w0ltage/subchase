package main

import (
	"reflect"
	"testing"
)

func Test_processFoundDomains(t *testing.T) {

    t.Run("Bring to lowercase, remove duplicates and schemes", func(t *testing.T) {
        got := []string{"patent.google.com", "FONTS.gOOgle.com", "support.google.com", "https://support.google.com"}
        result := processFoundDomains(got)

        // Convert got slice to []reflect.Value
        var expected []reflect.Value
        for _, str := range got {
            expected = append(expected, reflect.ValueOf(str))
        }

        if reflect.DeepEqual(result, expected) {
            t.Errorf("Got %v, expected %v, result %v", got, expected, result)
        }
    })

    t.Run("Passed empty []string slice", func(t *testing.T) {
        got := []string{}
        result := processFoundDomains(got)

        if len(result) != 0 {
            t.Errorf("Got %v, expected length %v, result length of slice %v", got, 0, len(result))
        }
    })
}

