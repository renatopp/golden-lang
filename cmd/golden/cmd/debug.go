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

	tokens, err := internal.Lex(file)
	if err != nil {
		return fmt.Errorf("error lexing file: %v", err)
	}

	println("## Tokens")
	for _, t := range tokens {
		fmt.Printf("- %s: %q\n", t.Kind, t.Literal)
	}

	return nil
}
