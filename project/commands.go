package project

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func AssembleCommandHierarchy() *cobra.Command {
	var root = &cobra.Command{
		Use:   "project",
		Short: "Manage a project within the current working directory.",
	}

	var init = &cobra.Command{
		Use:   "init",
		Short: "Initialize a new project.",
		Run: func(cmd *cobra.Command, args []string) {
			var name, _ = cmd.Flags().GetString("name")
			var basePackage, _ = cmd.Flags().GetString("package")

			// ensure a proejct doesn't already exist
			if ConfigExists() {
				fmt.Println("Config already exists")
				os.Exit(1)
			}

			fmt.Printf("Creating '%s'\n", name)

			// create a base config
			cfg := ProjectConfig{
				Name: name,
				Version: Version{
					Major:  0,
					Minor:  1,
					Patch:  0,
					Hotfix: nil,
				},
				BasePackage: basePackage,
			}

			// persist the config
			PersistConfig(&cfg)

			// write some example code
			println("Writing example code")
			WriteExampleCode(&cfg)

			println("Done.")
		},
	}
	init.Flags().StringP("name", "n", "", "Name of the project")
	init.Flags().StringP("package", "p", "org.example.myapp", "Base package of the application")
	init.MarkFlagRequired("name")
	root.AddCommand(init)

	return root
}
