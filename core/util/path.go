package util

import "os/user"

// GetEspressoDirectoryPath gets the path to the user's .espresso directory.
func GetEspressoDirectoryPath() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.HomeDir, nil
}
