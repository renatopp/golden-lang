package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	node, err := internal.Parse(tokens)
	if err != nil {
		return fmt.Errorf("parsing file:\n%v", err)
	}

	node.Traverse(func(n *internal.Node, depth int) bool {
		fmt.Printf("%s%s\n", strings.Repeat("  ", depth), n)
		return true
	})

	return nil
}
