package cmd

import (
	"flag"
	"fmt"

	"github.com/renatopp/golden/internal/backend/golang"
	"github.com/renatopp/golden/internal/backend/javascript"
	"github.com/renatopp/golden/internal/builder"
	"github.com/renatopp/golden/internal/helpers/debug"
	"github.com/renatopp/golden/internal/helpers/errors"
	"github.com/renatopp/golden/internal/helpers/fs"
	"github.com/renatopp/golden/internal/helpers/logger"
)

type Build struct{}

func (c *Build) Name() string {
	return "build"
}

func (c *Build) Description() string {
	return "Builds the project"
}

func (c *Build) Help() string {
	return "Builds the project"
}

func (c *Build) Run() error {
	flagDebug := flag.Bool("debug", false, "enable debug information")
	flagLevel := flag.String("log-level", "error", "log level")
	flagWorkingDir := flag.String("working-dir", ".", "working directory")
	flagTarget := flag.String("target", "go", "output backend")
	flagOutput := flag.String("output", "", "output file")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		return fmt.Errorf("no file specified")
	}

	logger.SetLevel(logger.LevelFromString(*flagLevel))

	file, _ := fs.GetAbsolutePath(args[0])
	opts := builder.NewBuildOptions(file)
	if *flagDebug {
		opts.OnTokensReady.Subscribe(debug.PrettyPrintTokens)
		opts.OnAstReady.Subscribe(debug.PrettyPrintAst)
		opts.OnDependencyGraphReady.Subscribe(printDependencyGraph)
		opts.OnTypeCheckReady.Subscribe(printTypedAst)
	}

	if flagTarget != nil {
		switch *flagTarget {
		case "go":
			opts.OutputTarget = golang.NewBackend()
		case "js":
			opts.OutputTarget = javascript.NewBackend()
		default:
			return fmt.Errorf("unknown target %s", *flagTarget)
		}
	}

	if flagWorkingDir != nil {
		abs, _ := fs.GetAbsolutePath(*flagWorkingDir)
		opts.WorkingDir = abs
	}

	if flagOutput != nil && *flagOutput != "" {
		opts.OutputFilePath, _ = fs.GetAbsolutePath(*flagOutput)
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
