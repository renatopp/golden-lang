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
