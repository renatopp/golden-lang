package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/renatopp/golden/internal"
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
		return fmt.Errorf("no project specified")
	}

	project, err := filepath.Abs(args[0])
	if err != nil {
		return fmt.Errorf("error getting absolute path: %v", err)
	}

	pkg, err := internal.ReadPackage(project)
	if err != nil {
		return fmt.Errorf("error reading package: %v", err)
	}

	println(pkg.Debug())
	_ = pkg
	return nil
}
