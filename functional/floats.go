package functional

import (
	"math"
)

// OneToOne takes a single argument and produces a single result.
type OneToOne func(float64) float64

// TwoToOne takes 2 arguments and produces a single result.
type TwoToOne func(float64, float64) float64

// ToBool takes a single argument and produces a boolean.
type ToBool func(float64) bool

// MapF applies f to every element in slice to produce a new slice.
func MapF(slice []float64, fns ...OneToOne) (out []float64) {
	for i := 0; i < len(slice); i++ {
		n := slice[i]
		for _, f := range fns {
			n = f(n)
		}
		out = append(out, n)
	}
	return
}

//// functions that can be used with Map ////

// Add produces a OneToOne that adds its argument by a number.
func Add(a float64) OneToOne {
	return func(b float64) float64 {
		return a + b
	}
}

// Mul produces a OneToOne that multiplies its argument by a number.
func Mul(a float64) OneToOne {
	return func(b float64) float64 {
		return a * b
	}
}

// Exp produces a OneToOne that exponentiates its argument by a number.
func Exp(exp float64) OneToOne {
	return func(base float64) float64 {
		return math.Pow(base, exp)
	}
}

// Inv produces a OneToOne that returns the multiplicative inverse of its argument.
func Inv() OneToOne {
	return func(n float64) float64 {
		return 1 / n
	}
}

// Truth produces a OneToOne that returns 1 if predicate is true for its
// argument, or 0 if predicate is false.
func Truth(predicate ToBool) OneToOne {
	return func(n float64) float64 {
		if predicate(n) {
			return 1
		}
		return 0
	}
}

// ReduceF applies f to every element in slice to produce an accumlated
// value `out`. in `f` the first argument is the accumulator, and the second
// argument is the current element from the slice.
func ReduceF(slice []float64, f TwoToOne) (out float64) {
	if len(slice) >= 1 {
		out = slice[0]
	}
	for i := 1; i < len(slice); i++ {
		out = f(out, slice[i])
	}
	return
}

//// functions that can be used with Reduce ////

// Sum is a TwoToOne that will add the accumulator and n.
func Sum(accumulator, n float64) float64 {
	return accumulator + n
}

// Prod is a TwoToOne that will multiply the accumulator and n.
func Prod(accumulator, n float64) float64 {
	return accumulator * n
}

// Max is a TwoToOne that will return the larger of its arguments.
func Max(max, n float64) float64 {
	return math.Max(max, n)
}

// Min is a TwoToOne that will return the smaller of its arguments.
func Min(min, n float64) float64 {
	return math.Min(min, n)
}

// FilterF applies f to every element in slice to produce a new slice containing
// the elements where f evaluates to true.
func FilterF(slice []float64, f ToBool) (out []float64) {
	for i := 0; i < len(slice); i++ {
		if f(slice[i]) {
			out = append(out, slice[i])
		}
	}
	return
}

//// functions that can be used with Filter ////

// LessThan returns a ToBool that returns true if its arg is less than a number.
func LessThan(a float64) ToBool {
	return func(b float64) bool {
		return b < a
	}
}

// GreaterThan returns a ToBool that returns true if its arg is greater than a number.
func GreaterThan(a float64) ToBool {
	return func(b float64) bool {
		return b > a
	}
}

// Equal returns a ToBool that returns true if its arg is equal to a number.
func Equal(a float64) ToBool {
	return func(b float64) bool {
		return b == a
	}
}

// Any returns a ToBool that returns true if any of the predicates evaluate to true.
func Any(predicates ...ToBool) ToBool {
	return func(n float64) bool {
		rval := false
		for _, f := range predicates {
			rval = rval || f(n)
		}
		return rval
	}
}

// All returns a ToBool that returns true if all of the predicates evaluate to true.
func All(predicates ...ToBool) ToBool {
	return func(n float64) bool {
		for _, f := range predicates {
			if !f(n) {
				return false
			}
		}
		return true
	}
}
