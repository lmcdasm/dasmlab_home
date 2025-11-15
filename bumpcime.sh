#!/bin/bash
# Script to toggle "bump" at end of README.md and commit/push

set -e

README="README.md"

# Check if README.md exists
if [ ! -f "$README" ]; then
    echo "ERROR: $README not found"
    exit 1
fi

# Check if "bump" is already at the end
if tail -1 "$README" | grep -q "^bump$"; then
    echo "Removing 'bump' from end of $README"
    # Remove the last line if it's just "bump"
    sed -i '$ { /^bump$/d; }' "$README"
else
    echo "Adding 'bump' to end of $README"
    echo "bump" >> "$README"
fi

# Git operations
echo "Staging all changes..."
git add -A .

echo "Committing with message 'CI bump'..."
git commit -m "CI bump"

echo "Pushing to origin main..."
git push origin main

echo "Done!"

