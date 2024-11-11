package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/renatopp/golden/internal/compiler/syntax"
	"github.com/renatopp/golden/internal/compiler/syntax/ast"
	"github.com/renatopp/golden/internal/core"
)

type Debug struct{}

func (c *Debug) Name() string {
	return "debug"
}

func (c *Debug) Description() string {
	return "Debug"
}

func (c *Debug) Help() string {
	return "Debug"
}

func (c *Debug) Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no file specified")
	}

	filename, err := filepath.Abs(args[0])
	if err != nil {
		return fmt.Errorf("error getting absolute path: %v", err)
	}

	file, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	println("## Lexer Output:\n")
	tokens, err := syntax.Lex(file)
	if err != nil {
		return fmt.Errorf("lexing file:\n%v", err)
	}

	for _, t := range tokens {
		fmt.Printf("- %s: %q\n", t.Kind, t.Literal)
	}
	println("\n")

	println("## Parser Output:\n")
	root, err := syntax.Parse(tokens)
	if err != nil {
		return fmt.Errorf("parsing file:\n%v", err)
	}

	ast := root.Data().(*ast.Module)
	for _, imp := range ast.Imports {
		if imp.Alias != "" {
			fmt.Printf("import %s as %s\n", imp.Path, imp.Alias)
		} else {
			fmt.Printf("import %s\n", imp.Path)
		}
	}

	for _, decl := range ast.Types {
		decl.Traverse(printNode)
	}
	for _, decl := range ast.Functions {
		decl.Traverse(printNode)
	}
	for _, decl := range ast.Variables {
		decl.Traverse(printNode)
	}

	println("\n")

	// println("## Analyzer Output:\n")

	// scope := internal.NewScope()
	// scope.SetType("Void", internal.Void)
	// scope.SetType("Bool", internal.Bool)
	// scope.SetType("Int", internal.Int)
	// scope.SetType("Float", internal.Float)
	// scope.SetType("String", internal.String)

	// module := internal.NewModule()
	// module.Scope = scope
	// module.Ast = ast
	// module.Node = root
	// module.Analyzer = internal.NewAnalyzer(module)
	// err = module.Analyzer.Analyze()
	// if err != nil {
	// 	return fmt.Errorf("analyzing module:\n%v", err)
	// }

	// for _, decl := range ast.Types {
	// 	println(decl.String())
	// }
	// for _, decl := range ast.Functions {
	// 	println(decl.String())
	// }
	// for _, decl := range ast.Variables {
	// 	println(decl.String())
	// }

	// println()
	// println(module.Scope.String())

	return nil
}

func printNode(node *core.AstNode, level int) {
	ident := strings.Repeat("  ", level)
	line := ident + node.Signature()
	comment := " -- " + ident + node.Tag()

	size := utf8.RuneCountInString(line)
	if size < 30 {
		println(line, strings.Repeat(" ", 30-utf8.RuneCountInString(line)), comment)
	} else {
		println(line, comment)
	}
}
