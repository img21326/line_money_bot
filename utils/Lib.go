package utils

func Abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
