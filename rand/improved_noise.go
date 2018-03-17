package rand

import (
	"math"
	"math/rand"
	"time"

	"github.com/quillaja/goutil/num"
)

// permuation array
var p [512]int

// initialize when the package is used
func init() {
	FillPermutation(nil)
}

// FillPermutation changes the permutation table used for noise generation.
// If a source is provided, it uses that. Otherwise, it uses a source seeded
// with the time.Now().UnixNano().
func FillPermutation(source rand.Source) {
	var r *rand.Rand
	if source == nil {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	} else {
		r = rand.New(source)
	}
	// get a permutation of [0,255]
	copy(p[:256], r.Perm(256))
	// fill 2nd half of array with duplicates
	copy(p[256:], p[:256])
}

// returns the dot product of vec3(x,y,z) (x,y,z in [0,1])
func grad(hash int, x, y, z float64) float64 {
	switch hash & 0xF {
	case 0x0:
		return x + y
	case 0x1:
		return -x + y
	case 0x2:
		return x - y
	case 0x3:
		return -x - y
	case 0x4:
		return x + z
	case 0x5:
		return -x + z
	case 0x6:
		return x - z
	case 0x7:
		return -x - z
	case 0x8:
		return y + z
	case 0x9:
		return -y + z
	case 0xA:
		return y - z
	case 0xB:
		return -y - z
	case 0xC:
		return y + x
	case 0xD:
		return -y + z
	case 0xE:
		return y - x
	case 0xF:
		return -y - z
	default:
		return 0 // never happens
	}
}

// Noise3 returns 3d perlin noise based on Ken Perlin's 2002
// "improved noise" algorithm.
//
// See: http://mrl.nyu.edu/~perlin/noise/
// Paper: http://mrl.nyu.edu/~perlin/paper445.pdf
func Noise3(x, y, z float64) float64 {
	// find unit cube that contains point
	xCube, yCube, zCube := 255&int(x), 255&int(y), 255&int(z)

	// x,y,z in [0,1] as a porportional location inside that cube
	x, y, z = x-math.Floor(x), y-math.Floor(y), z-math.Floor(z)

	// fade curves for x,y,z
	u, v, w := num.SmootherStep(x), num.SmootherStep(y), num.SmootherStep(z)

	// hash coordinates of the 8 cube corners
	A := p[xCube] + yCube
	AA := p[A] + zCube
	AB := p[A+1] + zCube
	B := p[xCube+1] + yCube
	BA := p[B] + zCube
	BB := p[B+1] + zCube

	// get gradients from point to corners of unit cube, using 'dot product',
	// then blend them together
	return num.UnitLerp(w, num.UnitLerp(v, num.UnitLerp(u, grad(p[AA], x, y, z),
		grad(p[BA], x-1, y, z)),
		num.UnitLerp(u, grad(p[AB], x, y-1, z),
			grad(p[BB], x-1, y-1, z))),
		num.UnitLerp(v, num.UnitLerp(u, grad(p[AA+1], x, y, z-1),
			grad(p[BA+1], x-1, y, z-1)),
			num.UnitLerp(u, grad(p[AB+1], x, y-1, z-1),
				grad(p[BB+1], x-1, y-1, z-1))))

	// above should be formatted as below, but go's formatter won't let it be
	// num.UnitLerp(w, num.UnitLerp(v, num.UnitLerp(u, grad(p[AA],   x,   y,   z  ),
	//                                                 grad(p[BA],   x-1, y,   z  )),
	// 	                               num.UnitLerp(u, grad(p[AB],   x,   y-1, z  ),
	// 		                                           grad(p[BB],   x-1, y-1, z  ))),
	// 	               num.UnitLerp(v, num.UnitLerp(u, grad(p[AA+1], x,   y,   z-1),
	// 		                                           grad(p[BA+1], x-1, y,   z-1)),
	// 		                           num.UnitLerp(u, grad(p[AB+1], x,   y-1, z-1),
	// 			                                       grad(p[BB+1], x-1, y-1, z-1))))
}
