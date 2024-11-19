package cmd

import (
	"flag"
	"fmt"
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
	// flagDebug := flag.Bool("debug", false, "enable debug information")
	// flagLevel := flag.String("log-level", "error", "log level")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		return fmt.Errorf("no file specified")
	}

	// logger.SetLevel(logger.LevelFromString(*flagLevel))

	// b := builder.NewBuilder2()
	// err := b.Build(build.Options{
	// 	InputFilePath:  args[0],
	// 	OutputFilePath: "out",
	// 	NumWorkers:     runtime.NumCPU(),
	// 	Debug:          *flagDebug,
	// })

	// if err != nil {
	// 	errors.PrettyPrint(err)
	// }

	return nil
}
