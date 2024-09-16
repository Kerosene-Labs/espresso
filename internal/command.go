package internal

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

// GetProjectCommand returns the pre-built Cobra Command 'project'
func GetProjectCommand() *cobra.Command {
	var root = &cobra.Command{
		Use:   "project",
		Short: "Manage a project within the current directory.",
	}

	var build = &cobra.Command{
		Use:     "build",
		Short:   "Build the project.",
		Aliases: []string{"b"},
		Run: func(cmd *cobra.Command, args []string) {
			// get the config
			cfg, err := GetConfig()
			if err != nil {
				fmt.Printf("An error occurred while reading the config: %s\n", err)
			}

			// discover source files
			files, err := DiscoverSourceFiles(cfg)
			if err != nil {
				fmt.Printf("An error occurred while discovering source files: %s\n", err)
			}
			fmt.Printf("Discovered %d source file(s)\n", len(files))

			// run the compiler on each source file
			var wg sync.WaitGroup
			for _, value := range files {
				wg.Add(1)
				go func(f *SourceFile) {
					defer wg.Done()
					println("Compiling: " + f.Path)
					CompileSourceFile(cfg, &value)
				}(&value)
			}
			wg.Wait()

			// package the project
			println("Done")
		},
	}
	root.AddCommand(build)

	var init = &cobra.Command{
		Use:     "init",
		Short:   "Initialize a new project.",
		Aliases: []string{"i"},
		Run: func(cmd *cobra.Command, args []string) {
			var name, _ = cmd.Flags().GetString("name")
			var basePackage, _ = cmd.Flags().GetString("package")

			// ensure JAVA_HOME is set
			javaHome, err := GetJavaHome()
			if err != nil {
				fmt.Println("JAVA_HOME is not set, do you have Java installed?")
			}

			// ensure a proejct doesn't already exist
			cfgExists, err := ConfigExists()
			if err != nil {
				fmt.Println("Error occurred while ensuring a config doesn't already exist")
				panic(err)
			} else if *cfgExists {
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
				Toolchain: Toolchain{
					Path: *javaHome,
				},
				Dependencies: Dependencies{
					Repositories: []Registry{{Url: "https://github.com/Kerosene-Labs/espresso-registry/archive/refs/heads/main.zip"}},
				},
			}

			// write some example code
			println("Creating base package, writing example code")
			WriteExampleCode(&cfg)

			// persist the config
			println("Persisting project configuration")
			PersistConfig(&cfg)

			println("Done.")
		},
	}
	init.Flags().StringP("name", "n", "", "Name of the project")
	init.Flags().StringP("package", "p", "org.example.myapp", "Base package of the application")
	init.MarkFlagRequired("name")
	root.AddCommand(init)

	return root
}

// GetDependencyCommand returns the pre build Cobra Command 'dependency'
func GetDependencyCommand() *cobra.Command {
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

	var add = &cobra.Command{
		Use:   "add",
		Short: "Add a dependency to the project configuration.",
		Run: func(cmd *cobra.Command, args []string) {
			println("TODO")
		},
	}
	root.AddCommand(add)
	return root
}
