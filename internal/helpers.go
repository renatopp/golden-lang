package internal

import (
	"fmt"
	"strings"

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

func ident(s string, i int) string {
	return strings.ReplaceAll(s, "\n", "\n"+strings.Repeat("  ", i))
}

func parseDelimeter[T any](p *parser, openKind, closeKind string, parseFn func() T) T {
	p.ExpectTokens(openKind)
	p.EatToken()
	res := parseFn()
	p.ExpectTokens(closeKind)
	p.EatToken()
	return res
}

func parseOptionalDelimeter[T any](p *parser, openKind, closeKind string, parseFn func() T) T {
	if p.IsNextTokens(openKind) {
		return parseDelimeter(p, openKind, closeKind, parseFn)
	}
	return parseFn()
}
