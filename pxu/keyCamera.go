package pxu

import (
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// KeyCamera provides a simple but relatively feature rich 2D camera for use with
// the Pixel graphics library.
//
// Position (panning) is controlled via the Position and PanSpeed members.
// By default, the Position is controled by the user via the Up, Down, Left, and
// Right keys. However, this can be customized by setting the UpButton,
// DownButton, LeftButton, and RightButton members to a pixelgl.Button value.
//
// Zoom is controlled by the Zoom and ZoomSpeed members. Zooming is based
// on multication, so a "zero" zoom is Zoom == 1. The ZoomSpeed is based on
// expotentiation where ZoomSpeed is the base which is raised to some "zoom
// level" power (ie ZoomSpeed**zoomlvl). This zoom level is usually the mouse's
// Y-scroll position, but can be customized by providing setting ZoomLevel to
// a function with the signature `func() float64`.
//
// Limits on panning and zooming are provided by XExtents, YExtents, and
// ZExtents which provide "Low" and "High" limits to which Position.X,
// Position.Y, and Position.Z are clamped to, respectively.
//
// The Reset() series of methods allow the camera position and/or zoom to be
// reset to their original values. If NewCamera() was used to instantiate, then
// those defaults are Pos(0,0) and Zoom(1). If NewCameraParams() was used, the
// initial values for position and zoom provided as parameters are saved within
// the camera and then used when Reset() is called.
//
// Example usage within the context of Pixel:
//
//     // ... prior initialization
//     cam := pixel.NewKeyCamera() // using defaults
//     // alternatively: NewCameraParams(width/2, height/2, 200, 1.1) for window center
//     loopStart := time.Now()
//     for !win.Closed() {
//         elapsed := time.Since(loopStart)
//         loopStart = time.Now()
//
//         // Reset() should be called before Update() and before
//         // setting the window's matrix (via SetMatrix())
//         if win.JustPressed(pixelgl.KeyHome) {
//             cam.Reset() // reset Position and Zoom to (0,0) and 1
//         }
//
//         cam.Update(win, elapsed.Seconds()) // camera updates position, zoom, etc
//         win.SetMatrix(cam.GetMatrix()) // provide transformation matrix to window
//         win.Clear()
//         //... redraw window, update state, etc
//
type KeyCamera struct {
	Position  pixel.Vec // X,Y location of the camera
	PanSpeed  float64   // pixels/sec speed of camera pan
	Zoom      float64   // current zoom factor of the camera
	ZoomSpeed float64   // speed at which the camera zooms, eg 2=2x zoom for each 'zoom level'

	XExtents clamp // min and max extents the camera can pan horizontally
	YExtents clamp // min and max extents the camera can pan vertically
	ZExtents clamp // min and max zoom factor

	UpButton, DownButton    pixelgl.Button // which buttons are to be used for up/down panning
	LeftButton, RightButton pixelgl.Button // which buttons are to be used for left/right panning
	ZoomLevel               func() float64 // func to get a zoom level, default to mouse y-scroll value

	// storage of values used internally
	prevWinBounds pixel.Rect
	origPosition  pixel.Vec
	origZoom      float64
	lastUpdate    time.Time
}

// NewCamera creates a new camera with sane defaults.
//
// Default values are:
//     Position: pixel.ZV // X=0, Y=0
//     PanSpeed: 200 // pixels/sec
//     Zoom: 1
//     ZoomSpeed: 1.1
//     XExtents.Low = -5000, XExtents.High = 5000
//     YExtents.Low = -5000, YExtents.High = 5000
//     ZExtents.Low = -50, ZExtents.High = 50
//
// Keyboard Up, Down, Left, and Right control panning, and the mouse wheel
// controls zoom.
func NewKeyCamera() *KeyCamera {
	return &KeyCamera{
		pixel.ZV,
		200,
		1,
		1.1,
		clamp{-5000, 5000},
		clamp{-5000, 5000},
		clamp{-50, 50},
		pixelgl.KeyUp, pixelgl.KeyDown,
		pixelgl.KeyLeft, pixelgl.KeyRight,
		nil,
		pixel.Rect{},
		pixel.ZV,
		1,
		time.Now()}
}

// NewCameraParams creates a new camera with the given parameters. "origPosition"
// and "origZoom" are stored and used for calls to Reset(). Other camera parameters
// are set according to the defaults (see NewCamera()) but can be changed.
func NewKeyCameraParams(origPosition pixel.Vec, origZoom, panSpeed, zoomSpeed float64) *KeyCamera {
	c := NewKeyCamera()
	c.Position = origPosition
	c.origPosition = origPosition
	c.Zoom = origZoom
	c.origZoom = origZoom
	c.PanSpeed = panSpeed
	c.ZoomSpeed = zoomSpeed
	return c
}

// Update recalculates the camera position and zoom, and is generally called
// each frame before setting the window's matrx. Update checks the keyboard for
// the status of the defined panning and zooming controls (see `Camera`).
// Position.X, Position.Y, and Zoom are clamped to XExtents, YExtents, and ZExtents
// respectively.
func (cam *KeyCamera) Update(win *pixelgl.Window) {
	// save window bounds (used in GetMatrix()) only when changed
	if cam.prevWinBounds != win.Bounds() {
		cam.prevWinBounds = win.Bounds()
	}
	timeElapsed := time.Since(cam.lastUpdate).Seconds()
	cam.lastUpdate = time.Now()
	// update pan
	if win.Pressed(cam.LeftButton) {
		cam.Position.X -= cam.PanSpeed * timeElapsed
	}
	if win.Pressed(cam.RightButton) {
		cam.Position.X += cam.PanSpeed * timeElapsed
	}
	if win.Pressed(cam.DownButton) {
		cam.Position.Y -= cam.PanSpeed * timeElapsed
	}
	if win.Pressed(cam.UpButton) {
		cam.Position.Y += cam.PanSpeed * timeElapsed
	}

	// update zoom based on either the user-defined "ZoomLevel()"
	// or the mouse wheel's Y position
	var zlvl float64
	if cam.ZoomLevel != nil {
		zlvl = cam.ZoomLevel()
	} else {
		zlvl = win.MouseScroll().Y
	}
	cam.Zoom *= math.Pow(cam.ZoomSpeed, zlvl)

	// clamp to extents
	cam.Position.X = pixel.Clamp(cam.Position.X, cam.XExtents.Low, cam.XExtents.High)
	cam.Position.Y = pixel.Clamp(cam.Position.Y, cam.YExtents.Low, cam.YExtents.High)
	cam.Zoom = pixel.Clamp(cam.Zoom, cam.ZExtents.Low, cam.ZExtents.High)

}

// GetMatrix gets the transformation matrix to apply the camera's settings to
// a window.
//
//     win.SetMatrix(cam.GetMatrix())
func (cam *KeyCamera) GetMatrix() pixel.Matrix {
	return pixel.IM.Scaled(cam.Position, cam.Zoom).
		Moved(cam.prevWinBounds.Center().Sub(cam.Position))
}

// Unproject will translate a point to its apparent position in the camera's
// view. This method is identical to:
//
//     m := cam.GetMatrix()
//     m.Unproject(point)
func (cam *KeyCamera) Unproject(point pixel.Vec) pixel.Vec {
	return cam.GetMatrix().Unproject(point)
}

// Reset restores the camera's Position and Zoom to its initial settings.
func (cam *KeyCamera) Reset() {
	cam.ResetPan()
	cam.ResetZoom()
}

// ResetPan restores the camera's Position to its initial settings.
func (cam *KeyCamera) ResetPan() {
	cam.ResetXPan()
	cam.ResetYPan()
}

// ResetXPan restores the camera's Position.X (horizontal pan) to its initial setting.
func (cam *KeyCamera) ResetXPan() { cam.Position.X = cam.origPosition.X }

// ResetYPan restores the camera's Position.Y (vertical pan) to its initial setting.
func (cam *KeyCamera) ResetYPan() { cam.Position.Y = cam.origPosition.Y }

// ResetZoom restores the camera's Zoom to its initial setting.
func (cam *KeyCamera) ResetZoom() { cam.Zoom = cam.origZoom }
