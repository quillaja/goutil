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
// "if I zoom with the mouse at the top right, then the same "world coordinates"
// would be under the mouse after the zoom."
type MouseCamera struct {
	Pos       pixel.Vec
	Zoom      float64
	ZoomSpeed float64

	XExtents clamp
	YExtents clamp
	ZExtents clamp

	DragButton pixelgl.Button

	viewMatrix   pixel.Matrix
	prevMousePos pixel.Vec
	origCenter   pixel.Vec
}

func NewMouseCamera(initialCenter pixel.Vec) *MouseCamera {
	return &MouseCamera{
		Pos:          initialCenter,
		Zoom:         1,
		ZoomSpeed:    1.1,
		XExtents:     clamp{-5000, 5000},
		YExtents:     clamp{-5000, 5000},
		ZExtents:     clamp{-50, 50},
		DragButton:   pixelgl.MouseButtonLeft,
		viewMatrix:   pixel.IM.Moved(initialCenter),
		prevMousePos: pixel.ZV,
		origCenter:   initialCenter}
}

func (c *MouseCamera) Update(win *pixelgl.Window) {
	// translate the matrix only when mouse is dragged, and then translate
	// by the difference between the new and previous mouse positions.
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		c.prevMousePos = win.MousePosition()
	} else if win.Pressed(pixelgl.MouseButtonLeft) {
		newMousePos := win.MousePosition()
		delta := newMousePos.Sub(c.prevMousePos)

		// todo: check that Pos + delta is within extents.
		// if NOT, then alter delta to land in extents, update Pos etc.
		c.Pos = c.Pos.Add(delta)
		c.prevMousePos = newMousePos

		c.viewMatrix = c.viewMatrix.Moved(delta)
	}

	// scale the matrix only when the mouse wheel is scrolled. MouseScroll()
	// returns the change since last window update, so zoomspeed^mouse.Y is
	// the change in zoom
	if win.MouseScroll().Y != 0 {
		delta := math.Pow(c.ZoomSpeed, win.MouseScroll().Y)
		//TODO: check that Zoom*delta is within extents.
		// if NOT, then alter delta to land in extents, update Zoom, etc.
		c.Zoom *= delta
		c.viewMatrix = c.viewMatrix.Scaled(win.MousePosition(), delta)
	}
}

func (c *MouseCamera) GetMatrix() pixel.Matrix {
	return c.viewMatrix
}

func (c *MouseCamera) ResetPosition() {
	c.viewMatrix = c.viewMatrix.Moved(c.origCenter.Sub(c.Pos))
	c.Pos = c.origCenter
}

func (c *MouseCamera) ResetZoom() {
	c.viewMatrix = pixel.IM.Moved(c.Pos)
	c.Zoom = 1
}

func (c *MouseCamera) Reset() {
	c.ResetPosition()
	c.ResetZoom()
}

/*
@elliotmr's original code.

cameraOrigin := pixel.ZV.Add(win.Bounds().Center())
scale := 1.0
dragOrigin := pixel.V(0, 0)
second := time.Tick(time.Second)
viewMatrix := pixel.IM
frames := 0
for !win.Closed() {
	if win.MouseScroll().Y != 0 {
		factor := math.Pow(1.2, win.MouseScroll().Y)
		zoomDeltaStart := viewMatrix.Unproject(win.MousePosition())
		scale *= factor
		cameraOrigin = zoomDeltaStart.Add(win.Bounds().Center().Sub(win.MousePosition().Scaled(1 / scale)))
	}
	if win.JustPressed(pixelgl.MouseButton1) {
		dragOrigin = win.MousePosition().Scaled(1 / scale)
	} else if win.Pressed(pixelgl.MouseButton1) {
		newOrigin := win.MousePosition().Scaled(1 / scale)
		cameraOrigin = cameraOrigin.Sub(newOrigin.Sub(dragOrigin))
		dragOrigin = newOrigin
	}
	viewMatrix = pixel.IM.Moved(win.Bounds().Center().Sub(cameraOrigin)).Scaled(pixel.ZV, scale)
*/
