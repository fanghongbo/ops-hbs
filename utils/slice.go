package utils

func HasContainInt64(list []int64, target int64) bool {
	for _, b := range list {
		if b == target {
			return true
		}
	}
	return false
}

func IsSameSlice(a []int64, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i, val := range a {
		if val != b[i] {
			return false
		}
	}
	return true
}

func HasContainSlice(a []int64, b []int64) bool {
	for _, i := range a {
		if !HasContainInt64(b, i) {
			return false
		}
	}
	return true
}
