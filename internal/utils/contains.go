package utils

func IntArrayContains(a []int, value int) bool {
	for _, av := range a {
		if av == value {
			return true
		}
	}

	return false
}
