package utility

import (
)

func ConcatenateSlices(slices ...[]string) []string {
    var result []string

    for _, slice := range slices {
        result = append(result, slice...)
    }

    return result
}

