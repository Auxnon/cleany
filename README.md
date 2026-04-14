# cleany

A Go CLI tool that scans a directory tree for empty folders and lets you delete them interactively via a terminal UI powered by [Charm](https://charm.sh/) libraries.

## Features

- Recursively finds **all empty directories** under a given path
- **All directories selected by default** — ready to delete in one key press
- Terminal checklist UI built with [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Space** — toggle selection of the current item
- **a** — select / deselect all
- **↑ / ↓** (or **k / j**) — move cursor
- **d** or **Delete** — delete selected directories (with confirmation prompt)
- **q** — quit without deleting
- Works on **Windows** and **Linux**

## Usage

```
# Scan the current directory
cleany

# Scan a specific directory
cleany /path/to/scan
```

## Install

```
go install github.com/Auxnon/cleany@latest
```

Or build from source:

```
git clone https://github.com/Auxnon/cleany
cd cleany
go build -o cleany .
```