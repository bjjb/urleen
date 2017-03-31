package base62

import (
	"math"
)

const (
	// Alphabet contains all characters for Base62
	Alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// Radix is the length of the Alphabet
	Radix = len(Alphabet)
)

var (
	quick = map[rune]int{} // maps chars to their position, for efficiency
)

func init() {
	for i, r := range Alphabet {
		quick[r] = i
	}
}

// Decode turns a base62 string into a uint64
func Decode(s string) uint64 {
	var x uint64
	max := len(s) - 1
	for i, c := range s {
		x = x + uint64(quick[c]*int(math.Pow(float64(Radix), float64(max-i))))
	}
	return x
}

// Encode turns a uint64 into a base62 string
func Encode(i uint64) string {
	if i == 0 {
		return "0"
	}
	r := uint64(Radix)
	b := []byte{}
	for i > 0 {
		b = append([]byte{Alphabet[i%r]}, b...)
		i /= r
	}
	return string(b)
}
