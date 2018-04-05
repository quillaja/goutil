package rand

import (
	"math/rand"
	"testing"

	"github.com/quillaja/goutil/data"
)

func TestCellNoise(t *testing.T) {

}

func BenchmarkCellNoise(b *testing.B) {
	FillPermutation(rand.NewSource(0))
	noise := CellNoise(2, 5, data.EuclideanSq)
	for n := 0; n < b.N; n++ {
		noise(4*float64(n)/float64(b.N), 4*float64(n)/float64(b.N))
	}
}
