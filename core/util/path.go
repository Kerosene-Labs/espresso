// Copyright (c) 2024 Kerosene Labs
// This file is part of Espresso, which is licensed under the MIT License.
// See the LICENSE file for details.

package util

import "os/user"

type Path struct {
	Absolute string
}

// DoesExist returns if the path does or does not exist.
func (p Path) DoesExist() (bool, error) {
	exists, err := DoesPathExist(p.Absolute)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// GetEspressoDirectoryPath gets the path to the user's .espresso directory.
func GetEspressoDirectoryPath() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.HomeDir + "/.espresso", nil
}
