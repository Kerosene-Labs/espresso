package context

import (
	"os"

	"kerosenelabs.com/espresso/core/util"
)

// getConfigPath gets the absolute path to the config file
func getConfigPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if util.IsDebugMode() {
		path := wd + "/ESPRESSO_DEBUG" + "/espresso.yml"
		return path, nil
	} else {
		path := wd + "/espresso.yml"
		return path, nil
	}
}

// DoesExist gets if the current context's project config exists
func DoesExist(ctx *EnvironmentContext) (bool, error) {
	prjCtx, err := ctx.GetProjectContext()
	if err != nil {
		return false, err
	}
	_, err = os.Stat(*prjCtx.CfgPath)
	resp := err == nil || !os.IsNotExist(err)
	return resp, nil
}

// Persist persists the configuration under the EnvironmentContext to the filesystem
func Persist(ctx *EnvironmentContext) error {
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
