package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

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

	logger.SetLevel(logger.ErrorLevel)

	builder := internal.NewBuilder()
	err := builder.Build(internal.BuildOptions{
		InputFilePath:  args[0],
		OutputFilePath: "out",
		NumWorkers:     runtime.NumCPU(),
	})

	if err != nil {
		println("Err!", err.Error())
	}

	cmd := exec.Command("./.tools/tcc/win/tcc.exe", "-run", ".out/main.c")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute command: %v, output: %s", err, string(output))
	}

	fmt.Println(string(output))

	return nil
}
