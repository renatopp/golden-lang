package main

import (
	"fmt"
	"os"

	"github.com/renatopp/golden/cmd/golden/cmd"
)

type Command interface {
	Name() string
	Description() string
	Help() string
	Run(args []string) int
}

var commands = []Command{
	&cmd.Version{},
	&cmd.Build{},
	&cmd.Run{},
}

func main() {
	if len(os.Args) < 2 || os.Args[1] == "help" {
		help()
		os.Exit(0)
	}

	for _, cmd := range commands {
		if cmd.Name() == os.Args[1] {
			os.Exit(cmd.Run(os.Args))
		}
	}

	fmt.Println("Unknown command:", os.Args[1])
	os.Exit(1)
}

func help() {
	if len(os.Args) < 3 {
		helpAll()
		return
	}

	for _, cmd := range commands {
		if cmd.Name() == os.Args[2] {
			fmt.Printf("Usage: golden %s [arguments]\n\n", cmd.Name())
			fmt.Println(cmd.Help())
			return
		}
	}
}

func helpAll() {
	fmt.Println("Usage: golden <command> [arguments]")
	fmt.Println()
	fmt.Println("The commands are:")
	fmt.Println()
	for _, cmd := range commands {
		fmt.Println("    " + padRight(cmd.Name(), 10) + "  " + cmd.Description())
	}
	fmt.Println()
	fmt.Println("Use \"golden help <command>\" for more information about a command.")
	fmt.Println()
}

func padRight(str string, length int) string {
	for len(str) < length {
		str = str + " "
	}
	return str
}