package service

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"hlafaille.xyz/espresso/v0/core/context"
	"hlafaille.xyz/espresso/v0/core/toolchain"
	"hlafaille.xyz/espresso/v0/core/util"
)

func CleanWorkspace(cmd *cobra.Command, args []string) {
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

	// get the build dir
	buildPath, err := toolchain.GetBuildPath(prjCtx.Cfg)
	if err != nil {
		util.ErrorQuit("An error occurred while getting the build path: %s", err)
	}

	// get the dist dir
	distPath, err := toolchain.GetDistPath(prjCtx.Cfg)
	if err != nil {
		util.ErrorQuit("An error occurred while getting the build path: %s", err)
	}

	// remove the build dir
	err = os.RemoveAll(*buildPath)
	if err != nil {
		util.ErrorQuit("An error occurred while deleting the build path: %s", err)
	}

	// remove the dist dir
	err = os.RemoveAll(*distPath)
	if err != nil {
		util.ErrorQuit("An error occurred while deleting the dist path: %s", err)
	}

	color.Green("Cleaned")
}
