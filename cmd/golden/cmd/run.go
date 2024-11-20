package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/renatopp/golden/internal/builder"
	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/core"
	"github.com/renatopp/golden/internal/helpers/errors"
	"github.com/renatopp/golden/internal/helpers/logger"
	"github.com/renatopp/golden/lang"
)

type Run struct{}

func (c *Run) Name() string {
	return "run"
}

func (c *Run) Description() string {
	return "Runs the project"
}

func (c *Run) Help() string {
	return "Runs the project"
}

func (c *Run) Run() error {
	flagDebug := flag.Bool("debug", false, "enable debug information")
	flagLevel := flag.String("log-level", "error", "log level")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		return fmt.Errorf("no file specified")
	}

	logger.SetLevel(logger.LevelFromString(*flagLevel))

	b := builder.NewBuilder(&builder.BuildOptions{
		EntryFilePath: args[0],
		OnTokensReady: func(module *core.Module, tokens []*lang.Token) {
			if *flagDebug {
				c.printTokens(module, tokens)
			}
		},
		OnAstReady: func(module *core.Module, root *ast.Module) {
			if *flagDebug {
				fmt.Printf("AST for file %s:\n", module.Path)
				fmt.Println(root.Id(), root.Token())
				println()
			}
		},
	})

	res, err := b.Build()
	if err != nil {
		errors.PrettyPrint(err)
		return nil
	}

	fmt.Println("Build completed in", res.Elapsed)
	return nil
}

func (c *Run) printTokens(module *core.Module, tokens []*lang.Token) {
	fmt.Printf("Tokens for file %s:\n", module.Path)
	for _, t := range tokens {
		fmt.Printf("- %s('%s')\n", t.Kind, strings.ReplaceAll(t.Literal, "\n", "\\n"))
	}
	println()
}
