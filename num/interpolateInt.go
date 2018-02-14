package num

// ClampInt clamps x between min and max.
func ClampInt(x, min, max int) int {
	if x <= min {
		return min
	}
	if x >= max {
		return max
	}
	return x
}
