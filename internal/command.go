package internal

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

func GetCleanCommand() *cobra.Command {
	var root = &cobra.Command{
		Use:     "clean",
		Short:   "Clean the build context",
		Aliases: []string{"c"},
		Run: func(cmd *cobra.Command, args []string) {
			// get the config
			cfg, err := GetConfig()
			if err != nil {
				fmt.Printf("An error occurred while reading the config: %s\n", err)
			}

			// get the build dir
			buildPath, err := GetBuildPath(cfg)
			if err != nil {
				fmt.Printf("An error occurred while getting the build path: %s\n", err)
				return
			}

			// get the dist dir
			distPath, err := GetDistPath(cfg)
			if err != nil {
				fmt.Printf("An error occurred while getting the build path: %s\n", err)
				return
			}

			// remove the build dir
			err = os.RemoveAll(*buildPath)
			if err != nil {
				fmt.Printf("An error occurred while deleting the build path: %s\n", err)
				return
			}

			// remove the dist dir
			err = os.RemoveAll(*distPath)
			if err != nil {
				fmt.Printf("An error occurred while deleting the dist path: %s\n", err)
				return
			}
		},
	}
	return root
}

// GetProjectCommand returns the pre-built Cobra Command 'project'
func GetBuildCommand() *cobra.Command {
	var root = &cobra.Command{
		Use:     "build",
		Short:   "Build the project, outputting a distributable.",
		Aliases: []string{"b"},
		Run: func(cmd *cobra.Command, args []string) {
			// get the config
			cfg, err := GetConfig()
			if err != nil {
				fmt.Printf("An error occurred while reading the config: %s\n", err)
			}
			fmt.Printf("Building '%s', please ensure you are compliant with all dependency licenses\n", cfg.Name)

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
			println("Packaging")
			err = PackageClasses(cfg)
			if err != nil {
				fmt.Printf("An error occurred while packaging the classes: %s\n", err)
				return
			}
			println("Done")
		},
	}
	return root
}

func GetInitCommand() *cobra.Command {
	var root = &cobra.Command{
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
					Registries: []Registry{{Name: "espresso-registry", Url: "https://github.com/Kerosene-Labs/espresso-registry/archive/refs/heads/main.zip"}},
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
	root.Flags().StringP("name", "n", "", "Name of the project")
	root.Flags().StringP("package", "p", "org.example.myapp", "Base package of the application")
	root.MarkFlagRequired("name")
	return root
}

// GetDependencyCommand returns the pre build Cobra Command 'dependency'
func GetDependencyCommand() *cobra.Command {
	var root = &cobra.Command{
		Use:   "registry",
		Short: "Manage package registries for the project within the current directory.",
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
		Use:   "pull",
		Short: "Sync the dependencies declared in the project configuration with dependencies on the local filesystem.",
		Run: func(cmd *cobra.Command, args []string) {
			// get the config
			cfg, err := GetConfig()
			if err != nil {
				fmt.Printf("An error occurred while reading the config: %s\n", err)
			}

			// iterate over each registry, invalidate it
			for _, reg := range cfg.Dependencies.Registries {
				fmt.Printf("Invalidating cache: %s\n", reg.Url)
				err = InvalidateRegistryCache(&reg)
				if err != nil {
					fmt.Printf("An error occurred while invalidaing cache(s): %s\n", err)
				}
			}

			// iterate over each registry, download the zip
			for _, reg := range cfg.Dependencies.Registries {
				fmt.Printf("Downloading archive: %s\n", reg.Url)
				err = CacheRegistry(&reg)
				if err != nil {
					fmt.Printf("An error occurred while downloading the registry archive: %s\n", err)
				}
			}
		},
	}
	root.AddCommand(sync)
	return root
}
