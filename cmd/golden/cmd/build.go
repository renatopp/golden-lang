package cmd

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

func (c *Build) Run(args []string) int {
	return 0
}
