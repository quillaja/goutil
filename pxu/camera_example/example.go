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
	flag.StringVar(&camType, "cam", "", "Type of camera to test.")
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

	// create cams
	var cam pxu.Camera
	switch camType {
	case "keycamera":
		cam = pxu.NewKeyCamera(win.Bounds().Center())
	case "mousecamera":
		cam = pxu.NewMouseCamera(win.Bounds().Center())
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
	}

}
