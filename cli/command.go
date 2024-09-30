package cli

import (
	"github.com/spf13/cobra"
	"kerosenelabs.com/espresso/core/service"
)

// GetVersionCommand gets the prepared "version" command for cobra
func GetVersionCommand() *cobra.Command {
	var root = &cobra.Command{
		Use:     "version",
		Short:   "Print the version of Espresso",
		Aliases: []string{"v"},
		Run: func(cmd *cobra.Command, args []string) {
			service.PrintVersion()
		},
	}
	return root
}

// GetCleanCommand gets the prepared "clean" command for cobra
func GetCleanCommand() *cobra.Command {
	var root = &cobra.Command{
		Use:     "clean",
		Short:   "Clean the build context",
		Aliases: []string{"c"},
		Run: func(cmd *cobra.Command, args []string) {
			service.CleanWorkspace()
		},
	}
	return root
}

// GetBuildCommand gets the prepared "build" command for cobra
func GetBuildCommand() *cobra.Command {
	var root = &cobra.Command{
		Use:     "build",
		Short:   "Build the project, outputting a distributable.",
		Aliases: []string{"b"},
		Run: func(cmd *cobra.Command, args []string) {
			service.BuildProject()
		},
	}
	return root
}

// GetInitCommand gets the prepared "init" command for cobra
func GetInitCommand() *cobra.Command {
	var root = &cobra.Command{
		Use:     "init",
		Short:   "Initialize a new project.",
		Aliases: []string{"i"},
		Run: func(cmd *cobra.Command, args []string) {
			var name, _ = cmd.Flags().GetString("name")
			var basePackage, _ = cmd.Flags().GetString("package")

			service.InitializeProject(name, &basePackage)
		},
	}
	root.Flags().StringP("name", "n", "", "Name of the project")
	root.Flags().StringP("package", "p", "org.example.myapp", "Base package of the application")
	root.MarkFlagRequired("name")
	return root
}

// GetRegistryCommand gets the prepared "registry" command for cobra
func GetRegistryCommand() *cobra.Command {
	var root = &cobra.Command{
		Use:   "registry",
		Short: "Manage package registries for the project within the current directory.",
	}

	var query = &cobra.Command{
		Use:     "query",
		Short:   "Query the registries for the given search term.",
		Aliases: []string{"q"},
		Run: func(cmd *cobra.Command, args []string) {
			var term, _ = cmd.Flags().GetString("term")
			service.QueryRegistries(term)
		},
	}
	query.Flags().StringP("term", "t", "", "Term to query by")
	query.MarkFlagRequired("term")
	root.AddCommand(query)

	var pull = &cobra.Command{
		Use:   "invalidate",
		Short: "Invalidate and recache the declared registries.",
		Run: func(cmd *cobra.Command, args []string) {
			service.InvalidateRegistryCaches()
		},
	}
	root.AddCommand(pull)
	return root
}

// GetDependencyCommand gets the prepared "dependency" command for cobra
func GetDependencyCommand() *cobra.Command {
	var root = &cobra.Command{
		Use:   "dependency",
		Short: "Manage dependencies for the project within the current directory.",
	}

	var sync = &cobra.Command{
		Use:     "sync",
		Short:   "Fetch dependencies from the appropriate registries, storing them within their caches for consumption at compile time.",
		Aliases: []string{"s"},
		Run: func(cmd *cobra.Command, args []string) {
			service.SyncDependencies()
		},
	}
	root.AddCommand(sync)

	return root
}
