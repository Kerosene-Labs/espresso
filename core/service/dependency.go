package service

import (
	"fmt"
	"sync"

	"github.com/fatih/color"
	"hlafaille.xyz/espresso/v0/core/context"
	"hlafaille.xyz/espresso/v0/core/dependency"
	"hlafaille.xyz/espresso/v0/core/project"
	"hlafaille.xyz/espresso/v0/core/util"
)

// SyncDependencies is a service function to iterate over each dependency and save it within the user's espresso cached packages
func SyncDependencies() {
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
}
