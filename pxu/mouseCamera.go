package pxu

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// MouseCamera is a camera type that uses only the mouse for camera control.
// Panning is done by a mouse drag action like that commonly found in photo
// editing software. Zooming is done such that the "zoom locus" is the point
// the mouse is currently over.
//
// For a basic example of how to use the camera, see the `Camera` interface,
// and read more in the documentation for `KeyCamera`.
type MouseCamera struct {
	Position  pixel.Vec // the 'world' location of the camera
	Zoom      float64   // the zoom factor of the camera
	ZoomSpeed float64   // how quickly the camera zooms. should be >1

	XExtents clamp // min and max extents the camera can pan horizontally
	YExtents clamp // min and max extents the camera can pan vertically
	ZExtents clamp // min and max extents for the zoom factor

	DragButton pixelgl.Button // the button which performs drag when pressed and held. default is left mouse buttom

	// used internally
	viewMatrix        pixel.Matrix
	worldZeroInWindow pixel.Vec
	prevMousePos      pixel.Vec
	origPos           pixel.Vec
	origZoom          float64
	corrected         bool
}

// NewMouseCamera creates a new camera with default values. A recommended setting
// for worldZeroInWindow is `win.Bounds().Center` (center of the window).
//
// Defaults:
//
//     Position:   pixel.ZV (0,0),
//     Zoom:       1,
//     ZoomSpeed:  1.1,
//     XExtents:   Low: -5000, High: 5000,
//     YExtents:   Low: -5000, High: 5000,
//     ZExtents:   Low: 0.1, High: 50,
//     DragButton: pixelgl.MouseButtonLeft
//
func NewMouseCamera(worldZeroInWindow pixel.Vec) *MouseCamera {
	return &MouseCamera{
		Position:          pixel.ZV,
		Zoom:              1,
		ZoomSpeed:         1.1,
		XExtents:          clamp{-5000, 5000},
		YExtents:          clamp{-5000, 5000},
		ZExtents:          clamp{0.1, 50},
		DragButton:        pixelgl.MouseButtonLeft,
		viewMatrix:        pixel.IM.Moved(worldZeroInWindow),
		worldZeroInWindow: worldZeroInWindow,
		prevMousePos:      pixel.ZV,
		origPos:           pixel.ZV,
		origZoom:          1}
}

// NewMouseCameraParams creates a camera with the parameters specified. Other
// paramaters can be still be modified. See `NewMouseCamera` for more info.
func NewMouseCameraParams(worldZeroInWindow, origPos pixel.Vec, origZoom, zoomSpeed float64) *MouseCamera {
	c := NewMouseCamera(worldZeroInWindow)
	c.origPos = origPos
	c.Zoom = origZoom
	c.ZoomSpeed = zoomSpeed
	c.origZoom = origZoom
	return c
}

// Update recalculates the camera's position and zoom, and is generally called
// each frame before setting the window's matrix.
func (c *MouseCamera) Update(win *pixelgl.Window) {
	if !c.corrected {
		// a bit hackish, but works for now
		c.Position = win.Bounds().Center().Sub(c.worldZeroInWindow)
		c.origPos = c.Position
		c.corrected = true
	}

	// translate the matrix only when mouse is dragged, and then translate
	// by the difference between the new and previous mouse positions.
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		c.prevMousePos = win.MousePosition()
	} else if win.Pressed(pixelgl.MouseButtonLeft) {
		newMousePos := win.MousePosition()
		delta := newMousePos.Sub(c.prevMousePos).Scaled(-1) // this delta is in the opposite direction we want to move "Position", so invert
		c.prevMousePos = newMousePos                        // todo? i wonder if there's a potential error here

		// check that new position is within extents.
		// if NOT then clamp
		newPos := c.Position.Add(delta.Scaled(1 / c.Zoom))
		newPos.X = pixel.Clamp(newPos.X, c.XExtents.Low, c.XExtents.High)
		newPos.Y = pixel.Clamp(newPos.Y, c.YExtents.Low, c.YExtents.High)

		// update matrix and position
		c.viewMatrix = c.viewMatrix.Moved(newPos.Sub(c.Position).Scaled(-c.Zoom))
		c.Position = newPos
	}

	// scale the matrix only when the mouse wheel is scrolled. MouseScroll()
	// returns the change since last window update, so zoomspeed^mouse.Y is
	// the change in zoom
	if win.MouseScroll().Y != 0 {
		delta := math.Pow(c.ZoomSpeed, win.MouseScroll().Y)
		// check that the new zoom is within extents.
		// if NOT then clamp
		newZoom := c.Zoom * delta
		newZoom = pixel.Clamp(newZoom, c.ZExtents.Low, c.ZExtents.High)

		// update matrix and zoom
		c.viewMatrix = c.viewMatrix.Scaled(win.MousePosition(), newZoom/c.Zoom)
		c.Zoom = newZoom

		// move position so that point under window center
		// accurately reflects "world" coordinate
		c.Position = c.viewMatrix.Unproject(win.Bounds().Center())
	}

}

// GetMatrix gets a transformation matrix to apply to the window.
//
//     win.SetMatrix(cam.GetMatrix())
func (c *MouseCamera) GetMatrix() pixel.Matrix {
	return c.viewMatrix
}

// ResetPosition resets the camera's position.
func (c *MouseCamera) ResetPosition() {
	c.viewMatrix = pixel.IM.Scaled(pixel.ZV, c.Zoom).Moved(c.worldZeroInWindow)
	c.Position = c.origPos
}

// ResetZoom resets the camera's zoom factor.
func (c *MouseCamera) ResetZoom() {
	c.viewMatrix = pixel.IM.Moved(c.worldZeroInWindow.Sub(c.Position))
	c.Zoom = c.origZoom
}

// Reset resets the camera's position and zoom factor.
func (c *MouseCamera) Reset() {
	// could do c.viewMatrix = pixel.IM.Moved(c.worldZeroInWindow) as well
	c.ResetZoom()
	c.ResetPosition()
}

// Unproject gets a "world" coordinate from a "screen" coordinate, and is
// indentical to:
//
//     cam.GetMatrix().Unproject(point)
func (c *MouseCamera) Unproject(point pixel.Vec) pixel.Vec {
	return c.viewMatrix.Unproject(point)
}

// GetPosition returns the camera's position.
func (c *MouseCamera) GetPosition() pixel.Vec {
	return c.Position
}

// GetZoom returns the camera's zoom factor.
func (c *MouseCamera) GetZoom() float64 {
	return c.Zoom
}
