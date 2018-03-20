package rand

import (
	sr "math/rand"
	"testing"
)

func TestNoise3(t *testing.T) {
	// generate N noise values twice, using the same offsets and
	// check that they're the same
	N := 256 * 20
	XDELTA := 0.1

	FillPermutation(sr.NewSource(1)) // use deterministic 'random' number for testing

	A, B := make([]float64, N), make([]float64, N)
	for i, x := 0, 0.0; i < N; i, x = i+1, x+XDELTA {
		A[i] = Noise3(x, 0, 0)
	}
	for i, x := 0, 0.0; i < N; i, x = i+1, x+XDELTA {
		B[i] = Noise3(x, 0, 0)
	}

	for i := 0; i < N; i++ {
		if A[i] != B[i] {
			t.Errorf("A[%d] != B[%d]: %f != %f", i, i, A[i], B[i])
		}
	}
}

func TestNoise3_Range(t *testing.T) {
	// just checks for values near 1 or -1
	N := 100000

	FillPermutation(nil)

	for i := 0; i < N; i++ {
		a := Noise3(Float64NM(0, 256), Float64NM(0, 256), Float64NM(0, 256))
		if a > 0.9 || a < -0.9 {
			t.Log(a)
		}
	}
}

func BenchmarkNoise3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		offset := float64(i) / float64(b.N)
		Noise3(offset, offset, offset)
	}
}

func BenchmarkNoise3Octaves_4_2_05(b *testing.B) {
	for i := 0; i < b.N; i++ {
		offset := float64(i) / float64(b.N)
		Noise3Octaves(offset, offset, offset, 4, 2, 0.5)
	}
}
