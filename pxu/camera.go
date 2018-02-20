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
// Example usage within the context of Pixel:
//
//     // ... prior initialization
//     var cam Camera = NewKeyCamera(win.Bounds().Center()) // using defaults
//     // alternatively: NewKeyCameraParams(<somex>, <somey>, 200, 1.1) for another starting point
//     // or NewMouseCamera() etc.
//
//     for !win.Closed() {
//
//         // Reset() should be called before Update() and before
//         // setting the window's matrix (via SetMatrix())
//         if win.JustPressed(pixelgl.KeyHome) {
//             cam.Reset() // reset Position and Zoom to (0,0) and 1
//         }
//
//         cam.Update(win) // camera updates position, zoom, etc
//         win.SetMatrix(cam.GetMatrix()) // provide transformation matrix to window
//         win.Clear()
//         //... redraw window, update state, etc
//
type Camera interface {
	Update(*pixelgl.Window)
	GetMatrix() pixel.Matrix
	Reset()
	Unproject(pixel.Vec) pixel.Vec
}
