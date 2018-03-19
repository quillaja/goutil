package pxu

// PixIndex is a convenience method to calculate the index of a pixel in
// a pixelgl.Canvas pixel array used in Pixels()/SetPixels(). The formula is
//     4*y*width + x*4
func PixIndex(x, y, width int) int {
	return y*4*width + x*4
}
