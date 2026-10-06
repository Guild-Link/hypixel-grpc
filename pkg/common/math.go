package common

import "math"

func RoundToTwo(value float64) float64 {
	return math.Round(value*100) / 100
}
