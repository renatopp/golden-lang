package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/renatopp/golden/internal/backend"
	"github.com/renatopp/golden/internal/backend/golang"
	"github.com/renatopp/golden/internal/backend/javascript"
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
	flagTargets := flag.String("targets", "go", "output backends, separated by comma")
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
		opts.OnDependencyGraphReady.Subscribe(printDependencyGraph)
		opts.OnTypeCheckReady.Subscribe(printTypedAst)
	}

	if flagTargets != nil {
		opts.OutputTargets = []backend.Backend{}
		registered := map[string]bool{}
		for _, o := range strings.Split(*flagTargets, ",") {
			if registered[o] {
				continue
			}
			registered[o] = true
			switch o {
			case "go":
				opts.OutputTargets = append(opts.OutputTargets, golang.NewBackend())
			case "js":
				opts.OutputTargets = append(opts.OutputTargets, javascript.NewBackend())
			default:
				return fmt.Errorf("unknown backend: %s", o)
			}
		}
	}

	if flagWorkingDir != nil {
		abs, _ := fs.GetAbsolutePath(*flagWorkingDir)
		opts.WorkingDir = abs
	}

	b := builder.NewBuilder(opts)
	res, err := b.Build()
	if err != nil {
		errors.PrettyPrint(err)
		return nil
	}
	fmt.Println("Build completed in", res.Elapsed)

	if err := b.Run(); err != nil {
		errors.PrettyPrint(err)
		return nil
	}

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
