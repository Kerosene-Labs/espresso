package config

import "os"

// DoesExist gets if a config exists at the current directory
func DoesExist() (*bool, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	_, err = os.Stat(*configPath)
	resp := err == nil || !os.IsNotExist(err)
	return &resp, nil
}

// GetConfig reads the config from the filesystem and returns a pointer to a ProjectConfig, or an error
func GetConfig() (*ProjectConfig, error) {
	// get the config path for this context
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// read the file
	file, err := os.ReadFile(*path)
	if err != nil {
		return nil, err
	}

	// unmarshal into a cfg struct
	cfg, err := UnmarshalConfig(string(file))
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// Persist persists the given ProjectConfig to the filesystem
func Persist(cfg *ProjectConfig) error {
	marshalData, err := MarshalConfig(cfg)
	if err != nil {
		return err
	}

	cfgPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// write the config
	cfgExists, err := DoesExist()
	if err != nil {
		return err
	}

	if *cfgExists {
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
