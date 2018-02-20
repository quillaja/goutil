package pxu

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// An interal type to store 2 related numbers, a low and high bound.
// pixel.Vec could have been used instead, but this is creates a more
// self-documenting interface for the user.
type clamp struct {
	Low, High float64
}

// Camera is the basic interface for all camera types.
type Camera interface {
	Update(*pixelgl.Window)
	GetMatrix() pixel.Matrix
	Reset()
}
