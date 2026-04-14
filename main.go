package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	recursive := flag.Bool("r", false, "recursively search for empty directories")
	flag.Parse()

	// Determine root directory.
	root := "."
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	// Resolve to absolute path for nicer display.
	abs, err := absPath(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cleany: %v\n", err)
		os.Exit(1)
	}
	root = abs

	fmt.Printf("Scanning %s …\n", root)

	dirs, err := findEmptyDirs(root, *recursive)
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
