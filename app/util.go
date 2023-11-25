package app

import "math"

func EaseOutSine(x float64) float64 {
	if x >= 1 {
		return 1
	}

	return 1 - math.Cos(x*math.Pi/2)
}
