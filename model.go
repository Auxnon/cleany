package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ---------- styles ----------

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	cursorStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("214"))
	dimStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	confirmStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("208"))
)

// ---------- messages ----------

type deleteDoneMsg struct{ deleted int }
type errMsg struct{ err error }

// ---------- view state ----------

type viewState int

const (
	stateList    viewState = iota // main checklist
	stateConfirm                  // "really delete?" prompt
	stateDone                     // finished – quit
)

// ---------- model ----------

type model struct {
	dirs     []string // all empty dirs found
	selected []bool   // parallel slice: true = will be deleted
	cursor   int      // list cursor
	offset   int      // first visible item index (for scrolling)
	state    viewState
	deleted  int    // count of successfully deleted dirs
	lastErr  string // last error message
	width    int
	height   int
}

func newModel(dirs []string) model {
	sel := make([]bool, len(dirs))
	for i := range sel {
		sel[i] = true // select all by default
	}
	return model{
		dirs:     dirs,
		selected: sel,
	}
}

// ---------- init ----------

func (m model) Init() tea.Cmd {
	return nil
}

// ---------- update ----------

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case deleteDoneMsg:
		m.deleted = msg.deleted
		m.state = stateDone
		return m, tea.Quit

	case errMsg:
		m.lastErr = msg.err.Error()
		m.state = stateList

	case tea.KeyMsg:
		switch m.state {

		case stateList:
			return m.handleListKey(msg)

		case stateConfirm:
			return m.handleConfirmKey(msg)
		}
	}

	return m, nil
}

func (m model) handleListKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "Q":
		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
			m.clampOffset()
		}

	case "down", "j":
		if m.cursor < len(m.dirs)-1 {
			m.cursor++
			m.clampOffset()
		}

	case " ":
		if len(m.dirs) > 0 {
			m.selected[m.cursor] = !m.selected[m.cursor]
		}

	case "a", "A":
		// Toggle all: if any is unselected, select all; otherwise deselect all.
		anyUnselected := false
		for _, s := range m.selected {
			if !s {
				anyUnselected = true
				break
			}
		}
		for i := range m.selected {
			m.selected[i] = anyUnselected
		}

	case "d", "D", "delete":
		// Only proceed if at least one dir is selected.
		anySelected := false
		for _, s := range m.selected {
			if s {
				anySelected = true
				break
			}
		}
		if anySelected {
			m.state = stateConfirm
			m.lastErr = ""
		}
	}

	return m, nil
}

// clampOffset adjusts m.offset so the cursor is always within the visible window.
func (m *model) clampOffset() {
	visible := m.visibleRows()
	if visible <= 0 {
		return
	}
	if m.cursor < m.offset {
		m.offset = m.cursor
	} else if m.cursor >= m.offset+visible {
		m.offset = m.cursor - visible + 1
	}
}

// visibleRows returns the number of list rows that fit in the terminal.
// headerLines accounts for: title line, blank line, stats/hints line, blank line.
// footerLines accounts for: blank line, error message line.
const headerLines = 4
const footerLines = 2

func (m model) visibleRows() int {
	if m.height == 0 {
		return 20 // sensible default before first WindowSizeMsg
	}
	h := m.height - headerLines
	if m.lastErr != "" {
		h -= footerLines
	}
	if h < 1 {
		return 1
	}
	return h
}

func (m model) handleConfirmKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y", "enter":
		return m, deleteSelected(m.dirs, m.selected)
	case "n", "N", "esc", "ctrl+c":
		m.state = stateList
	}
	return m, nil
}

// deleteSelected returns a Cmd that deletes selected directories.
func deleteSelected(dirs []string, selected []bool) tea.Cmd {
	return func() tea.Msg {
		deleted := 0
		var firstErr error
		for i, path := range dirs {
			if !selected[i] {
				continue
			}
			if err := os.Remove(path); err != nil {
				if firstErr == nil {
					firstErr = err
				}
				continue
			}
			deleted++
		}
		if firstErr != nil {
			return errMsg{firstErr}
		}
		return deleteDoneMsg{deleted}
	}
}

// ---------- view ----------

func (m model) View() string {
	switch m.state {
	case stateConfirm:
		return m.confirmView()
	case stateDone:
		return m.doneView()
	default:
		return m.listView()
	}
}

func (m model) listView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("cleany — empty directory scanner") + "\n\n")

	if len(m.dirs) == 0 {
		b.WriteString(dimStyle.Render("No empty directories found.") + "\n")
	} else {
		// Count selected
		numSelected := 0
		for _, s := range m.selected {
			if s {
				numSelected++
			}
		}

		// Scroll indicator suffix
		scrollInfo := ""
		visible := m.visibleRows()
		total := len(m.dirs)
		if total > visible {
			scrollText := fmt.Sprintf("[%d-%d/%d]", m.offset+1, min(m.offset+visible, total), total)
			scrollInfo = fmt.Sprintf("  %s", dimStyle.Render(scrollText))
		}

		b.WriteString(fmt.Sprintf("%s  %s%s\n\n",
			dimStyle.Render(fmt.Sprintf("%d/%d selected", numSelected, total)),
			helpStyle.Render("[↑/↓] move  [space] toggle  [a] all  [d/del] delete  [q] quit"),
			scrollInfo,
		))

		end := m.offset + visible
		if end > total {
			end = total
		}
		for i := m.offset; i < end; i++ {
			dir := m.dirs[i]
			cursor := "  "
			if i == m.cursor {
				cursor = cursorStyle.Render("▶ ")
			}

			checkbox := "[ ]"
			line := dir
			if m.selected[i] {
				checkbox = selectedStyle.Render("[✓]")
				line = selectedStyle.Render(dir)
			} else {
				checkbox = dimStyle.Render("[ ]")
				line = dimStyle.Render(dir)
			}

			b.WriteString(fmt.Sprintf("%s%s %s\n", cursor, checkbox, line))
		}
	}

	if m.lastErr != "" {
		b.WriteString("\n" + errorStyle.Render("Error: "+m.lastErr) + "\n")
	}

	return b.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (m model) confirmView() string {
	numSelected := 0
	for _, s := range m.selected {
		if s {
			numSelected++
		}
	}
	return fmt.Sprintf(
		"%s\n\n%s\n",
		confirmStyle.Render(fmt.Sprintf("Delete %d director%s? This cannot be undone.", numSelected, plural(numSelected))),
		helpStyle.Render("[y/enter] confirm   [n/esc] cancel"),
	)
}

func (m model) doneView() string {
	return fmt.Sprintf("%s\n",
		titleStyle.Render(fmt.Sprintf("Done! Deleted %d director%s.", m.deleted, plural(m.deleted))),
	)
}

func plural(n int) string {
	if n == 1 {
		return "y"
	}
	return "ies"
}
