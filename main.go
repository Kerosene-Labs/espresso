package main

import (
	"github.com/spf13/cobra"
	"kerosenelabs.com/espresso/cli"
)

var CommitSha string
var Version string

func main() {
	var root = &cobra.Command{
		Use:   "espresso",
		Short: "A modern Java build tool.",
	}

	root.AddCommand(cli.GetVersionCommand())
	root.AddCommand(cli.GetCleanCommand())
	root.AddCommand(cli.GetBuildCommand())
	root.AddCommand(cli.GetInitCommand())
	root.AddCommand(cli.GetRegistryCommand())
	root.AddCommand(cli.GetDependencyCommand())

	// execute
	root.Execute()
}
