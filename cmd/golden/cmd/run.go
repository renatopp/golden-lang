package cmd

type Run struct{}

func (c *Run) Name() string {
	return "Run"
}

func (c *Run) Description() string {
	return "Runs the project"
}

func (c *Run) Help() string {
	return "Runs the project"
}

func (c *Run) Run(args []string) int {
	return 0
}
