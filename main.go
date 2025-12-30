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
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA")).MarginLeft(4)
)

func printHelp() {
	fmt.Println(titleStyle.Render("🧹 TIDY CLI HELP"))
	fmt.Println("\nUsage: tidy <command> [flags]")
	fmt.Println("\nCommands:")
	fmt.Printf("  %-10s %s\n", "cup", "Clean Up: Organize files into folders by type.")
	fmt.Printf("  %-10s %s\n", "undo", "Revert the last cleanup operation.")
	fmt.Printf("  %-10s %s\n", "ls", "Pretty-print directory contents with interactive UI.")
	fmt.Printf("  %-10s %s\n", "help", "Show this help menu.")

	fmt.Println("\nFlags for 'cup':")
	fmt.Printf("  %-20s %s\n", "--ext \"png jpg\"", "Move ONLY these extensions.")
	fmt.Printf("  %-20s %s\n", "--skip \"mp4 exe\"", "Move everything EXCEPT these extensions.")

	fmt.Println("\nExample:")
	fmt.Println(helpStyle.Render("tidy cup --ext \"pdf docx\""))
}

func runCleanup(targetExts string, skipExts string) {
	files, _ := os.ReadDir(".")
	history, err := os.Create(historyFile)
	if err != nil {
		fmt.Println("Error recording history")
		return
	}
	defer history.Close()

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

		if len(targets) > 0 && !targets[ext] {
			continue
		}
		if skips[ext] {
			continue
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
	extFlag := cupCmd.String("ext", "", "Only move specific extensions")
	skipFlag := cupCmd.String("skip", "", "Skip specific extensions")

	if len(os.Args) < 2 {
		printHelp()
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
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Printf(warnStyle.Render("Unknown command: %s\n\n"), os.Args[1])
		printHelp()
	}
}
