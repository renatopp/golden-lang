package debug

import (
	"fmt"
	"strings"

	"github.com/renatopp/golden/internal/builder"
	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/lang"
)

func esc(s string) string {
	return strings.ReplaceAll(s, "\n", "\\n")
}

func PrettyPrintTokens(module *builder.Module, tokens []*lang.Token) {
	fmt.Printf("Tokens for module %s:\n", module.Path)
	for _, token := range tokens {
		fmt.Printf("- %s (%s)\n", token.Kind, esc(token.Literal))
	}
	println()
}

func PrettyPrintAst(module *builder.Module, root *ast.Module) {
	fmt.Printf("AST for module %s:\n", module.Path)

	printer := NewAstPrinter()
	root.Accept(printer)
	println()
}
