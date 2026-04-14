package main

import (
	"os"
	"path/filepath"
)

// isEmptyDir reports whether a directory contains no entries (files or subdirs).
func isEmptyDir(path string) (bool, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}

// findEmptyDirs walks root and returns all paths to empty directories.
// It walks depth-first so that deeply-nested empty dirs are listed before
// their parents; this prevents "directory not empty" errors when callers
// delete in list order.
func findEmptyDirs(root string) ([]string, error) {
	var results []string

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			// Skip directories we cannot read.
			return nil
		}
		if !d.IsDir() {
			return nil
		}
		if path == root {
			return nil
		}

		empty, err := isEmptyDir(path)
		if err != nil {
			return nil
		}
		if empty {
			results = append(results, path)
			// Skip descending into empty dirs (there's nothing inside).
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}
