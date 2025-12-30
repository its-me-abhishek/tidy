# Output configuration
Output demo.gif
Set FontSize 18
Set Width 1200
Set Height 800
Set Padding 40
Set Theme "Catppuccin Macchiato"

# 1. Setup: Create the mess
Type "cd testing"
Enter
Type "./sample.sh"
Enter
Sleep 1s

# 2. Show the mess with ls
Type "ls"
Enter
Sleep 2s

# 3. Run Tidy Cup
Type "tidy cup"
Enter
Sleep 2s

# 4. Show the result
Type "ls"
Enter
Sleep 2s

# 5. Show off the UI (ls mode)
Type "tidy ls"
Enter
Sleep 1s
Down 2
Sleep 500ms
Down 3
Sleep 500ms
Up 1
Sleep 1s
Type "q"
Enter

# 6. Show the Undo feature
Type "tidy undo"
Enter
Sleep 2s

# 7. Final Check
Type "ls"
Enter
Sleep 3s