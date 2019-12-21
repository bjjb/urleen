package base62

import (
	"math"
	"regexp"
)

// Base62 runes
const Base62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Pattern against which strings must match
var Pattern *regexp.Regexp

var base62 = map[rune]int{}

// Decode a Base62 string
func Decode(s string) uint64 {
	var x uint64
	max := len(s) - 1
	for i, c := range s {
		x = x + uint64(base62[c]*int(math.Pow(float64(62), float64(max-i))))
	}
	return x
}

// Encode i as a Base62 string
func Encode(i uint64) string {
	if i == 0 {
		return "0"
	}
	r := uint64(62)
	b := []byte{}
	for i > 0 {
		b = append([]byte{Base62[i%r]}, b...)
		i /= r
	}
	return string(b)
}

func init() {
	for i, r := range Base62 {
		base62[r] = i
	}
}
