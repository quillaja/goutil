package main

import (
	"flag"
	"fmt"

	"github.com/quillaja/goutil/pxu"

	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var camType string

func main() {
	flag.StringVar(&camType, "cam", "", "Type of camera to test. Either 'keycamera' or 'mousecamera'.")
	flag.Parse()

	if camType == "" {
		panic("a camera type must be specified")
	}

	pixelgl.Run(run)

}

func run() {
	// create window
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:  "Camera test",
		Bounds: pixel.R(0, 0, 800, 600),
		VSync:  true,
	})
	if err != nil {
		panic(err)
	}

	// draw a reference image using "world" coordinates
	// it's a 1000x1000 "world coords" square with diagonals and a
	// 5px radius red circle in the center at (0,0).
	img := imdraw.New(nil)
	img.Color = colornames.Red
	img.Push(pixel.ZV)
	img.Circle(5, 0)
	pts := []pixel.Vec{
		{-500, 500},
		{-500, -500},
		{500, -500},
		{500, 500},
		{-500, -500},
		{500, -500},
		{-500, 500},
		{500, 500},
	}
	img.Color = colornames.Blue
	img.Push(pts...)
	img.Polygon(2)

	// create cams using Camera interface
	var cam pxu.Camera
	switch camType {
	case "keycamera":
		cam = pxu.NewKeyCamera(pixel.V(0, win.Bounds().H()/2)) //win.Bounds().Center())
		cam.(*pxu.KeyCamera).XExtents.High = 400
		cam.(*pxu.KeyCamera).YExtents.Low = -400
	case "mousecamera":
		cam = pxu.NewMouseCamera(pixel.V(win.Bounds().W()/3, win.Bounds().H()/3))
		cam.(*pxu.MouseCamera).XExtents.High = 400
		cam.(*pxu.MouseCamera).YExtents.Low = -400
	default:
		panic(fmt.Errorf("'%s' not a valid camera", camType))
	}

	for !win.Closed() {

		// camera parts
		if win.JustPressed(pixelgl.KeySpace) {
			cam.Reset()
		}
		cam.Update(win)
		win.SetMatrix(cam.GetMatrix())

		// draw reference image
		win.Clear(colornames.White)
		img.Draw(win)
		win.Update()

		win.SetTitle(fmt.Sprintf("cam.pos(%0.1f, %0.1f) cam.zoom(%0.1f)",
			cam.GetPosition().X, cam.GetPosition().Y, cam.GetZoom()))
		// example way to access concrete type when using various camera types.
		// switch c := cam.(type) {
		// case *pxu.KeyCamera:
		// 	p := c.Position
		// 	win.SetTitle(fmt.Sprintf("keycam.pos(%0.1f, %0.1f)", p.X, p.Y))
		// case *pxu.MouseCamera:
		// 	p := c.Position
		// 	win.SetTitle(fmt.Sprintf("mousecam.pos(%0.1f, %0.1f)", p.X, p.Y))
		// }
	}

}
