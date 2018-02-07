package rand

import (
	"fmt"
	"math/rand"
)

// Float64NM produced a random float64 in the range [low,high) from the
// standard source. Assumes source is already seeded. Panics if
// low >= high.
func Float64NM(low, high float64) float64 {
	if low >= high {
		panic(fmt.Errorf("Invalid params: %g not <= %g", low, high))
	}
	return rand.Float64()*(high-low) + low
}
