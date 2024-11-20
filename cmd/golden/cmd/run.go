package cmd

import (
	"flag"
	"fmt"

	"github.com/renatopp/golden/internal/builder"
	"github.com/renatopp/golden/internal/helpers/debug"
	"github.com/renatopp/golden/internal/helpers/errors"
	"github.com/renatopp/golden/internal/helpers/logger"
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

	opts := builder.NewBuildOptions(args[0])
	if *flagDebug {
		opts.OnTokensReady.Subscribe(debug.PrettyPrintTokens)
		opts.OnAstReady.Subscribe(debug.PrettyPrintAst)
	}

	b := builder.NewBuilder(opts)
	res, err := b.Build()
	if err != nil {
		errors.PrettyPrint(err)
		return nil
	}

	fmt.Println("Build completed in", res.Elapsed)
	return nil
}
