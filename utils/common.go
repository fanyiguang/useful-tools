package utils

func InArray[T comparable](val T, arr []T) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}
