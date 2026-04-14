package main

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestFindEmptyDirs_none(t *testing.T) {
	root := t.TempDir()
	// Create a file inside root so root itself is non-empty.
	f, err := os.CreateTemp(root, "file")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	dirs, err := findEmptyDirs(root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(dirs) != 0 {
		t.Fatalf("expected 0 empty dirs, got %v", dirs)
	}
}

func TestFindEmptyDirs_single(t *testing.T) {
	root := t.TempDir()
	empty := filepath.Join(root, "empty")
	if err := os.Mkdir(empty, 0o755); err != nil {
		t.Fatal(err)
	}

	dirs, err := findEmptyDirs(root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(dirs) != 1 || dirs[0] != empty {
		t.Fatalf("expected [%s], got %v", empty, dirs)
	}
}

func TestFindEmptyDirs_nestedEmpty(t *testing.T) {
	// Layout:
	//   root/
	//     a/       (non-empty – contains b/)
	//       b/     (empty)
	//     c/       (empty)
	root := t.TempDir()
	a := filepath.Join(root, "a")
	b := filepath.Join(a, "b")
	c := filepath.Join(root, "c")
	for _, d := range []string{a, b, c} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	dirs, err := findEmptyDirs(root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sort.Strings(dirs)
	want := []string{b, c}
	sort.Strings(want)

	if len(dirs) != len(want) {
		t.Fatalf("expected %v, got %v", want, dirs)
	}
	for i := range want {
		if dirs[i] != want[i] {
			t.Fatalf("expected %v, got %v", want, dirs)
		}
	}
}

func TestFindEmptyDirs_skipRootItself(t *testing.T) {
	// An empty root should not include itself.
	root := t.TempDir()
	dirs, err := findEmptyDirs(root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, d := range dirs {
		if d == root {
			t.Fatalf("root itself should not appear in results")
		}
	}
}

func TestIsEmptyDir(t *testing.T) {
	root := t.TempDir()

	empty, err := isEmptyDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if !empty {
		t.Fatal("newly created tempdir should be empty")
	}

	f, err := os.CreateTemp(root, "x")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	empty, err = isEmptyDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if empty {
		t.Fatal("dir with a file should not be empty")
	}
}

func TestAbsPath(t *testing.T) {
	abs, err := absPath("/tmp")
	if err != nil {
		t.Fatal(err)
	}
	if abs != "/tmp" {
		t.Fatalf("expected /tmp, got %s", abs)
	}

	wd, _ := os.Getwd()
	abs2, err := absPath(".")
	if err != nil {
		t.Fatal(err)
	}
	if abs2 != wd {
		t.Fatalf("expected %s, got %s", wd, abs2)
	}
}
