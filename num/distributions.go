package num

import (
	"math"
)

// simple factorial function
func factorial(n int) (f int) {
	f = 1
	for ; n > 1; n-- {
		f *= n
	}
	return
}

// Poisson returns the probability of an event happening k times when the
// expected frequency is lambda.
func Poisson(lambda, k int) float64 {
	return math.Pow(math.E, -float64(lambda)) * math.Pow(float64(lambda), float64(k)) / float64(factorial(k))
}
