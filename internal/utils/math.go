package utils

func Clamp(value, minValue, maxValue int32) int32 {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

