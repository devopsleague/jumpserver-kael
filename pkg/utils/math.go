package utils

func GetMaxInt() int {
	maxInt := int(^uint(0) >> 1)
	return maxInt
}
