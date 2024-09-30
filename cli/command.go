package cli

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"hlafaille.xyz/espresso/v0/core/dependency"
	"hlafaille.xyz/espresso/v0/core/project"
	"hlafaille.xyz/espresso/v0/core/registry"
	"hlafaille.xyz/espresso/v0/core/toolchain"
	"hlafaille.xyz/espresso/v0/core/util"
)

func GetCleanCommand() *cobra.Command {
	var root = &cobra.Command{
		Use:     "clean",
		Short:   "Clean the build context",
		Aliases: []string{"c"},
		Run: func(cmd *cobra.Command, args []string) {
			// get the config
			cfg, err := project.GetConfig()
			if err != nil {
				fmt.Printf("An error occurred while reading the config: %s\n", err)
			}

			// get the build dir
			buildPath, err := toolchain.GetBuildPath(cfg)
			if err != nil {
				fmt.Printf("An error occurred while getting the build path: %s\n", err)
				return
			}

			// get the dist dir
			distPath, err := toolchain.GetDistPath(cfg)
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
			cfg, err := project.GetConfig()
			if err != nil {
				panic(fmt.Sprintf("An error occurred while reading the config: %s\n", err))
			}
			fmt.Printf("Building '%s', please ensure you are compliant with all dependency licenses\n", cfg.Name)

			// discover source files
			files, err := project.DiscoverSourceFiles(cfg)
			if err != nil {
				panic(fmt.Sprintf("An error occurred while discovering source files: %s\n", err))
			}
			fmt.Printf("Discovered %d source file(s)\n", len(files))

			// run the compiler on each source file
			var wg sync.WaitGroup
			for _, value := range files {
				wg.Add(1)
				go func(f *project.SourceFile) {
					defer wg.Done()
					println("Compiling: " + f.Path)
					toolchain.CompileSourceFile(cfg, &value)
				}(&value)
			}
			wg.Wait()

			// package the project
			println("Packaging")
			err = toolchain.PackageClasses(cfg)
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
			javaHome, err := util.GetJavaHome()
			if err != nil {
				fmt.Println("JAVA_HOME is not set, do you have Java installed?")
			}

			// ensure a proejct doesn't already exist
			cfgExists, err := project.ConfigExists()
			if err != nil {
				fmt.Println("Error occurred while ensuring a config doesn't already exist")
				panic(err)
			} else if *cfgExists {
				fmt.Println("Config already exists")
				os.Exit(1)
			}

			fmt.Printf("Creating '%s'\n", name)

			// create a base config
			cfg := project.ProjectConfig{
				Name: name,
				Version: project.Version{
					Major:  0,
					Minor:  1,
					Patch:  0,
					Hotfix: nil,
				},
				BasePackage: basePackage,
				Toolchain: project.Toolchain{
					Path: *javaHome,
				},
				Registries:   []project.Registry{{Name: "espresso-registry", Url: "https://github.com/Kerosene-Labs/espresso-registry/archive/refs/heads/main.zip"}},
				Dependencies: []project.Dependency{},
			}

			// write some example code
			println("Creating base package, writing example code")
			project.WriteExampleCode(&cfg)

			// persist the config
			println("Persisting project configuration")
			project.PersistConfig(&cfg)

			println("Done.")
		},
	}
	root.Flags().StringP("name", "n", "", "Name of the project")
	root.Flags().StringP("package", "p", "org.example.myapp", "Base package of the application")
	root.MarkFlagRequired("name")
	return root
}

// GetRegistryCommand returns the pre build Cobra Command 'dependency'
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

			// get the config
			cfg, err := project.GetConfig()
			if err != nil {
				panic(fmt.Sprintf("An error occurred while reading the config: %s\n", err))
			}

			// iterate over each registry, get its packages
			var filteredPkgs []registry.Package = []registry.Package{}
			for _, reg := range cfg.Registries {
				color.Blue("Checking '%s'", reg.Name)
				regPkgs, err := registry.GetRegistryPackages(reg)
				if err != nil {
					panic(fmt.Sprintf("An error occurred while fetching packages from the '%s' registry cache: %s", reg.Name, err))
				}

				// filter by name
				for _, pkg := range regPkgs {
					if term == "*" ||
						strings.Contains(strings.ToLower(pkg.Name), strings.ToLower(term)) ||
						strings.Contains(strings.ToLower(pkg.Description), strings.ToLower(term)) {
						filteredPkgs = append(filteredPkgs, pkg)
					}
				}
			}

			// print out our packages
			color.Cyan("Found %v package(s):", len(filteredPkgs))
			data := [][]string{}
			for _, filtered := range filteredPkgs {
				data = append(data, []string{
					filtered.Group,
					filtered.Name,
					filtered.Versions[len(filtered.Versions)-1].Number,
				})
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Group", "Package", "Latest Version"})

			for _, v := range data {
				table.Append(v)
			}
			table.Render()
		},
	}
	query.Flags().StringP("term", "t", "", "Term to query by")
	query.MarkFlagRequired("term")
	root.AddCommand(query)

	var pull = &cobra.Command{
		Use:   "invalidate",
		Short: "Invalidate and recache the declared registries.",
		Run: func(cmd *cobra.Command, args []string) {
			// get the config
			cfg, err := project.GetConfig()
			if err != nil {
				ErrorQuit(fmt.Sprintf("An error occurred while reading the config: %s\n", err))
			}

			// iterate over each registry, invalidate it
			for _, reg := range cfg.Registries {
				fmt.Printf("Invalidating cache: %s\n", reg.Url)
				err = registry.InvalidateRegistryCache(&reg)
				if err != nil {
					ErrorQuit(fmt.Sprintf("An error occurred while invalidaing cache(s): %s\n", err))
				}
			}

			// iterate over each registry, download the zip
			for _, reg := range cfg.Registries {
				fmt.Printf("Downloading archive: %s\n", reg.Url)
				err = registry.CacheRegistry(&reg)
				if err != nil {
					ErrorQuit(fmt.Sprintf("An error occurred while downloading the registry archive: %s\n", err))
				}
			}
		},
	}
	root.AddCommand(pull)
	return root
}

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
			// get the config
			cfg, err := project.GetConfig()
			if err != nil {
				ErrorQuit(fmt.Sprintf("An error occurred while reading the config: %s\n", err))
			}

			// iterate over the dependencies
			for _, dep := range cfg.Dependencies {
				color.Cyan("Looking for '%s:%s:%s'...", dep.Group, dep.Name, dep.Version)
				color.Blue("Not found, fetching")
				dependency.FetchDependency()
			}
		},
	}
	root.AddCommand(sync)

	return root
}
