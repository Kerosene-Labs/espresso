package service

import "fmt"

var CommitSha string
var Version string

// PrintVersion prints the version and commit sha to the standard output
func PrintVersion() {
	fmt.Printf("ver='%s' commit='%s'\n", Version, CommitSha)
}
