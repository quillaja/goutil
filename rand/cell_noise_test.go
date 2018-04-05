package rand

import (
	"math/rand"
	"testing"

	"github.com/quillaja/goutil/data"
)

func TestCellNoise(t *testing.T) {

}

func BenchmarkCellNoiseSlow(b *testing.B) {
	FillPermutation(rand.NewSource(0))
	const m = 6
	noise := CellNoiseSlow(2, 5, data.EuclideanSq)
	for n := 0; n < b.N; n++ {
		noise(m*float64(n)/float64(b.N), m*float64(n)/float64(b.N))
	}
}

func BenchmarkCellNoise2D(b *testing.B) {
	FillPermutation(rand.NewSource(0))
	const m = 6
	noise := CellNoise2D(m, m, 2, 5, data.EuclideanSq)
	for n := 0; n < b.N; n++ {
		noise(m*float64(n)/float64(b.N), m*float64(n)/float64(b.N))
	}
}
