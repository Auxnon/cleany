package main

import (
	"os"
	"path/filepath"
)

// absPath returns the absolute version of path. On all platforms it resolves
// relative paths against the current working directory.
func absPath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Clean(filepath.Join(wd, path)), nil
}
