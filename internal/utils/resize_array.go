package utils

import (
	"fmt"
)

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

func main() {
	intSlice := []int{1, 2, 3}
	resizedIntSlice := ResizeArray(intSlice, 5)
	fmt.Println(resizedIntSlice)

	stringSlice := []string{"a", "b"}
	resizedStringSlice := ResizeArray(stringSlice, 4)
	fmt.Println(resizedStringSlice)

	type Bitset32 struct {
		value uint32
	}

	bitsetSlice := []*Bitset32{&Bitset32{value: 1}, &Bitset32{value: 2}}
	resizedBitsetSlice := ResizeArray(bitsetSlice, 5)
	fmt.Println(resizedBitsetSlice)
}
