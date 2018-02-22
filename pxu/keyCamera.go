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
// level" power (ie ZoomSpeed**zoomlvl). This zoom level controlled by the
// + and - keys by default, but can be customized by setting ZoomLevel to
// a function with the signature `func() float64`. For example, this library
// provides `GetZoomLevelFromMouseScroll()` which creates such a function to
// get the zoom level from the mouse's scroll wheel when provided with a
// pointer to a pixelgl.Window.
//
// Limits on panning and zooming are provided by XExtents, YExtents, and
// ZExtents which provide "Low" and "High" limits to which Position.X,
// Position.Y, and Position.Z are clamped to, respectively.
//
// The Reset() series of methods allow the camera position and/or zoom to be
// reset to their original values. If NewKeyCamera() was used to instantiate, then
// those defaults are the "origCenter" parameter and Zoom(1). If
// NewKeyCameraParams() was used, the initial values for position and zoom
// provided as parameters are saved within the camera and then used when
// Reset() is called.
type KeyCamera struct {
	Position  pixel.Vec // X,Y location in "world" of the camera
	PanSpeed  float64   // pixels/sec speed of camera pan
	Zoom      float64   // current zoom factor of the camera
	ZoomSpeed float64   // speed at which the camera zooms, eg 2=2x zoom for each 'zoom level'

	XExtents clamp // min and max extents the camera can pan horizontally
	YExtents clamp // min and max extents the camera can pan vertically
	ZExtents clamp // min and max zoom factor

	UpButton, DownButton        pixelgl.Button // which buttons are to be used for up/down panning
	LeftButton, RightButton     pixelgl.Button // which buttons are to be used for left/right panning
	ZoomInButton, ZoomOutButton pixelgl.Button // which buttons are to be used for zoom in/out
	ZoomLevel                   func() float64 // func to get a zoom level as alternative to buttons

	// storage of values used internally
	viewMatrix        pixel.Matrix
	worldZeroInWindow pixel.Vec
	origPosition      pixel.Vec
	origZoom          float64
	lastUpdate        time.Time
	corrected         bool // default of false
}

// NewKeyCamera creates a new camera with sane defaults.
//
// A recommended worldZeroInWindow is `win.Bounds().Center()` (center of window).
// Default values are:
//     PanSpeed: 200 // pixels/sec
//     Zoom: 1
//     ZoomSpeed: 1.1
//     XExtents.Low = -5000, XExtents.High = 5000
//     YExtents.Low = -5000, YExtents.High = 5000
//     ZExtents.Low = 0.1, ZExtents.High = 50
//
// Keyboard Up, Down, Left, and Right control panning, and the + and - buttons
// (aka = and -) control zooming.
func NewKeyCamera(worldZeroInWindow pixel.Vec) *KeyCamera {
	return &KeyCamera{
		Position:  pixel.ZV,
		PanSpeed:  200,
		Zoom:      1,
		ZoomSpeed: 1.1,
		XExtents:  clamp{-5000, 5000},
		YExtents:  clamp{-5000, 5000},
		ZExtents:  clamp{0.1, 50},
		UpButton:  pixelgl.KeyUp, DownButton: pixelgl.KeyDown,
		LeftButton: pixelgl.KeyLeft, RightButton: pixelgl.KeyRight,
		ZoomInButton: pixelgl.KeyEqual, ZoomOutButton: pixelgl.KeyMinus,
		ZoomLevel:         nil,
		viewMatrix:        pixel.IM,
		worldZeroInWindow: worldZeroInWindow,
		origPosition:      pixel.ZV,
		origZoom:          1,
		lastUpdate:        time.Now()}
}

// NewKeyCameraParams creates a new camera with the given parameters. "origPosition"
// and "origZoom" are stored and used for calls to Reset(). Other camera parameters
// are set according to the defaults (see NewKeyCamera()) but can be changed.
func NewKeyCameraParams(worldZeroInWindow, origPosition pixel.Vec, origZoom, panSpeed, zoomSpeed float64) *KeyCamera {
	c := NewKeyCamera(worldZeroInWindow)
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
// the status of the defined panning and zooming controls (see `KeyCamera`).
// Position.X, Position.Y, and Zoom are clamped to XExtents, YExtents, and ZExtents
// respectively.
func (cam *KeyCamera) Update(win *pixelgl.Window) {
	if !cam.corrected {
		// a bit hackish, but works for now
		cam.Position = win.Bounds().Center().Sub(cam.worldZeroInWindow)
		cam.origPosition = cam.Position
		cam.corrected = true
	}
	// need to know seconds elapsed since last call for correct cam movement
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
	// or the default behavior (keys)
	var zlvl float64
	if cam.ZoomLevel != nil {
		zlvl = cam.ZoomLevel()
	} else {
		// default behavior
		if win.Pressed(cam.ZoomInButton) {
			zlvl = timeElapsed * 30 // timeElapsed will be about 1/60 sec @ 60 fps
		}
		if win.Pressed(cam.ZoomOutButton) {
			zlvl = timeElapsed * -30
		}
	}
	cam.Zoom *= math.Pow(cam.ZoomSpeed, zlvl)

	// clamp to extents
	cam.Position.X = pixel.Clamp(cam.Position.X, cam.XExtents.Low, cam.XExtents.High)
	cam.Position.Y = pixel.Clamp(cam.Position.Y, cam.YExtents.Low, cam.YExtents.High)
	cam.Zoom = pixel.Clamp(cam.Zoom, cam.ZExtents.Low, cam.ZExtents.High)

	cam.viewMatrix = pixel.IM.Moved(win.Bounds().Center().Sub(cam.Position)).
		Scaled(win.Bounds().Center(), cam.Zoom)
}

// GetMatrix gets the transformation matrix to apply the camera's settings to
// a window.
//
//     win.SetMatrix(cam.GetMatrix())
func (cam *KeyCamera) GetMatrix() pixel.Matrix {
	return cam.viewMatrix
}

// Unproject will translate a point from its window position to its "world"
// position. This method is identical to:
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

// GetPosition returns the camera's position.
func (cam *KeyCamera) GetPosition() pixel.Vec {
	return cam.Position
}

// GetZoom returns the camera's zoom factor.
func (cam *KeyCamera) GetZoom() float64 {
	return cam.Zoom
}

// GetZoomLevelFromMouseScroll provides a convenient way to make a KeyCamera's
// "ZoomLvl" property, which is a func() float64.
//
//     cam.ZoomLvl = GetZoomLevelFromMouseScroll(win)
//     for !win.Closed() {
//         ...
func GetZoomLevelFromMouseScroll(win *pixelgl.Window) func() float64 {
	return func() float64 {
		return win.MouseScroll().Y
	}
}
