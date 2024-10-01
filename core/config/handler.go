// Copyright (c) 2024 Kerosene Labs
// This file is part of Espresso, which is licensed under the MIT License.
// See the LICENSE file for details.

package config

import (
	"os"

	"kerosenelabs.com/espresso/core/context"
)

// DoesExist gets if the current context's project config exists
func DoesExist(ctx *context.EnvironmentContext) (bool, error) {
	prjCtx, err := ctx.GetProjectContext()
	if err != nil {
		return false, err
	}
	_, err = os.Stat(*prjCtx.CfgPath)
	resp := err == nil || !os.IsNotExist(err)
	return resp, nil
}

// Persist persists the configuration under the EnvironmentContext to the filesystem
func Persist(ctx *context.EnvironmentContext) error {
	// get the project context
	prjCtx, err := ctx.GetProjectContext()
	if err != nil {
		return err
	}

	marshalData, err := MarshalConfig(prjCtx.Cfg)
	if err != nil {
		return err
	}

	cfgPath, err := prjCtx.CfgPath
	if err != nil {
		return err
	}

	// write the config
	cfgExists, err := DoesExist()
	if err != nil {
		return err
	}

	if cfgExists {
		file, err := os.Open(*cfgPath)
		if err != nil {
			return err
		}
		file.WriteString(*marshalData)
		defer file.Close()

	} else {
		file, err := os.Create(*cfgPath)
		if err != nil {
			return err
		}
		file.WriteString(*marshalData)
		defer file.Close()
	}

	return nil
}
