// Copyright (c) 2024 Kerosene Labs
// This file is part of Espresso, which is licensed under the MIT License.
// See the LICENSE file for details.

package util

import "os/user"

// GetEspressoDirectoryPath gets the path to the user's .espresso directory.
func GetEspressoDirectoryPath() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.HomeDir + "/.espresso", nil
}
