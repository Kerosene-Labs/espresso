package cli

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"hlafaille.xyz/espresso/v0/core/config"
	"hlafaille.xyz/espresso/v0/core/context"
	"hlafaille.xyz/espresso/v0/core/dependency"
	"hlafaille.xyz/espresso/v0/core/project"
	"hlafaille.xyz/espresso/v0/core/registry"
	"hlafaille.xyz/espresso/v0/core/service"
	"hlafaille.xyz/espresso/v0/core/toolchain"
	"hlafaille.xyz/espresso/v0/core/util"
)

func GetCleanCommand() *cobra.Command {
	var root = &cobra.Command{
		Use:     "clean",
		Short:   "Clean the build context",
		Aliases: []string{"c"},
		Run:     service.CleanWorkspace,
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
			// get the environment context
			ctx, err := context.GetEnvironmentContext()
			if err != nil {
				util.ErrorQuit("An error occurred while getting the environment context: %s", err)
			}

			// get our project context
			prjCtx, err := ctx.GetProjectContext()
			if err != nil {
				util.ErrorQuit("An error occurred while getting the project context: %s", err)
			}

			color.Cyan("-- Building '%s', please ensure you are compliant with all dependency licenses\n", prjCtx.Cfg.Name)

			// discover source files
			files, err := project.DiscoverSourceFiles(prjCtx.Cfg)
			if err != nil {
				util.ErrorQuit(fmt.Sprintf("An error occurred while discovering source files: %s\n", err))
			}

			// run the compiler on each source file
			var wg sync.WaitGroup
			for _, value := range files {
				wg.Add(1)
				go func(f *project.SourceFile) {
					defer wg.Done()
					color.Cyan("-- Compiling: " + f.Path)
					err := toolchain.CompileSourceFile(prjCtx.Cfg, &value)
					if err != nil {
						util.ErrorQuit("An error occurred while compiling a source file: %s\n", err)
					}
				}(&value)
			}
			wg.Wait()

			// package the project
			color.Cyan("-- Packaging distributable")
			err = toolchain.PackageClasses(prjCtx.Cfg)
			if err != nil {
				util.ErrorQuit(fmt.Sprintf("An error occurred while packaging the classes: %s\n", err))
			}
			color.Blue("Finished packaging distributable")

			// iterate over each dependency, resolve it and copy it
			distPath, err := toolchain.GetDistPath(prjCtx.Cfg)
			if err != nil {
				util.ErrorQuit(fmt.Sprintf("Unable to get dist path: %s", err))
			}
			os.MkdirAll(*distPath+"/libs", 0755)
			var depCopyWg sync.WaitGroup
			color.Cyan("-- Copying packages to distributable")
			for _, dep := range prjCtx.Cfg.Dependencies {
				depCopyWg.Add(1)
				go func() {
					defer depCopyWg.Done()
					color.Cyan("-- Beginning copy of '%s:%s' to distributable", dep.Group, dep.Name)
					resolved, err := dependency.ResolveDependency(&dep, &prjCtx.Cfg.Registries)
					if err != nil {
						util.ErrorQuit(fmt.Sprintf("Unable to resolve dependency: %s", err))
					}

					// calculate the should-be location of this jar locally
					espressoPath, err := util.GetEspressoDirectoryPath()
					if err != nil {
						util.ErrorQuit(fmt.Sprintf("Unable to get the espresso home: %s", espressoPath))
					}
					signature := registry.CalculatePackageSignature(resolved.Package, resolved.PackageVersion)
					cachedPackageHome := espressoPath + "/cachedPackages" + signature + ".jar"

					// copy the file
					util.CopyFile(cachedPackageHome, *distPath+"/libs")

					color.Blue("Copied '%s:%s' to distributable", dep.Group, dep.Name)
				}()
			}
			depCopyWg.Wait()
			color.Green("Done!")
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
			cfgExists, err := config.DoesExist()
			if err != nil {
				fmt.Println("Error occurred while ensuring a config doesn't already exist")
				panic(err)
			} else if *cfgExists {
				fmt.Println("Config already exists")
				os.Exit(1)
			}

			fmt.Printf("Creating '%s'\n", name)

			// create a base config
			cfg := config.ProjectConfig{
				Name: name,
				Version: config.Version{
					Major:  0,
					Minor:  1,
					Patch:  0,
					Hotfix: nil,
				},
				BasePackage: basePackage,
				Toolchain: config.Toolchain{
					Path: *javaHome,
				},
				Registries:   []config.Registry{{Name: "espresso-registry", Url: "https://github.com/Kerosene-Labs/espresso-registry/archive/refs/heads/main.zip"}},
				Dependencies: []config.Dependency{},
			}

			// write some example code
			println("Creating base package, writing example code")
			project.WriteExampleCode(&cfg)

			// persist the config
			println("Persisting project configuration")
			config.Persist(&cfg)

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
			cfg, err := config.GetConfig()
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
			// get the environment context
			ctx, err := context.GetEnvironmentContext()
			if err != nil {
				util.ErrorQuit("An error occurred while getting the environment context: %s", err)
			}

			// get our project context
			prjCtx, err := ctx.GetProjectContext()
			if err != nil {
				util.ErrorQuit("An error occurred while getting the project context: %s", err)
			}

			// iterate over each registry, invalidate it
			for _, reg := range prjCtx.Cfg.Registries {
				fmt.Printf("Invalidating cache: %s\n", reg.Url)
				err = registry.InvalidateRegistryCache(&reg)
				if err != nil {
					util.ErrorQuit(fmt.Sprintf("An error occurred while invalidaing cache(s): %s\n", err))
				}
			}

			// iterate over each registry, download the zip
			for _, reg := range prjCtx.Cfg.Registries {
				fmt.Printf("Downloading archive: %s\n", reg.Url)
				err = registry.CacheRegistry(&reg)
				if err != nil {
					util.ErrorQuit(fmt.Sprintf("An error occurred while downloading the registry archive: %s\n", err))
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
			// get the environment context
			ctx, err := context.GetEnvironmentContext()
			if err != nil {
				util.ErrorQuit("An error occurred while getting the environment context: %s", err)
			}

			// get our project context
			prjCtx, err := ctx.GetProjectContext()
			if err != nil {
				util.ErrorQuit("An error occurred while getting the project context: %s", err)
			}

			// iterate over the dependencies
			var wg sync.WaitGroup
			for _, dep := range prjCtx.Cfg.Dependencies {
				wg.Add(1)
				go func() {
					defer wg.Done()
					displayStr := fmt.Sprintf("%s:%s:%s", dep.Group, dep.Name, project.GetVersionAsString(&dep.Version))
					color.Cyan("[%s] Resolving", displayStr)
					rdep, err := dependency.ResolveDependency(&dep, &prjCtx.Cfg.Registries)
					if err != nil {
						util.ErrorQuit(fmt.Sprintf("[%s] An error occurred while resolving dependencies: %s\n", displayStr, err))
					}

					// cache the resolved dependency
					err = dependency.CacheResolvedDependency(rdep)
					if err != nil {
						util.ErrorQuit(fmt.Sprintf("[%s] An error occurred while caching the resolved dependency: %s\n", displayStr, err))
					}
					color.Green("[%s] Cached", displayStr)
				}()
			}
			wg.Wait()
		},
	}
	root.AddCommand(sync)

	return root
}
