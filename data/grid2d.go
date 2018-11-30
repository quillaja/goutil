package data

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/quillaja/goutil/num"
)

// hash "hashes" a point into a string representing its place in a Grid2D
// with buckets of `width` and `height`.
func hash(point mgl64.Vec2, width, height int) string {
	return fmt.Sprintf("%d,%d",
		int(num.RoundTo(point.X(), float64(width))),
		int(num.RoundTo(point.Y(), float64(height))))
}

type Grid2D struct {
}
