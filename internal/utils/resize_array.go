package utils

func ResizeArray[T any](arr []T, newSize int) []T {
	if newSize <= len(arr) {
		return arr[:newSize]
	}
	resized := make([]T, newSize)
	copy(resized, arr)
	for i := len(arr); i < newSize; i++ {
		var zeroValue T
		resized[i] = zeroValue
	}
	return resized
}
