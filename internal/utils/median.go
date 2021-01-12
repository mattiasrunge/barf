package utils

import "sort"

// Median returns the median value of the supplied numbers
func Median(numbers []float64) float64 {
	sort.Float64s(numbers)

	index := len(numbers) / 2

	if len(numbers)%2 != 0 {
		return numbers[index]
	}

	return (numbers[index-1] + numbers[index]) / 2
}
