// Copyright (c) 2024 Kerosene Labs
// This file is part of Espresso, which is licensed under the MIT License.
// See the LICENSE file for details.

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
