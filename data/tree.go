package data

import (
	"math"
)

// Interface defines methods required for type wishing to use SpacialTree.
type Interface interface {
	// get the n-dimensional location of an item
	Location() []float64
}

// SpacialTree defines the interface for use in n-dimension space
// partitoning trees, such as k-d trees and quad/oct-trees.
type SpacialTree interface {
	// get the number of dimensions used in the tree
	Dimensions() int
	// get the total number of items in the tree
	Len() int
	// get a list of items in the tree
	Items() []Interface
	// build tree from item(s)
	Build(items []Interface)
	// check if item is in the tree, based on its location
	QueryPoint(item Interface) bool
	// get all items within the region defined by the list of mins/maxs
	QueryRange(ranges [][2]float64) []Interface
	// get the 1 nearest neighbor
	NearestNeighbor(dist DistanceMetric, point ...float64) Interface
	// get the k nearest neighbors. may return fewer than k
	NearestNeighbors(dist DistanceMetric, k int, point ...float64) []Interface
}

// DistanceMetric is a type of function that calculates the distance
// between 2 n-dimensional points. Both arguments to the function should
// be equal length.
type DistanceMetric func([]float64, []float64) float64

// EuclideanSq is a DistanceMetric func which computes the
// euclidean/cartesian/geometric distance. It actually returns the sum of
// squares, without taking the square root.
func EuclideanSq(a, b []float64) float64 {
	if len(a) != len(b) {
		panic("a and b are different lengths")
	}
	sum := 0.0
	for i := 0; i < len(a); i++ {
		diff := b[i] - a[i]
		sum += diff * diff
	}
	return sum
}

// Euclidean is the same as EuclideanSq() but takes the square root.
func Euclidean(a, b []float64) float64 {
	return math.Sqrt(EuclideanSq(a, b))
}

// Manhattan is a DistanceMetric func which computes the
// manhattan/taxi cab/snake distance.
func Manhattan(a, b []float64) float64 {
	if len(a) != len(b) {
		panic("a and b are different lengths")
	}
	sum := 0.0
	for i := 0; i < len(a); i++ {
		diff := b[i] - a[i]
		sum += math.Abs(diff)
	}
	return sum
}

// Chebyshev is a DistanceMetric func which computes the chebyshev distance,
// where the distance is the single most significant of the components.
func Chebyshev(a, b []float64) float64 {
	if len(a) != len(b) {
		panic("a and b are different lengths")
	}
	max := 0.0
	for i := 0; i < len(a); i++ {
		diff := b[i] - a[i]
		max = math.Max(max, math.Abs(diff))
	}
	return max
}

// Canberra is a DistanceMetric func which computes the canberra distance:
// Sum( |b-a| / |b|+|a| )
func Canberra(a, b []float64) float64 {
	if len(a) != len(b) {
		panic("a and b are different lengths")
	}
	sum := 0.0
	for i := 0; i < len(a); i++ {
		diff := b[i] - a[i]
		sum += math.Abs(diff) / (math.Abs(b[i]) + math.Abs(a[i]))
	}
	return sum
}
