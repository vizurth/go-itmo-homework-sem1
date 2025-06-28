package go_digest

import (
	"math"
	"math/cmplx"
	"math/rand"
	"strings"
	"unsafe"
)

// GetCharByIndex returns the i-th character from the given string.
func GetCharByIndex(str string, idx int) rune {
	if len(str) == 0 {
		panic("empty string")
	} else if idx > len(str) {
		panic("index out of range")
	} else if idx < 0 {
		panic("index out of range")
	}

	for _, r := range str {
		if idx == 0 {
			return r

		}
		idx--
	}
	return 0
}

// GetStringBySliceOfIndexes returns a string formed by concatenating specific characters from the input string based
// on the provided indexes.
func GetStringBySliceOfIndexes(str string, indexes []int) string {
	b := strings.Builder{}
	b.Grow(len(indexes))
	runes := []rune(str)
	for _, idx := range indexes {
		b.WriteRune(runes[idx])
	}
	return b.String()
}

// ShiftPointer shifts the given pointer by the specified number of bytes using unsafe.Add.
func ShiftPointer(pointer **int, shift int) {
	*pointer = (*int)(unsafe.Add(unsafe.Pointer(*pointer), shift))
}

// IsComplexEqual compares two complex numbers and determines if they are equal.
func IsComplexEqual(a, b complex128) bool {
	const epsilon = 1e-5

	ra, rb := real(a), real(b)
	ia, ib := imag(a), imag(b)

	// если хоть одна часть NaN — всегда false
	if math.IsNaN(ra) || math.IsNaN(rb) || math.IsNaN(ia) || math.IsNaN(ib) {
		return false
	}

	// если обе части Inf и равны — считаем равными
	if ra == rb && ia == ib {
		return true
	}

	return math.Abs(ra-rb) < epsilon && math.Abs(ia-ib) < epsilon
}

// GetRootsOfQuadraticEquation returns two roots of a quadratic equation ax^2 + bx + c = 0.
func GetRootsOfQuadraticEquation(a, b, c float64) (complex128, complex128) {
	D := complex(b*b-4*a*c, 0)
	sqrtD := cmplx.Sqrt(D)
	denom := complex(2*a, 0)

	x1 := (-complex(b, 0) + sqrtD) / denom
	x2 := (-complex(b, 0) - sqrtD) / denom

	return x1, x2
}

// Sort sorts in-place the given slice of integers in ascending order.
//func Sort(source []int) {
//	sort.Ints(source)
//}

func Sort(nums []int) {
	quickSort(nums, 0, len(nums)-1)
}

func quickSort(nums []int, left, right int) {
	if left >= right {
		return
	}

	pivot := partition(nums, left, right)
	quickSort(nums, left, pivot-1)
	quickSort(nums, pivot+1, right)
}

func partition(nums []int, left, right int) int {
	// Random pivot to avoid worst-case
	randIndex := left + rand.Intn(right-left+1)
	nums[right], nums[randIndex] = nums[randIndex], nums[right]

	pivot := nums[right]
	i := left - 1

	for j := left; j < right; j++ {
		if nums[j] <= pivot {
			i++
			nums[i], nums[j] = nums[j], nums[i]
		}
	}

	nums[i+1], nums[right] = nums[right], nums[i+1]
	return i + 1
}

// ReverseSliceOne in-place reverses the order of elements in the given slice.
func ReverseSliceOne(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// ReverseSliceTwo returns a new slice of integers with elements in reverse order compared to the input slice.
// The original slice remains unmodified.
func ReverseSliceTwo(s []int) []int {
	reversed := make([]int, 0)
	for i := len(s) - 1; i >= 0; i-- {
		reversed = append(reversed, s[i])
	}
	return reversed
}

// SwapPointers swaps the values of two pointers.
func SwapPointers(a, b *int) {
	*a, *b = *b, *a
}

// IsSliceEqual compares two slices of integers and returns true if they contain the same elements in the same order.
func IsSliceEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	} else if len(a) == 0 && len(b) == 0 {
		return true
	} else if (a == nil && len(b) == 0) || (len(a) == 0 && b == nil) {
		return true
	} else if len(a) == len(b) {
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
	}
	return true
}

// DeleteByIndex deletes the element at the specified index from the slice and returns a new slice.
// The original slice remains unmodified.
func DeleteByIndex(s []int, idx int) []int {
	var newSlice []int
	newSlice = append(newSlice, s[:idx]...)
	newSlice = append(newSlice, s[idx+1:]...)
	return newSlice
}
