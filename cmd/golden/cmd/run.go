package cmd

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/renatopp/golden/internal/builder"
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
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		return fmt.Errorf("no file specified")
	}

	logger.SetLevel(logger.ErrorLevel)

	b := builder.NewBuilder()
	err := b.Build(builder.BuildOptions{
		InputFilePath:  args[0],
		OutputFilePath: "out",
		NumWorkers:     runtime.NumCPU(),
		Debug:          *flagDebug,
	})

	if err != nil {
		println("Err!", err.Error())
	}

	return nil
}
