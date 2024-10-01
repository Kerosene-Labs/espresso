// Copyright (c) 2024 Kerosene Labs
// This file is part of Espresso, which is licensed under the MIT License.
// See the LICENSE file for details.

package context

import (
	"os"

	"kerosenelabs.com/espresso/core/config"
)

// singular reference to the environment context for this instance
var envCtx *EnvironmentContext = nil

// ProjectContext provides context for the running instance's current project (if any)
type ProjectContext struct {
	Cfg     *config.ProjectConfig
	CfgPath *string
}

// EnvironmentContext provides context for the given runtime environment. For example,
// it contains information about the current project (if any), current toolchain, etc.
type EnvironmentContext struct {
	Wd     *string
	prjCtx *ProjectContext
}

// GetProjectContext returns a pointer to the cached project context, or loads it incase it doesn't exist.
func (eCtx EnvironmentContext) GetProjectContext() (*ProjectContext, error) {
	if eCtx.prjCtx != nil {
		return eCtx.prjCtx, nil
	}

	// get our config path
	cfgPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	// get our config
	cfg, err := getConfig()
	if err != nil {
		return nil, err
	}
	return &ProjectContext{Cfg: cfg, CfgPath: &cfgPath}, nil
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
		Wd: &wd,
	}

	return envCtx, nil
}
