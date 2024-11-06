package internal

import (
	"fmt"

	"github.com/renatopp/golden/lang/strutils"
)

var f = fmt.Sprintf
var esc = strutils.EscapeNewlines

func appendAll[T any](arrays ...[]T) []T {
	var out []T
	for _, arr := range arrays {
		out = append(out, arr...)
	}
	return out
}

func isData[T any](node *Node) bool {
	_, ok := node.Data.(T)
	return ok
}
