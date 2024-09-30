package service

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"kerosenelabs.com/espresso/core/config"
	"kerosenelabs.com/espresso/core/context"
	"kerosenelabs.com/espresso/core/registry"
	"kerosenelabs.com/espresso/core/util"
)

// QueryRegistries is a service function for querying all registries declared within a project
func QueryRegistries(term string) {
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
}

// InvalidateRegistryCaches is a service function that invalidates and recaches all registries declared within the project.
func InvalidateRegistryCaches() {
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
}
