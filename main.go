package main

import (
	"github.com/spf13/cobra"
	"hlafaille.xyz/espresso/v0/dependency"
	"hlafaille.xyz/espresso/v0/project"
)

func main() {
	var root = &cobra.Command{
		Use:   "espresso",
		Short: "A modern Java build tool.",
	}

	// project commands
	root.AddCommand(project.AssembleCommandHierarchy())
	root.AddCommand(dependency.AssembleCommandHierarchy())

	// execute
	root.Execute()
}
