# 🧹 tidy

A minimal, beautiful CLI tool built in Go using Bubble Tea to keep your directories clean and organized.

# Features

- tidy cup: Instantly moves loose files into folders based on their extension.
- tidy undo: Reverses the last cleanup and deletes the empty folders.
- tidy ls: A prettier, interactive alternative to ls built for the terminal.
- Filters: Move exactly what you want, or skip what you don’t.

# Installation
1. Clone the repository
```
git clone https://github.com/its-me-abhishek/tidy.git
cd tidy

2. Initialize & build
go mod init tidy
go mod tidy
go build -o tidy main.go

3. Global access (optional)
sudo mv tidy /usr/local/bin/
```

# Usage
## Basic commands
- Command	Action
  - ./tidy cup	Organizes all files in the current directory
  - ./tidy undo	Restores files to their original state
  - ./tidy ls	Opens the interactive Bubble Tea file explorer

- Advanced filtering
  - Target specific file types or exclude them:
    - Only specific types
      ```
       ./tidy cup --ext "jpg png gif"
      ```
    - Everything except specific types
      ```
      ./tidy cup --skip "txt md"
      ```
      
# Development & Testing

Helper scripts are included to safely test the tool (useful in WSL or Linux):

- sample.sh: Generates a messy directory with various file types (.txt, .md, .html, .obscure)
- cleanup.sh: Hard reset script that moves files back to root and removes folders
- Keybindings for tidy ls
  - j / Down → Move selection down
  - k / Up → Move selection up
  - q / Ctrl+C → Quit explorer
