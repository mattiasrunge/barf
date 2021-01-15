package utils

import "sort"

// Median returns the median value of the supplied numbers
func Median(numbers []float64) float64 {
	worknumbers := make([]float64, len(numbers))
	copy(worknumbers, numbers)

	sort.Float64s(worknumbers)

	index := len(worknumbers) / 2

	if len(worknumbers)%2 != 0 {
		return worknumbers[index]
	}

	return (worknumbers[index-1] + worknumbers[index]) / 2
}
