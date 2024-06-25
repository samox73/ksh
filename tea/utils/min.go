package utils

import "math"

func minInt(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}