package cmd

import (
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

func (c *Run) Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no project specified")
	}

	return nil
}
