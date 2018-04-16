package rand

import (
	"testing"

	"github.com/quillaja/goutil/data"
)

// These are probably the crappiest benchmarks ever.

func BenchmarkCellNoiseSlow(b *testing.B) {
	FillPermutation(0)
	const m = 10
	noise := CellNoiseSlow(2, 5, data.EuclideanSq)
	for n := 0; n < b.N; n++ {
		offset := m * float64(n) / float64(b.N)
		noise(offset, offset)
	}
}

func BenchmarkCellNoise2D(b *testing.B) {
	const m = 10
	noise := NewCellNoise2D(0, 2, 5, data.EuclideanSq)
	for n := 0; n < b.N; n++ {
		offset := m * float64(n) / float64(b.N)
		noise.Noise(offset, offset)
	}
}

func BenchmarkCellNoise3D(b *testing.B) {
	const m = 10
	noise := NewCellNoise3D(0, 2, 5, data.EuclideanSq)
	for n := 0; n < b.N; n++ {
		offset := m * float64(n) / float64(b.N)
		noise.Noise(offset, offset, offset)
	}
}
