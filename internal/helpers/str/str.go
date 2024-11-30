package str

import (
	"strings"

	"github.com/renatopp/golden/internal/helpers/iter"
)

func HumanList(items []string, last string) string {
	if len(items) == 0 {
		return ""
	}

	if len(items) == 1 {
		return items[0]
	}

	res := strings.Join(items[:len(items)-1], ", ") + " " + last + " " + items[len(items)-1]
	return res
}

func MapHumanList[T any](items []T, fn func(T) string, last string) string {
	strs := iter.Map(items, fn)
	return HumanList(strs, last)
}

// Pads the string with spaces to the left until the string reaches the desired
// length.
//
// If the string is already longer than the desired length, it will be returned
// as is.
func PadLeft(s string, n int) string {
	return PadLeftWith(s, n, " ")
}

// Pads the string with the specified character to the left until the string
// reaches the desired length.
//
// If the string is already longer than the desired length, it will be returned
// as is.
func PadLeftWith(s string, n int, with string) string {
	n = n - len(s)
	if n <= 0 {
		return s
	}
	return strings.Repeat(with, n) + s
}

// Pads the string with spaces to the right until the string reaches the
// desired
//
// If the string is already longer than the desired length, it will be returned
// as is.
func PadRight(s string, n int) string {
	return PadRightWith(s, n, " ")
}

// Pads the string with the specified character to the right until the string
// reaches the desired length.
//
// If the string is already longer than the desired length, it will be returned
// as is.
func PadRightWith(s string, n int, with string) string {
	n = n - len(s)
	if n <= 0 {
		return s
	}
	return s + strings.Repeat(with, n)
}

// Pads the string with spaces to the left and right until the string reaches
// the desired length. In case of an odd number of characters, the left side
// will have one more character than the right side.
//
// If the string is already longer than the desired length, it will be returned
// as is.
func PadCenter(s string, n int) string {
	return PadCenterWith(s, n, " ")
}

// Pads the string with the specified character to the left and right until the
// string reaches the desired length. In case of an odd number of characters,
// the left side will have one more character than the right side.
//
// If the string is already longer than the desired length, it will be returned
// as is.
func PadCenterWith(s string, n int, with string) string {
	n = n - len(s)
	if n <= 0 {
		return s
	}
	left := n / 2
	right := left
	if left%2 == 1 {
		left++
	}
	return strings.Repeat(with, left) + s + strings.Repeat(with, right)
}

// Repeats the string n times. If n is less than or equal to 0, an empty string
// will be returned.
func Repeat(s string, n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(s, n)
}
