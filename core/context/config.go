package context

import (
	"os"

	"kerosenelabs.com/espresso/core/config"
)

// getConfig reads the config from the filesystem and returns a pointer to a ProjectConfig, or an error
func getConfig() (*config.ProjectConfig, error) {
	// get the config path for this context
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	// read the file
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// unmarshal into a cfg struct
	cfg, err := config.UnmarshalConfig(string(file))
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
