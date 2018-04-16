// Package main is a simple visualization of perlin noise.
package main

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	colorful "github.com/lucasb-eyer/go-colorful"

	"github.com/quillaja/goutil/misc"
	"github.com/quillaja/goutil/num"
	"github.com/quillaja/goutil/pxu"
	"github.com/quillaja/goutil/rand"
)

const title = "Perlin Noise Example"
const scale = 4

func run() {
	cfg := pixelgl.WindowConfig{
		Title:   title,
		Bounds:  pixel.R(0, 0, 800, 800),
		VSync:   true,
		Monitor: pixelgl.PrimaryMonitor(), // fullscreen
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	canvas := pixelgl.NewCanvas(pixel.R(0, 0, win.Bounds().W()/scale, win.Bounds().H()/scale))

	// perlin noise generation params
	xoff, yoff, zoff := 0.0, 0.0, 0.0
	xdelta, ydelta, zdelta := 0.02, 0.02, 0.01

	// example of using FillPermutation
	// rand.FillPermutation(1)

	// for fps output
	fps := time.NewTicker(time.Second)
	frames := 0

	// main loop
	avg := misc.NewAverager("perlin loop (ms)", 100)
	for !win.Closed() && !(win.JustPressed(pixelgl.KeyQ)) {
		winh, winw := int(canvas.Bounds().H()), int(canvas.Bounds().W())
		start := time.Now()
		pixels := canvas.Pixels()
		for y := 0; y < winh; y++ {
			for x := 0; x < winw; x++ {
				// h := num.Lerp(rand.Noise3Octaves(xoff, yoff, zoff, 2, 2, 0.5), -1, 1, 0, 360)
				// h := uint8(num.Lerp(math.Pow(1+rand.Noise3(xoff, yoff, zoff), 2), 0, 4, 0, 255))
				h := num.Lerp(rand.Noise3(xoff, yoff, zoff), -1, 1, 0, 360)
				r, g, b := colorful.Hsv(h, 1, 1).RGB255()
				i := pxu.PixIndex(x, y, winw)
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
		avg.Add(time.Since(start).Seconds() * 1000) //profiling

		// win.Clear(colornames.Blue)
		canvas.Draw(win, pixel.IM.Scaled(pixel.ZV, scale).Moved(win.Bounds().Center()))
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
