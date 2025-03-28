package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/renatopp/golden/internal/backend/golang"
	"github.com/renatopp/golden/internal/backend/javascript"
	"github.com/renatopp/golden/internal/backend/interpreter"
	"github.com/renatopp/golden/internal/builder"
	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/compiler/env"

	// "github.com/renatopp/golden/internal/compiler/env"
	"github.com/renatopp/golden/internal/helpers/debug"
	"github.com/renatopp/golden/internal/helpers/errors"
	"github.com/renatopp/golden/internal/helpers/fs"
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
	flagWorkingDir := flag.String("working-dir", ".", "working directory")
	flagTarget := flag.String("target", "eval", "output backend")
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
		opts.OnOptimizationReady.Subscribe(debug.PrettyPrintAst)
	}

	if flagTarget != nil {
		switch *flagTarget {
		case "go":
			opts.OutputTarget = golang.NewBackend()
		case "js":
			opts.OutputTarget = javascript.NewBackend()
		case "eval":
			opts.OutputTarget = interpreter.NewBackend()
		default:
			return fmt.Errorf("unknown backend: %s", *flagTarget)
		}
	}

	if flagWorkingDir != nil {
		abs, _ := fs.GetAbsolutePath(*flagWorkingDir)
		opts.WorkingDir = abs
	}

	b := builder.NewBuilder(opts)
	res, err := b.Run()
	if err != nil {
		errors.PrettyPrint(err)
		return nil
	}
	fmt.Println("Run completed in", res.Elapsed)

	return nil
}

func printDependencyGraph(order []*builder.File) {
	deps := []string{}
	for _, p := range order {
		deps = append(deps, p.Path)
	}
	names := strings.Join(deps, "\n- ")
	fmt.Printf("Order of dependencies:\n- %s\n", names)
	println()
}

func printTypedAst(mod *builder.File, a *ast.Module, scope *env.Scope) {
	debug.PrettyPrintAst(mod, a)
	debug.PrettyPrintScope(scope)
}
