package dependency

import "github.com/spf13/cobra"

func AssembleCommandHierarchy() *cobra.Command {
	var root = &cobra.Command{
		Use:   "dependency",
		Short: "Manage dependencies for the project within the current directory.",
	}

	var query = &cobra.Command{
		Use:     "query",
		Short:   "Query the registries for the given search term.",
		Aliases: []string{"q"},
		Run: func(cmd *cobra.Command, args []string) {
			println("TODO")
		},
	}
	query.Flags().StringP("term", "t", "", "Term to query by")
	query.MarkFlagRequired("term")
	root.AddCommand(query)

	var sync = &cobra.Command{
		Use:   "sync",
		Short: "Sync dependencies declared in the project configuration with dependencies on the local filesystem.",
		Run: func(cmd *cobra.Command, args []string) {
			println("TODO")
		},
	}
	root.AddCommand(sync)
	return root
}
