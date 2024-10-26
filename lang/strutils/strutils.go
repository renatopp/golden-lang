package strutils

import "strings"

func Escape(s string) string {
	return strings.ReplaceAll(s, "\n", "\\n")
}

func PadLeft(s string, n int) string {
	return strings.Repeat(" ", n) + s
}

func FillLeft(s string, n int, with string) string {
	return strings.Repeat(with, n-len(s)) + s
}

func PadRight(s string, n int) string {
	return s + strings.Repeat(" ", n)
}

func FillRight(s string, n int, with string) string {
	return s + strings.Repeat(with, n-len(s))
}

func PadCenter(s string, n int) string {
	return strings.Repeat(" ", n/2) + s + strings.Repeat(" ", n/2)
}

func FillCenter(s string, n int, with string) string {
	return strings.Repeat(with, (n-len(s))/2) + s + strings.Repeat(with, (n-len(s))/2)
}

func Repeat(s string, n int) string {
	return strings.Repeat(s, n)
}
