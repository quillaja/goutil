package rand

import (
	"fmt"
	"math/rand"
)

// IntNM produced a random int in the range [low,high) from the
// standard source. Assumes source is already seeded. Panics if
// low >= high.
func IntNM(low, high int) int {
	if low >= high {
		panic(fmt.Errorf("Invalid params: %d not <= %d", low, high))
	}
	return int(rand.Float64()*float64(high-low)) + low
}
