package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Determine root directory.
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	// Resolve to absolute path for nicer display.
	abs, err := absPath(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cleany: %v\n", err)
		os.Exit(1)
	}
	root = abs

	fmt.Printf("Scanning %s …\n", root)

	dirs, err := findEmptyDirs(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cleany: %v\n", err)
		os.Exit(1)
	}

	m := newModel(dirs)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "cleany: %v\n", err)
		os.Exit(1)
	}
}
