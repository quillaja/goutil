package rand

import (
	"math"

	"github.com/quillaja/goutil/data"

	"github.com/quillaja/goutil/num"
)

const numCells = 3
const numZCells = 2

// deprecated
func hash(x, y int) int {
	//a >= b ? a * a + a + b : a + b * b;  where a, b >= 0
	a, b := x%(1<<16), y%(1<<16)
	if a >= b {
		return a*a + a + b
	}
	return a + b*b
}

// holds low and high bounds (or range, if that wasn't a keyword)
type bound struct {
	low, high int
}

func (b bound) in(n int) bool {
	return b.low <= n && n <= b.high
}

func (b bound) array() [2]int {
	return [2]int{b.low, b.high}
}

// 2D point for cell noise
type point [2]float64

func (p *point) Location() []float64 {
	return p[:]
}

// 3D point for 3d cell noise
type point3 [3]float64

func (p *point3) Location() []float64 {
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

	return func(x, y float64) float64 {
		nearest := math.Inf(0)

		// for the cell and 8 cells in its neighborhood
		for r := 1; r >= -1; r-- {
			for c := -1; c <= 1; c++ {
				// 1. determine which cell the evaluation point is in
				xc, yc := math.Floor(x)+float64(c), math.Floor(y)+float64(r)

				// 2. generate a reproducible RNG for the cube (create seed by hashing)
				seed := p[p[int(xc)&0xFF]+int(yc)&0xFF]

				// 3. determine how many feature points are in the cube
				npts := max
				selection := float64(p[seed]) / 256
				for i := 0; i < len(cdf); i++ {
					if selection <= cdf[i] {
						npts = i
						break
					}
				}
				npts = num.ClampInt(npts, 1, max)

				// 4. place random feature points in the cube and
				// 5. keep track of the closest one
				for ; npts > 0; npts-- {
					d := dist(
						[]float64{x, y},
						[]float64{
							xc + float64(p[seed+(npts*2)])/256,
							yc + float64(p[seed+(npts*2-1)])/256})
					if d < nearest {
						nearest = d
					}
				}

			}
		}

		// 6 return nearest distance, clamped
		return num.ClampFloat(nearest, 0, 1)
	}
}

// CellNoise2D contains the configuration data for a 2D cell noise generator.
type CellNoise2D struct {
	cdf            []float64
	tree           *data.KDTree
	xrange, yrange bound
	maxPtsPerCell  int
	dist           data.DistanceMetric
	perm           *[512]int
}

// NewCellNoise2D creates a pointer to a new CellNoise2D struct.
//
// Lambda and max determine the number of feature points in a given unit cube, where
// lambda is the average number per cell and maxPtsPerCell is the maximum number per cell.
// DistanceMetric dist provides the notion of 'distance'.
func NewCellNoise2D(seed int64, lambda, maxPtsPerCell int, dist data.DistanceMetric) *CellNoise2D {
	cdf := make([]float64, maxPtsPerCell+1, maxPtsPerCell+1)
	for k, t := 0, 0.0; k <= maxPtsPerCell; k++ {
		p := num.Poisson(lambda, k)
		t += p
		cdf[k] = t
	}

	conf := &CellNoise2D{
		cdf:           cdf,
		tree:          data.NewKDTree(2),
		xrange:        bound{0, numCells},
		yrange:        bound{0, numCells},
		maxPtsPerCell: maxPtsPerCell,
		dist:          dist,
		perm:          MakePermutation(seed),
	}

	conf.tree.Build(MakeNoisePoints2D(
		conf.xrange.array(), conf.yrange.array(),
		conf.maxPtsPerCell, conf.cdf, conf.perm))

	return conf
}

// Noise gets a noise value at the given point (x, y).
func (conf *CellNoise2D) Noise(x, y float64) float64 {
	// if x and y go past already calculated ranges of cells,
	// new cells need to be calculated and added.
	rebuild := false
	xc, yc := int(x), int(y) // calculate current cell x,y is in
	if !conf.xrange.in(xc) {
		rebuild = true
		offset := xc / numCells // find which 'group' of precalc cells we're in
		conf.xrange.low = (numCells * offset) + 2
		conf.xrange.high = numCells * (offset + 1)
		// ^ adding 2 to min range is a hack to avoid generating duplicate
		// points when the range is expanded, since in MakeNoisePoints()
		// the cells are looped -1 to +1 the range.
	}
	if !conf.yrange.in(yc) {
		rebuild = true
		offset := yc / numCells
		conf.yrange.low = (numCells * offset) + 2
		conf.yrange.high = numCells * (offset + 1)
	}

	if rebuild {
		// add new points to old points, then rebuild tree
		// start := time.Now() // for debug info
		newPoints := MakeNoisePoints2D(conf.xrange.array(), conf.yrange.array(),
			conf.maxPtsPerCell, conf.cdf, conf.perm)
		conf.tree.Build(append(conf.tree.Items(), newPoints...))
		// expand ranges by resetting the min
		conf.xrange.low, conf.yrange.low = 0, 0
		// show debug info
		// fmt.Println("add+rebuilt tree", x, y, "expanded ranges:", xrange, yrange)
		// fmt.Println(" num points", tree.Len(), "rebuild took (ms):", time.Since(start).Seconds()*1000)
		// fmt.Println(" theoretical size of points (KB):", tree.Len()*16/1024) // 8bytes per float64 * 2 per point * points / bytes/kb
	}

	// 6??? return dist to nearest neighbor (or Nth nearest, or points themselves...or?)
	nearest := conf.tree.NearestNeighbor(conf.dist, x, y) // could but should not return nil
	return num.ClampFloat(conf.dist([]float64{x, y}, nearest.Location()), 0, 1)

}

// MakeNoisePoints2D generates all the points for all the cells given the parameters.
func MakeNoisePoints2D(xrange, yrange [2]int, maxPtsPerCell int, cdf []float64, p *[512]int) []data.Interface {
	points := make([]data.Interface, 0, 500) // TODO: fix arbitrary size
	// for the cell and 8 cells in its neighborhood
	for yc := yrange[0] - 1; yc <= yrange[1]+1; yc++ {
		for xc := xrange[0] - 1; xc <= xrange[1]+1; xc++ {
			// 1. determine which cell the evaluation point is in
			// xc, yc := math.Floor(x)+float64(c), math.Floor(y)+float64(r)

			// 2. generate a reproducible RNG for the cube (create seed by hashing)
			seed := p[p[xc&0xFF]+yc&0xFF]
			// if seed < 0 {
			// 	seed = -seed
			// }

			// 3. determine how many feature points are in the cube
			npts := maxPtsPerCell
			selection := float64(p[seed]) / 256
			for i, cump := range cdf {
				if selection <= cump {
					npts = i
					break
				}
			}
			npts = num.ClampInt(npts, 1, maxPtsPerCell)

			// 4. place random feature points in the cube
			for ; npts > 0; npts-- {
				points = append(points, &point{
					float64(xc) + float64(p[seed+(npts*2)])/256,
					float64(yc) + float64(p[seed+(npts*2-1)])/256,
				})
			}
		}
	}

	return points
}

// MakeNoisePoints3D generates all the points for all the cells given the parameters.
func MakeNoisePoints3D(xrange, yrange, zrange [2]int, maxPtsPerCell int, cdf []float64, p *[512]int) []data.Interface {
	// totalCells := (maxX / cellSize * maxY / cellSize * maxZ / cellSize) // TODO: figure out better way to determine this
	points := make([]data.Interface, 0, 5000) //int(totalCells))

	// for the cell and 8 cells in its neighborhood
	for zc := zrange[0] - 1; zc <= zrange[1]+1; zc++ {
		for yc := yrange[0] - 1; yc <= yrange[1]+1; yc++ {
			for xc := xrange[0] - 1; xc <= xrange[1]+1; xc++ {
				// 1. determine which cell the evaluation point is in
				// xc, yc := math.Floor(x)+float64(c), math.Floor(y)+float64(r)

				// 2. generate a reproducible RNG for the cube (create seed by hashing)
				seed := p[p[p[xc&0xFF]+yc&0xFF]+zc&0xFF]
				// seed := hash(hash(xc, yc), zc) % 256
				// if seed < 0 {
				// 	seed = -seed
				// }

				// 3. determine how many feature points are in the cube
				npts := maxPtsPerCell
				selection := float64(p[seed]) / 256
				for i, cump := range cdf {
					if selection <= cump {
						npts = i
						break
					}
				}
				npts = num.ClampInt(npts, 1, maxPtsPerCell)

				// 4. place random feature points in the cube
				// can't multiply offset with 3 because the p-array is
				// only 512 long, and seed is [0,255]
				for ; npts > 0; npts-- {
					points = append(points, &point3{
						float64(xc) + float64(p[seed+(npts*2)])/256,
						float64(yc) + float64(p[seed+(npts*2-1)])/256,
						float64(zc) + float64(p[seed+(npts*2-2)])/256,
					})
				}
			}
		}
	}

	return points
}

// CellNoise3D contains the configuration data for a 3D cell noise generator.
type CellNoise3D struct {
	cdf                    []float64
	tree                   *data.KDTree
	xrange, yrange, zrange bound
	maxPtsPerCell          int
	dist                   data.DistanceMetric
	perm                   *[512]int
}

// NewCellNoise3D creates a pointer to a new CellNoise3D struct.
//
// Lambda and max determine the number of feature points in a given unit cube, where
// lambda is the average number per cell and maxPtsPerCell is the maximum number per cell.
// DistanceMetric dist provides the notion of 'distance'.
func NewCellNoise3D(seed int64, lambda, maxPtsPerCell int, dist data.DistanceMetric) *CellNoise3D {
	cdf := make([]float64, maxPtsPerCell+1, maxPtsPerCell+1)
	for k, t := 0, 0.0; k <= maxPtsPerCell; k++ {
		p := num.Poisson(lambda, k)
		t += p
		cdf[k] = t
	}

	conf := &CellNoise3D{
		cdf:           cdf,
		tree:          data.NewKDTree(3),
		xrange:        bound{0, numCells},
		yrange:        bound{0, numCells},
		zrange:        bound{0, numZCells},
		maxPtsPerCell: maxPtsPerCell,
		dist:          dist,
		perm:          MakePermutation(seed),
	}

	conf.tree.Build(MakeNoisePoints3D(
		conf.xrange.array(), conf.yrange.array(), conf.zrange.array(),
		conf.maxPtsPerCell, conf.cdf, conf.perm))

	return conf
}

// Noise generates a noise value at the (x,y,z) location.
func (conf *CellNoise3D) Noise(x, y, z float64) float64 {
	rebuild := false
	zrebuild := false
	xc, yc, zc := int(x), int(y), int(z)
	if !conf.xrange.in(xc) {
		rebuild = true
		offset := xc / numCells
		conf.xrange.low = (numCells * offset) + 2
		conf.xrange.high = numCells * (offset + 1)
	}
	if !conf.yrange.in(yc) {
		rebuild = true
		offset := yc / numCells
		conf.yrange.low = (numCells * offset) + 2
		conf.yrange.high = numCells * (offset + 1)
	}
	if !conf.zrange.in(zc) {
		rebuild, zrebuild = true, true
		offset := zc / numZCells
		conf.zrange.low = (numZCells * offset)
		conf.zrange.high = numZCells * (offset + 1)
		// z range doesn't need the same hack as x and y to avoid generating
		// duplicate points, because expansions in the z range will rebuild
		// the entire tree
	}

	if rebuild {
		// start := time.Now() // used with debug info
		newPoints := MakeNoisePoints3D(
			conf.xrange.array(), conf.yrange.array(), conf.zrange.array(),
			conf.maxPtsPerCell, conf.cdf, conf.perm)
		// if the rebuild is because of a change in the z direction, then all points
		// in the tree will be rebuilt. x and y will still generate in the expanded
		// ranges established in previous iterations.
		if zrebuild {
			conf.tree.Build(newPoints)
		} else {
			conf.tree.Build(append(conf.tree.Items(), newPoints...))
		}
		// expand x and y ranges. z range doesn't expand.
		conf.xrange.low, conf.yrange.low = 0, 0
		// debug info
		// fmt.Println("add+rebuilt tree (", x, y, z, ") | expanded ranges:", conf.xrange, conf.yrange, conf.zrange)
		// fmt.Println(" num points", conf.tree.Len(), "rebuild took (ms):", time.Since(start).Seconds()*1000)
		// fmt.Println(" theoretical size of points (KB):", conf.tree.Len()*24/1024) // 8bytes per float64 * 3 per point * points / bytes/kb
	}

	// 6??? return nearest neighbor, Nth nearest, or their distances or?
	nearest := conf.tree.NearestNeighbor(conf.dist, x, y, z) // could but should not return nil
	return num.ClampFloat(conf.dist([]float64{x, y, z}, nearest.Location()), 0, 1)
}

// DEPRECATED
// Old original implementation using function closure. Algorithm is same as struct-based approach
// CellNoise3D makes all the points at once, and therefore runs faster.
// func CellNoise3D(lambda, maxPtsPerCell int, dist data.DistanceMetric) func(x, y, z float64) float64 {
// 	cdf := make([]float64, maxPtsPerCell+1, maxPtsPerCell+1)
// 	for k, t := 0, 0.0; k <= maxPtsPerCell; k++ {
// 		p := num.Poisson(lambda, k)
// 		t += p
// 		cdf[k] = t
// 	}

// 	// 5. find the closest N neighbors to the evaluation point,
// 	//    including in the 8 neighboring cells.
// 	tree := data.NewKDTree(3)
// 	xrange, yrange := [2]int{0, numCells}, [2]int{0, numCells}
// 	zrange := [2]int{0, numZCells}
// 	tree.Build(MakeNoisePoints3D(xrange, yrange, zrange, maxPtsPerCell, cdf, &p))

// 	return func(x, y, z float64) float64 {
// 		rebuild := false
// 		zrebuild := false
// 		xc, yc, zc := int(x), int(y), int(z)
// 		if !(xrange[0] <= xc && xc <= xrange[1]) {
// 			rebuild = true
// 			offset := xc / numCells
// 			xrange = [2]int{(numCells * offset) + 2, numCells * (offset + 1)}
// 		}
// 		if !(yrange[0] <= yc && yc <= yrange[1]) {
// 			rebuild = true
// 			offset := yc / numCells
// 			yrange = [2]int{(numCells * offset) + 2, numCells * (offset + 1)}
// 		}
// 		if !(zrange[0] <= zc && zc <= zrange[1]) {
// 			rebuild, zrebuild = true, true
// 			offset := zc / numZCells
// 			zrange = [2]int{(numZCells * offset), numZCells * (offset + 1)}
// 			// z range doesn't need the same hack as x and y to avoid generating
// 			// duplicate points, because expansions in the z range will rebuild
// 			// the entire tree
// 		}
// 		if rebuild {
// 			// start := time.Now()
// 			newPoints := MakeNoisePoints3D(xrange, yrange, zrange, maxPtsPerCell, cdf, &p)
// 			// if the rebuild is because of a change in the z direction, then all points
// 			// in the tree will be rebuilt. x and y will still generate in the expanded
// 			// ranges established in previous iterations.
// 			if zrebuild {
// 				tree.Build(newPoints)
// 			} else {
// 				tree.Build(append(tree.Items(), newPoints...))
// 			}
// 			// expand x and y ranges. z range doesn't expand.
// 			xrange[0], yrange[0] = 0, 0
// 			// fmt.Println("add+rebuilt tree (", x, y, z, ") | expanded ranges:", xrange, yrange, zrange)
// 			// fmt.Println(" num points", tree.Len(), "rebuild took (ms):", time.Since(start).Seconds()*1000)
// 			// fmt.Println(" theoretical size of points (KB):", tree.Len()*24/1024) // 8bytes per float64 * 3 per point * points / bytes/kb
// 		}

// 		// 6??? return N neighbors or their distances or?
// 		nearest := tree.NearestNeighbor(dist, x, y, z) // could but should not return nil
// 		return num.ClampFloat(dist([]float64{x, y, z}, nearest.Location()), 0, 1)
// 	}
// }
