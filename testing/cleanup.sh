#!/bin/bash

# Define the extensions we expect tidy to have created as folders
folders=("txt" "md" "html" "obscure" "jpg" "csv" "zip" "misc")

for dir in "${folders[@]}"; do
    if [ -d "$dir" ]; then
        # Move files back to root before deleting folder (optional safety)
        mv "$dir"/* . 2>/dev/null
        rmdir "$dir"
        echo "🗑️ Removed folder: $dir"
    fi
done

# Also remove the hidden history file if it exists
rm -f .tidy_history

echo "✨ Environment reset."