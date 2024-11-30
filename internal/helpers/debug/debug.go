package debug

import (
	"fmt"
	"strings"

	"github.com/renatopp/golden/internal/builder"
	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/compiler/token"
)

func Escape(s string) string {
	return strings.ReplaceAll(s, "\n", "\\n")
}

func PrettyPrintTokens(file *builder.File, tokens []*token.Token) {
	fmt.Printf("Tokens for module %s:\n", file.Path)
	for _, token := range tokens {
		fmt.Printf("- [%s] (%s)\n", token.Display(), Escape(token.Literal))
	}
	println()
}

func PrettyPrintAst(file *builder.File, root ast.Module) {
	fmt.Printf("AST for module %s:\n", file.Path)
	root.Visit(NewAstPrinter())
	println()
}

// func PrettyPrintScope(scope *env.Scope) {
// 	if scope == nil {
// 		return
// 	}

// 	println("Scope:")
// 	for k, v := range scope.Types.Bindings {
// 		if v.Type == nil {
// 			fmt.Printf("- (T) %s → %s\n", k, "<nil>")
// 		} else {
// 			fmt.Printf("- (T) %s → %s\n", k, v.Type.GetSignature())
// 		}
// 	}
// 	for k, v := range scope.Values.Bindings {
// 		if v.Type == nil {
// 			fmt.Printf("- (V) %s → %s\n", k, "<nil>")
// 		} else {
// 			fmt.Printf("- (V) %s → %s\n", k, v.Type.GetSignature())
// 			// n := reflect.TypeOf(v.Node).String()
// 			// fmt.Printf("- (V) %s:%s → %s\n", k, n, v.Type.GetSignature())
// 		}
// 	}

// 	if scope.Parent != nil {
// 		print("\nParent ")
// 		PrettyPrintScope(scope.Parent)
// 	}

// 	println()
// }
