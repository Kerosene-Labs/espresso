// Copyright (c) 2024 Kerosene Labs
// This file is part of Espresso, which is licensed under the MIT License.
// See the LICENSE file for details.

package service

import "fmt"

var CommitSha string
var Version string

// PrintVersion prints the version and commit sha to the standard output
func PrintVersion() {
	fmt.Printf("ver='%s' commit='%s'\n", Version, CommitSha)
}
