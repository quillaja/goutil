package rand

import (
	"fmt"
	"math"

	"github.com/quillaja/goutil/data"

	"github.com/quillaja/goutil/num"
)

func hash(x, y int) int {
	//a >= b ? a * a + a + b : a + b * b;  where a, b >= 0
	a, b := x%(1<<16), y%(1<<16)
	if a >= b {
		return a*a + a + b
	}
	return a + b*b
}

// 2D point for cell noise
type point [2]float64

func (p *point) Location() []float64 {
	return p[:]
}

// CellNoiseSlow creates a function which gives 2D cell noise. Lambda and max determine
// the number of feature points in a given unit cube, which lambda is the average
// number per cell and max is the maximum number per cell. DistanceMetric dist
// provides the notion of 'distance'.
func CellNoiseSlow(lambda, max int, dist data.DistanceMetric) func(x, y float64) float64 {
	cdf := make([]float64, max+1, max+1)
	for k, t := 0, 0.0; k <= max; k++ {
		p := num.Poisson(lambda, k)
		t += p
		cdf[k] = t
	}
	// fmt.Println(cdf)

	return func(x, y float64) float64 {
		points := make([]data.Interface, 0, max*9)
		// pchan := make(chan data.Interface, max*9)
		// wg := sync.WaitGroup{}
		// wg.Add(9)

		// for the cell and 8 cells in its neighborhood
		for r := 1; r >= -1; r-- {
			for c := -1; c <= 1; c++ {
				// go func(r, c int) {
				// 1. determine which cell the evaluation point is in
				xc, yc := math.Floor(x)+float64(c), math.Floor(y)+float64(r)

				// 2. generate a reproducible RNG for the cube (create seed by hashing)
				seed := hash(int(xc), int(yc)) % 256
				if seed < 0 {
					seed = -seed
				}
				// rng := rand.New(rand.NewSource(int64(seed)))

				// 3. determine how many feature points are in the cube
				npts := max
				selection := float64(p[seed]) / 255
				for i, cump := range cdf {
					if selection <= cump {
						npts = i
						break
					}
				}
				if npts > max { // clamp
					npts = max
				}

				// 4. place random feature points in the cube
				for ; npts > 0; npts-- {
					points = append(points, &point{
						xc + float64(p[seed+(npts*2)])/255,
						yc + float64(p[seed+(npts*2-1)])/255}) //point{xc + rng.Float64(), yc + rng.Float64()}
					// pchan <- point{xc + float64(p[seed+(npts*2)])/255, yc + float64(p[seed+(npts*2-1)])/255}
				}

				// wg.Done()
				// }(r, c)
			}
		}

		// wg.Wait()
		// close(pchan)
		// for p := range pchan {
		// 	points = append(points, p)
		// }

		if len(points) != 0 {
			// 5. find the closest N neighbors to the evaluation point,
			//    including in the 8 neighboring cells.
			tree := data.NewKDTree(2)
			tree.Build(points)

			// 6??? return N neighbors or their distances or?
			nearest := tree.NearestNeighbor(dist, x, y) // could but should not return nil
			return num.ClampFloat(dist([]float64{x, y}, nearest.Location()), 0, 1)
		}

		return 1 // default 'max' distance if no points
	}
}

// CellNoise2D creates a function which gives 2D cell noise by creating all the
// points at once, thereby running faster. Lambda and max determine
// the number of feature points in a given unit cube, which lambda is the average
// number per cell and max is the maximum number per cell. DistanceMetric dist
// provides the notion of 'distance'.
func CellNoise2D(maxX, maxY float64, lambda, maxPtsPerCell int, dist data.DistanceMetric) func(x, y float64) float64 {
	cdf := make([]float64, maxPtsPerCell+1, maxPtsPerCell+1)
	for k, t := 0, 0.0; k <= maxPtsPerCell; k++ {
		p := num.Poisson(lambda, k)
		t += p
		cdf[k] = t
	}

	// 5. find the closest N neighbors to the evaluation point,
	//    including in the 8 neighboring cells.
	tree := data.NewKDTree(2)
	tree.Build(MakeNoisePoints(maxX, maxY, maxPtsPerCell, cdf))

	return func(x, y float64) float64 {
		// 6??? return N neighbors or their distances or?
		nearest := tree.NearestNeighbor(dist, x, y) // could but should not return nil
		return num.ClampFloat(dist([]float64{x, y}, nearest.Location()), 0, 1)
	}
}

// MakeNoisePoints generates all the points for all the cells given the parameters.
func MakeNoisePoints(maxX, maxY float64, maxPtsPerCell int, cdf []float64) []data.Interface {
	points := make([]data.Interface, 0, 2000) // TODO: fix arbitrary size

	// for the cell and 8 cells in its neighborhood
	for yc := -1; yc < int(maxY+1); yc++ {
		for xc := -1; xc < int(maxX+1); xc++ {
			// 1. determine which cell the evaluation point is in
			// xc, yc := math.Floor(x)+float64(c), math.Floor(y)+float64(r)

			// 2. generate a reproducible RNG for the cube (create seed by hashing)
			seed := hash(int(xc), int(yc)) % 256
			if seed < 0 {
				seed = -seed
			}
			// rng := rand.New(rand.NewSource(int64(seed)))

			// 3. determine how many feature points are in the cube
			npts := maxPtsPerCell
			selection := float64(p[seed]) / 255
			for i, cump := range cdf {
				if selection <= cump {
					npts = i
					break
				}
			}
			if npts > maxPtsPerCell { // clamp
				npts = maxPtsPerCell
			}

			// 4. place random feature points in the cube
			for ; npts > 0; npts-- {
				points = append(points, &point{
					float64(xc) + float64(p[seed+(npts*2)])/255,
					float64(yc) + float64(p[seed+(npts*2-1)])/255,
				})
			}
		}
	}

	return points
}

// for 3d cell noise
type point3 [3]float64

func (p *point3) Location() []float64 {
	return p[:]
}

// CellNoise3D makes all the points at once, and therefore runs faster.
func CellNoise3D(maxX, maxY, maxZ float64, lambda, maxPtsPerCell int, dist data.DistanceMetric) func(x, y, z float64) float64 {
	cdf := make([]float64, maxPtsPerCell+1, maxPtsPerCell+1)
	for k, t := 0, 0.0; k <= maxPtsPerCell; k++ {
		p := num.Poisson(lambda, k)
		t += p
		cdf[k] = t
	}

	// 5. find the closest N neighbors to the evaluation point,
	//    including in the 8 neighboring cells.
	tree := data.NewKDTree(3)
	tree.Build(MakeNoisePoints3D(maxX, maxY, maxZ, maxPtsPerCell, cdf))
	fmt.Println("num points", tree.Len())

	return func(x, y, z float64) float64 {
		// 6??? return N neighbors or their distances or?
		nearest := tree.NearestNeighbor(dist, x, y, z) // could but should not return nil
		return num.ClampFloat(dist([]float64{x, y, z}, nearest.Location()), 0, 1)
	}
}

// MakeNoisePoints3D generates all the points for all the cells given the parameters.
func MakeNoisePoints3D(maxX, maxY, maxZ float64, maxPtsPerCell int, cdf []float64) []data.Interface {
	// totalCells := (maxX / cellSize * maxY / cellSize * maxZ / cellSize) // TODO: figure out better way to determine this
	points := make([]data.Interface, 0, 2000) //int(totalCells))
	fmt.Println("points 3d", len(points), cap(points))

	// for the cell and 8 cells in its neighborhood
	for yc := -1; yc < int(maxX+1); yc++ {
		for xc := -1; xc < int(maxX+1); xc++ {
			for zc := -1; zc < int(maxZ+1); zc++ {
				// 1. determine which cell the evaluation point is in
				// xc, yc := math.Floor(x)+float64(c), math.Floor(y)+float64(r)

				// 2. generate a reproducible RNG for the cube (create seed by hashing)
				seed := hash(hash(int(xc), int(yc)), int(zc)) % 256
				if seed < 0 {
					seed = -seed
				}
				// rng := rand.New(rand.NewSource(int64(seed)))

				// 3. determine how many feature points are in the cube
				npts := maxPtsPerCell
				selection := float64(p[seed]) / 255
				for i, cump := range cdf {
					if selection <= cump {
						npts = i
						break
					}
				}
				if npts > maxPtsPerCell { // clamp
					npts = maxPtsPerCell
				}

				// 4. place random feature points in the cube
				for ; npts > 0; npts-- {
					points = append(points, &point3{
						float64(xc) + float64(p[seed+(npts*3)])/255,
						float64(yc) + float64(p[seed+(npts*3-1)])/255,
						float64(zc) + float64(p[seed+(npts*3-2)])/255,
					})
				}
			}
		}
	}

	return points
}
