package cmd

import (
	"flag"
	"fmt"

	"github.com/renatopp/golden/internal/builder"
	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/helpers/debug"
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
		OnTokensReady: func(module *builder.Module, tokens []*lang.Token) {
			if *flagDebug {
				debug.PrettyPrintTokens(module, tokens)
			}
		},
		OnAstReady: func(module *builder.Module, root *ast.Module) {
			if *flagDebug {
				debug.PrettyPrintAst(module, root)
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
