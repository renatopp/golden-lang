package debug

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/renatopp/golden/internal/core"
	"github.com/renatopp/golden/lang"
)

func PrintTokens(tokens []*lang.Token) {
	for _, t := range tokens {
		fmt.Printf("    - %s: %q\n", t.Kind, t.Literal)
	}
	println()
}

func PrintAst(root *core.AstNode) {
	root.Traverse(func(node *core.AstNode, level int) {
		ident := strings.Repeat("  ", level)
		line := "    " + ident + node.Signature()
		comment := " -- " + ident + node.Tag()

		size := utf8.RuneCountInString(line)
		if size < 50 {
			println(line, strings.Repeat(" ", 50-size), comment)
		} else {
			println(line, comment)
		}
	})
	println()
}
