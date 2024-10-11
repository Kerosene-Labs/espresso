package project

import "sync"

// ProjectContext provides context for this instances current project (if any)
type ProjectContext struct {
	Config     ProjectConfig
	ConfigPath string
	SourcePath string
}

var projectContext *ProjectContext = nil
var once sync.Once

// GetProjectContext lazily loads the ProjectContext singleton
func GetProjectContext() (*ProjectContext, error) {
	var err error
	once.Do(func() {
		// get the config path
		configPath, e := GetConfigPath()
		if e != nil {
			err = e
			return
		}

		// get our config
		config, e := readConfigFromFileSystem()
		if e != nil {
			err = e
			return
		}

		// set the project context
		projectContext = &ProjectContext{Config: config, ConfigPath: configPath}
	})
	if err != nil {
		return nil, err
	}

	return projectContext, nil
}
