package cmd

import (
	"fmt"

	"github.com/renatopp/golden/internal"
	"github.com/renatopp/golden/internal/logger"
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

func (c *Run) Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no file specified")
	}

	logger.SetLevel(logger.TraceLevel)

	builder := internal.NewBuilder()
	err := builder.Build(internal.BuildOptions{
		InputFilePath:  args[0],
		OutputFilePath: "out",
		NumWorkers:     4, //runtime.NumCPU(),
	})

	if err != nil {
		println("Err!", err.Error())
	}

	return nil
}
