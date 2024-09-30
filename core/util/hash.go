package util

import (
	"crypto/sha256"
	"fmt"
)

// GetFileChecksum gets the SHA-256 checksum of a file at a given directory.
func GetChecksum(content string) (string, error) {
	hash := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%x", hash), nil
}
