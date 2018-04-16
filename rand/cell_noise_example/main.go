// Package main is a simple demo of cell noise.
// 2d by default, but uncomment the "z"
// lines of code and use NewCellNoise3D() instead to produce an
// animation of 3d noise.
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"time"

	"github.com/quillaja/goutil/data"
	qr "github.com/quillaja/goutil/rand"
)

func main() {
	const (
		width  = 800 // output image width
		height = 800 // output image height
		//depth     = 16  //* 60 // 60s at 16 fps
		pxPerCell = 150 // px/cell
	)

	// create 3 different noise generators
	seed := time.Now().UnixNano()
	red := qr.NewCellNoise2D(seed, 2, 5, data.Euclidean)
	green := qr.NewCellNoise2D(seed+1, 2, 5, data.Euclidean)
	blue := qr.NewCellNoise2D(seed+2, 2, 5, data.Euclidean)

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	start := time.Now()
	fmt.Println("Processing...")
	//z := 0
	// use z for-loop to create 'animated' noise
	// for z := 0; z < depth; z++ {
	// zs := float64(z) / pxPerCell
	for y := 0; y < height; y++ {
		ys := float64(y) / pxPerCell
		for x := 0; x < width; x++ {
			xs := float64(x) / pxPerCell
			img.Set(x, y, color.RGBA{
				255 - uint8(255*red.Noise(xs, ys)),   //, zs)),
				255 - uint8(255*green.Noise(xs, ys)), //, zs)),
				255 - uint8(255*blue.Noise(xs, ys)),  //, zs)),
				255})
		}
	}
	file, _ := os.Create("output.jpg") //fmt.Sprintf("output/%04d.jpg", z)) // for animation
	jpeg.Encode(file, img, &jpeg.Options{Quality: 95})
	file.Close()
	// fmt.Printf("\u001b[100D%0.1f%% complete.", (float32(z)+1)/depth*100) // for animation
	// }
	fmt.Println("image gen took (s)", time.Since(start).Seconds())
}
