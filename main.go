package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const historyFile = ".tidy_history"

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00ADD8")).MarginLeft(2)
	doneStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#02BA59")).Bold(true).PaddingLeft(2)
	warnStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")).Bold(true).PaddingLeft(2)
	dirStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))
)

func runCleanup(targetExts string, skipExts string) {
	files, _ := os.ReadDir(".")
	history, err := os.Create(historyFile)
	if err != nil {
		fmt.Println("Error recording history")
		return
	}
	defer history.Close()

	// Parse flags into maps for quick lookup
	targets := make(map[string]bool)
	if targetExts != "" {
		for _, e := range strings.Fields(targetExts) {
			targets[strings.TrimPrefix(e, ".")] = true
		}
	}

	skips := make(map[string]bool)
	if skipExts != "" {
		for _, e := range strings.Fields(skipExts) {
			skips[strings.TrimPrefix(e, ".")] = true
		}
	}

	count := 0
	for _, file := range files {
		if file.IsDir() || file.Name() == "main.go" || file.Name() == "tidy" || file.Name() == historyFile {
			continue
		}

		ext := strings.TrimPrefix(filepath.Ext(file.Name()), ".")
		if ext == "" {
			ext = "misc"
		}

		// Logic for Filter Flags
		if len(targets) > 0 && !targets[ext] {
			continue // Skip if not in the target list
		}
		if skips[ext] {
			continue // Skip if explicitly blacklisted
		}

		os.MkdirAll(ext, os.ModePerm)
		newPath := filepath.Join(ext, file.Name())

		if err := os.Rename(file.Name(), newPath); err == nil {
			fmt.Fprintf(history, "%s|%s\n", file.Name(), newPath)
			count++
		}
	}
	fmt.Printf("%s Tidied up %d files.\n", doneStyle.Render("✔"), count)
}

// ... runUndo and Bubble Tea Model/Update/View remain the same as previous version ...

func runUndo() {
	file, err := os.Open(historyFile)
	if err != nil {
		fmt.Println(warnStyle.Render("No history found! Nothing to undo."))
		return
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	count := 0
	dirsToCleanup := make(map[string]bool)
	for _, line := range lines {
		parts := strings.Split(line, "|")
		if len(parts) == 2 {
			if err := os.Rename(parts[1], parts[0]); err == nil {
				dirsToCleanup[filepath.Dir(parts[1])] = true
				count++
			}
		}
	}
	for dir := range dirsToCleanup {
		entries, _ := os.ReadDir(dir)
		if len(entries) == 0 {
			os.Remove(dir)
		}
	}
	os.Remove(historyFile)
	fmt.Printf("%s Undo complete! Restored %d files.\n", doneStyle.Render("↺"), count)
}

type model struct {
	files  []fs.DirEntry
	cursor int
}

func (m model) Init() tea.Cmd { return nil }
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyMsg); ok {
		switch k.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.files)-1 {
				m.cursor++
			}
		}
	}
	return m, nil
}
func (m model) View() string {
	s := titleStyle.Render("TIDY EXPLORER") + "\n\n"
	for i, f := range m.files {
		cursor := "  "
		name := f.Name()
		if f.IsDir() {
			name = dirStyle.Render(name + "/")
		}
		if m.cursor == i {
			cursor = "> "
		}
		s += fmt.Sprintf("%s%s\n", cursor, name)
	}
	return s + "\n (q to quit)"
}

func main() {
	cupCmd := flag.NewFlagSet("cup", flag.ExitOnError)
	extFlag := cupCmd.String("ext", "", "Only move specific extensions (e.g. 'jpg png')")
	skipFlag := cupCmd.String("skip", "", "Skip specific extensions (e.g. 'pdf exe')")

	if len(os.Args) < 2 {
		fmt.Println("Usage: tidy [cup | ls | undo]")
		return
	}

	switch os.Args[1] {
	case "cup":
		cupCmd.Parse(os.Args[2:])
		runCleanup(*extFlag, *skipFlag)
	case "undo":
		runUndo()
	case "ls":
		files, _ := os.ReadDir(".")
		p := tea.NewProgram(model{files: files})
		p.Run()
	default:
		fmt.Println("Unknown command")
	}
}
