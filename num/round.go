package num

import (
	"math"
)

// RoundTo rounds x to the nearest multiple of n.
func RoundTo(x, n float64) float64 {
	return n * math.Round(x/n)
}

// RoundToInt rounds x to the nearest multiple of n.
func RoundToInt(x, n int) int {
	return n * int(math.Round(float64(x)/float64(n)))
}
