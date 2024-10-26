package cmd

import "github.com/renatopp/golden"

type Version struct{}

func (c *Version) Name() string { return "version" }

func (c *Version) Description() string {
	return "Prints the version of Golden"
}

func (c *Version) Help() string {
	return c.Description()
}

func (c *Version) Run(args []string) error {
	println(golden.Version)
	return nil
}
