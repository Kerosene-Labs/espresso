package context

import (
	"os"

	"hlafaille.xyz/espresso/v0/core/config"
)

// singular reference to the environment context for this instance
var envCtx *EnvironmentContext = nil

// ProjectContext provides context for the running instance's current project (if any)
type ProjectContext struct {
	Cfg *config.ProjectConfig
}

// EnvironmentContext provides context for the given runtime environment. For example,
// it contains information about the current project (if any), current toolchain, etc.
type EnvironmentContext struct {
	Wd     string
	prjCtx *ProjectContext
}

// GetProjectContext returns a pointer to the cached project context, or loads it incase it doesn't exist.
func (eCtx EnvironmentContext) GetProjectContext() (*ProjectContext, error) {
	if eCtx.prjCtx != nil {
		return eCtx.prjCtx, nil
	}
	// TODO build our loading the project context
	return nil, nil
}

// GetEnvironmentContext returns a pointer to context for the current environment.
func GetEnvironmentContext() (*EnvironmentContext, error) {
	if envCtx != nil {
		return envCtx, nil
	}

	// get the working dir
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	envCtx = &EnvironmentContext{
		Wd: wd,
	}

	return envCtx, nil
}
