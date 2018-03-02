package main

import (
	"fmt"
	"math"

	f "github.com/quillaja/goutil/functional"
)

func main() {
	nums := []float64{9, 8, 2, 1, 6, 5.5}
	fmt.Println("original nums[]", nums)

	// test Map
	fmt.Println("add 1 to nums[]", f.MapF(nums, f.Add(1)))
	fmt.Println("times 2 and subtract 10", f.MapF(nums, f.Mul(2), f.Add(-10)))
	// Map using funcs in the go standard library
	fmt.Println("sin(x)", f.MapF(nums, math.Sin))
	fmt.Println("floor(x)", f.MapF(nums, math.Floor))

	// test complex composition
	fmt.Println("how many are in [1,6)",
		f.ReduceF(
			f.MapF(nums,
				f.Truth(
					f.Any(
						f.All(f.GreaterThan(1), f.Equal(1)),
						f.LessThan(6)))),
			f.Sum))

	// could also do this
	betweenIncExc := func(nums []float64, low, high float64) float64 {
		return f.ReduceF(
			f.MapF(nums,
				f.Truth(
					f.Any(
						f.All(f.GreaterThan(low), f.Equal(low)),
						f.LessThan(high)))),
			f.Sum)
	}

	fmt.Println("how many are in [1,6)", betweenIncExc(nums, 1, 6))

	// test filter
	fmt.Println("those in [1,6)", f.FilterF(nums, f.Any(
		f.All(f.GreaterThan(1), f.Equal(6)), f.LessThan(6))))

	fmt.Println("those == 5.5", f.FilterF(nums, f.Equal(5.5)))

	// test Reduce
	fmt.Println("min of set", f.ReduceF(nums, f.Min))
	fmt.Println("max of set", f.ReduceF(nums, f.Max))
	fmt.Println("avg of set", f.ReduceF(nums, f.Sum)/float64(len(nums)))
	fmt.Println("product of set", f.ReduceF(nums, f.Prod))
	fmt.Println("min using std lib", f.ReduceF(nums, math.Min))

}
