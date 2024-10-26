package strutils

import (
	"strings"
	"unicode"
)

// Replaces all occurrences of the "\n" character with the escaped version
// ("\\n").
func EscapeNewlines(s string) string {
	return strings.ReplaceAll(s, "\n", "\\n")
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

// Reverses the string.
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Replaces all occurrences of the old string with the new string.
func Replace(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

// Converts the string to lowercase. Example: "Hello World" -> "hello world".
func ToLowerCase(s string) string {
	return strings.ToLower(s)
}

// Converts the string to uppercase. Example: "Hello World" -> "HELLO WORLD".
func ToUpperCase(s string) string {
	return strings.ToUpper(s)
}

// Converts the string to title case. Example: "hello world" -> "Hello World".
func ToTitleCase(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		runes := []rune(word)
		if len(runes) > 0 {
			runes[0] = unicode.ToUpper(runes[0])
			for j := 1; j < len(runes); j++ {
				runes[j] = unicode.ToLower(runes[j])
			}
		}
		words[i] = string(runes)
	}
	return strings.Join(words, " ")
}

// TODO:
// 		camelCase
// 		PascalCase
// 		snake_case
// 		kebab-case
// 		Train-Case
// 		flatcase
// 		dot.case
// 		path/case
// 		Sentence case
