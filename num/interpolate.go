// Package num provides some numeric tools such as interpolation and clamping.
package num

import "math"

// NormalizeAngle takes an angle in radians and scales it to [-2PI,2PI].
func NormalizeAngle(theta float64) float64 {
	if -2*math.Pi <= theta && theta <= 2*math.Pi {
		return theta
	}
	f := theta / (2 * math.Pi)
	if f < 0 {
		return 2 * math.Pi * (f - math.Ceil(f))
	}
	return 2 * math.Pi * (f - math.Floor(f))
}

// NormalizeAngleDeg takes an angle in degrees and scales it to [-360, 360].
func NormalizeAngleDeg(theta float64) float64 {
	if -360 <= theta && theta <= 360 {
		return theta
	}
	f := theta / 360
	if f < 0 {
		return 360 * (f - math.Ceil(f))
	}
	return 360 * (f - math.Floor(f))
}

// UnitLerp does a linear interpolation of x to between toMin and toMax, with the
// assumption that x is in the range [0.0, 1,0]. Essentially, a special case
// of general linear interpolation.
func UnitLerp(x, toMin, toMax float64) float64 {
	return (1.0-x)*toMin + x*toMax
}

// ReverseUnitLerp interpolates an x in the range [xMin, xMax] to the range [0,1].
func ReverseUnitLerp(x, xMin, xMax float64) float64 {
	return (x - xMax) / (xMax - xMin)
}

// Lerp does a linear interpolation of x to between toMin and toMax where xMin
// and xMax are the lower and upper bounds of x.
func Lerp(x, xMin, xMax, toMin, toMax float64) float64 {
	return toMin + (toMax-toMin)*((xMax-x)/(xMax-xMin))
}

// SmoothStep uses a sigmoid-like iterpolation to produce a smooth interpolation
// of x to the range [0,1] when x is also in the range [0,1]. To interpolate
// any x into a suitable argument for this function, use ReverseUnitLerp() first.
// See: https://en.wikipedia.org/wiki/Smoothstep
func SmoothStep(x float64) float64 {
	if x <= 0 {
		return 0
	}
	if 1 <= x {
		return 1
	}
	return x * x * (3 - 2*x)
}

// Sigmoid returns the interpolation of x to the range [0,1] according to the
// logistic function S(x) = 1 / (1 + e^-x).
// See: https://en.wikipedia.org/wiki/Sigmoid_function
func Sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

// ClampFloat clamps x between min and max.
func ClampFloat(x, min, max float64) float64 {
	if x <= min {
		return min
	}
	if x >= max {
		return max
	}
	return x
}
