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

// findEmptyDirs returns all empty directories found under root.
// When recursive is false only immediate children of root are examined.
// When recursive is true the full directory tree is walked depth-first so
// that deeply-nested empty dirs appear before their parents; this prevents
// "directory not empty" errors when callers delete in list order.
func findEmptyDirs(root string, recursive bool) ([]string, error) {
	if !recursive {
		return findEmptyDirsShallow(root)
	}

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

// findEmptyDirsShallow checks only the immediate children of root and returns
// those that are empty directories.
func findEmptyDirsShallow(root string) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	var results []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		path := filepath.Join(root, e.Name())
		empty, err := isEmptyDir(path)
		if err != nil {
			continue
		}
		if empty {
			results = append(results, path)
		}
	}
	return results, nil
}
