package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/renatopp/golden/internal"
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
	tokens, err := internal.Lex(file)
	if err != nil {
		return fmt.Errorf("lexing file:\n%v", err)
	}

	for _, t := range tokens {
		fmt.Printf("- %s: %q\n", t.Kind, t.Literal)
	}
	println("\n")

	println("## Parser Output:\n")
	module, err := internal.Parse(tokens)
	if err != nil {
		return fmt.Errorf("parsing file:\n%v", err)
	}

	for _, imp := range module.Imports {
		if imp.Alias != "" {
			fmt.Printf("import %s as %s\n", imp.Path, imp.Alias)
		} else {
			fmt.Printf("import %s\n", imp.Path)
		}
	}

	for _, decl := range module.Types {
		println(decl.String())
	}
	for _, decl := range module.Functions {
		println(decl.String())
	}
	for _, decl := range module.Variables {
		println(decl.String())
	}

	println(module.Temp.String())
	println("\n")

	println("## Analyzer Output:\n")

	scope := internal.NewScope()
	scope.Set("Void", createPrimitive("Void"))
	scope.Set("Bool", createPrimitive("Bool"))
	scope.Set("Int", createPrimitive("Int"))
	scope.Set("Float", createPrimitive("Float"))
	scope.Set("String", createPrimitive("String"))

	module.Scope = scope.New()
	err = internal.Analyze(module, module.Scope)
	if err != nil {
		return fmt.Errorf("analyzing module:\n%v", err)
	}

	for _, decl := range module.Types {
		println(decl.String())
	}
	for _, decl := range module.Functions {
		println(decl.String())
	}
	for _, decl := range module.Variables {
		println(decl.String())
	}

	println(module.Temp.String())

	println()
	println(module.Scope.String())

	return nil
}

func createPrimitive(name string) *internal.Node {
	return internal.NewEmptyNode().WithType(internal.NewPrimitiveType(name))
}
