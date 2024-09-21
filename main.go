package main

import (
	"github.com/spf13/cobra"
	"hlafaille.xyz/espresso/v0/internal"
)

func main() {
	var root = &cobra.Command{
		Use:   "espresso",
		Short: "A modern Java build tool.",
	}

	root.AddCommand(internal.GetCleanCommand())
	root.AddCommand(internal.GetBuildCommand())
	root.AddCommand(internal.GetInitCommand())
	root.AddCommand(internal.GetRegistryCommand())
	root.AddCommand(internal.GetDependencyCommand())

	// execute
	root.Execute()
}
