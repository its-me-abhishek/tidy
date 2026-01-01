# 🧹 tidy

<img width=600 src="./demo.gif"/>

A minimal, beautiful CLI tool built in Go using Bubble Tea to keep your directories clean and organized. Moves all files inside a directory, to separate directories, based on the file extension!

# Features

- `tidy cup`: Instantly moves loose files into folders based on their extension. (Cup is short for cleanup!)
  - **Preview Mode**: Use `-p` or `--preview` to see exactly what will happen before touching a single file.
- `tidy undo`: Reverses the last cleanup and restores your directory to its original state.
  - **How it works**: Each `tidy cup` run generates a hidden `.tidy_history` file inside that same directory to track moved files.
  - **Note**: Deleting `.tidy_history` makes the reorganisation using cleanup permanent and prevents further undos for that session.
- `tidy ls`: A prettier, interactive alternative to ls built for the terminal.
- Filters: Move exactly what you want, or skip what you don’t.

# Installation

## The Quickest Way (Go Install)
If you have Go installed, this will compile and install tidy to your $GOPATH/bin automatically.

```
go install github.com/its-me-abhishek/tidy@latest
```

To make tidy work just by typing its name, you need to tell WSL/Linux where Go puts its binaries. Run these two commands:
```
# Add Go bin to your current session
export PATH=$PATH:$(go env GOPATH)/bin

# Make it permanent for every time you open WSL
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
```

## Pre-built Binaries

Grab the latest executable for Windows, macOS, or Linux from the Releases Page.

## Building from source
1. Clone the repository
```
git clone https://github.com/its-me-abhishek/tidy.git
cd tidy
```
2. Initialize & build
```
go mod init tidy
go mod tidy
go build -o tidy main.go
```
3. Global access (optional)
```
sudo mv tidy /usr/local/bin/
```

# Usage
## Basic commands
- Command	Action
  - tidy cup: Organizes all files in the current directory
  - tidy undo: Restores files to their original state
  - tidy ls: Opens the interactive Bubble Tea file explorer
  - tidy help: Lists all commands available for usage

- Advanced filtering
  - Target specific file types or exclude them:
    - Only specific types
      ```
       tidy cup --ext "jpg png gif"
      ```
    - Everything except specific types
      ```
      tidy cup --skip "txt md"
      ```
      
# Development & Testing

Helper scripts are included to safely test the tool (useful in WSL or Linux):

- sample.sh: Generates a messy directory with various file types (.txt, .md, .html, .obscure)
- cleanup.sh: Hard reset script that moves files back to root and removes folders
- Keybindings for tidy ls
  - j / Down → Move selection down
  - k / Up → Move selection up
  - q / Ctrl+C → Quit explorer
