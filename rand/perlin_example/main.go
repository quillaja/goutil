// Package main is a simple visualization of perlin noise.
package main

import (
	"fmt"
	"time"

	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/quillaja/goutil/pxu"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/quillaja/goutil/num"
	"github.com/quillaja/goutil/rand"
)

const title = "Perlin Noise Example"

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  title,
		Bounds: pixel.R(0, 0, 800, 600),
		VSync:  true,
		// Monitor: pixelgl.PrimaryMonitor(), // fullscreen
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	canvas := pixelgl.NewCanvas(win.Bounds())

	// perlin noise generation params
	xoff, yoff, zoff := 0.0, 0.0, 0.0
	xdelta, ydelta, zdelta := 0.005, 0.005, 0.015

	// example of using FillPermutation
	// rand.FillPermutation(sr.NewSource(1))

	// run perlin junk in a goroutine to keep main thread uninteruppted
	// this draws perlin noise as colors
	go func() {
		for !win.Closed() && !(win.JustPressed(pixelgl.KeyQ)) {
			pixels := canvas.Pixels()
			for y := 0; y < int(canvas.Bounds().H()); y++ {
				for x := 0; x < int(canvas.Bounds().W()); x++ {
					// h := num.Lerp(rand.Noise3Octaves(xoff, yoff, zoff, 4, 2, 0.5), -1, 1, 0, 360)
					h := num.Lerp(rand.Noise3(xoff, yoff, zoff), -1, 1, 0, 360)
					r, g, b := colorful.Hsl(h, 1, 0.5).RGB255()
					i := pxu.PixIndex(x, y, int(canvas.Bounds().W()))
					pixels[i+0] = r   // r
					pixels[i+1] = g   // g
					pixels[i+2] = b   // b
					pixels[i+3] = 255 // a

					xoff += xdelta
				}
				yoff += ydelta
				xoff = 0.0
			}
			yoff = 0.0
			zoff += zdelta
			canvas.SetPixels(pixels)
		}
	}()

	// for fps output
	fps := time.NewTicker(time.Second)
	frames := 0

	// main loop
	for !win.Closed() && !(win.JustPressed(pixelgl.KeyQ)) {

		// win.Clear(colornames.Blue)
		canvas.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		win.Update()

		select {
		case <-fps.C:
			win.SetTitle(fmt.Sprintf("%s - %d fps", title, frames))
			// fmt.Println(frames)
			frames = 0
		default:
			frames++
		}
	}
}

func main() {
	pixelgl.Run(run)
}
