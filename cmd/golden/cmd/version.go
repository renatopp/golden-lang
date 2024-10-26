package cmd

import "github.com/renatopp/golden"

type Version struct{}

func (v *Version) Name() string { return "version" }

func (v *Version) Description() string {
	return "Prints the version of Golden"
}

func (v *Version) Help() string {
	return v.Description()
}

func (v *Version) Run(args []string) int {
	println(golden.Version)
	return 0
}
