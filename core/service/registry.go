// Copyright (c) 2024 Kerosene Labs
// This file is part of Espresso, which is licensed under the MIT License.
// See the LICENSE file for details.

package service

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"kerosenelabs.com/espresso/core/context/project"
	"kerosenelabs.com/espresso/core/registry"
	"kerosenelabs.com/espresso/core/util"
)

// QueryRegistries is a service function for querying all registries declared within a project
func QueryRegistries(term string) {
	// get our project context
	projectContext, err := project.GetProjectContext()
	if err != nil {
		util.ErrorQuit("An error occurred while getting the project context: %s", err)
	}

	// iterate over each registry, get its packages
	var filteredPkgs []registry.Package = []registry.Package{}
	for _, reg := range projectContext.Config.Registries {
		color.Blue("Checking '%s'", reg.Name)
		regPkgs, err := registry.GetRegistryPackages(reg)
		if err != nil {
			util.ErrorQuit("An error occurred while fetching packages from the '%s' registry cache: %s", reg.Name, err)
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
	// get our project context
	projectContext, err := project.GetProjectContext()
	if err != nil {
		util.ErrorQuit("An error occurred while getting the project context: %s", err)
	}

	// iterate over each registry, invalidate it
	var invalWg sync.WaitGroup
	for _, reg := range projectContext.Config.Registries {
		invalWg.Add(1)
		go func() {
			defer invalWg.Done()
			color.Cyan("[%s] Invalidating cache", reg.Name)
			err = registry.InvalidateRegistryCache(reg)
			if err != nil {
				util.ErrorQuit("An error occurred while invalidaing cache(s): %s", err)
			}
			color.Blue("[%s] Invalidated", reg.Name)
		}()
	}
	invalWg.Wait()

	// iterate over each registry, download the zip
	var dlWg sync.WaitGroup
	for _, reg := range projectContext.Config.Registries {
		dlWg.Add(1)
		go func() {
			defer dlWg.Done()
			color.Cyan("[%s] Downloading archive", reg.Name)
			err = registry.CacheRegistry(reg)
			if err != nil {
				util.ErrorQuit(fmt.Sprintf("An error occurred while downloading the registry archive: %s\n", err))
			}
			color.Blue("[%s] Finished caching", reg.Name)
		}()
	}
	dlWg.Wait()
}
