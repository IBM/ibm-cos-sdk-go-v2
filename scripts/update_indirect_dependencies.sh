#!/bin/bash

# Script to find all go.mod files in the workspace and update indirect dependencies
# Usage:
#   ./scripts/update_indirect_dependencies.sh              - List all go.mod files
#   ./scripts/update_indirect_dependencies.sh --update     - Update all indirect dependencies in each go.mod
#   ./scripts/update_indirect_dependencies.sh --tidy       - Run go mod tidy on all go.mod files

set -e  # Exit on error

UPDATE_MODE=false
TIDY_MODE=false

# Parse command line arguments
if [[ "$1" == "--update" ]]; then
    UPDATE_MODE=true
elif [[ "$1" == "--tidy" ]]; then
    TIDY_MODE=true
fi

echo "Searching for go.mod files in: $(pwd)"
echo "================================"
echo ""

# Find all go.mod files recursively
GO_MOD_FILES=$(find . -name "go.mod" -type f | sort)

if [ -z "$GO_MOD_FILES" ]; then
    echo "No go.mod files found!"
    exit 1
fi

# Display all go.mod files
echo "$GO_MOD_FILES"

echo ""
echo "================================"
TOTAL_COUNT=$(echo "$GO_MOD_FILES" | wc -l | tr -d ' ')
echo "Total go.mod files found: $TOTAL_COUNT"

# Update dependencies if requested
if [ "$UPDATE_MODE" = true ]; then
    echo ""
    echo "================================"
    echo "Updating indirect dependencies..."
    echo "================================"
    echo ""
    
    SUCCESS_COUNT=0
    FAIL_COUNT=0
    
    while IFS= read -r go_mod_file; do
        # Get the directory containing go.mod
        dir=$(dirname "$go_mod_file")
        
        echo "Processing: $go_mod_file"
        echo "  Directory: $dir"
        
        # Change to the directory and update dependencies
        if (cd "$dir" && go get -u all && go mod tidy); then
            echo "  ✓ Successfully updated"
            ((SUCCESS_COUNT++))
        else
            echo "  ✗ Failed to update"
            ((FAIL_COUNT++))
        fi
        echo ""
    done <<< "$GO_MOD_FILES"
    
    echo "================================"
    echo "Update Summary:"
    echo "  Successful: $SUCCESS_COUNT"
    echo "  Failed: $FAIL_COUNT"
    echo "  Total: $TOTAL_COUNT"
    echo "================================"
elif [ "$TIDY_MODE" = true ]; then
    echo ""
    echo "================================"
    echo "Running go mod tidy..."
    echo "================================"
    echo ""
    
    SUCCESS_COUNT=0
    FAIL_COUNT=0
    
    while IFS= read -r go_mod_file; do
        # Get the directory containing go.mod
        dir=$(dirname "$go_mod_file")
        
        echo "Processing: $go_mod_file"
        echo "  Directory: $dir"
        
        # Change to the directory and run go mod tidy
        if (cd "$dir" && go mod tidy); then
            echo "  ✓ Successfully tidied"
            ((SUCCESS_COUNT++))
        else
            echo "  ✗ Failed to tidy"
            ((FAIL_COUNT++))
        fi
        echo ""
    done <<< "$GO_MOD_FILES"
    
    echo "================================"
    echo "Tidy Summary:"
    echo "  Successful: $SUCCESS_COUNT"
    echo "  Failed: $FAIL_COUNT"
    echo "  Total: $TOTAL_COUNT"
    echo "================================"
else
    echo ""
    echo "Available commands:"
    echo "  ./scripts/update_indirect_dependencies.sh --update    Update all indirect dependencies"
    echo "  ./scripts/update_indirect_dependencies.sh --tidy      Run go mod tidy on all modules"
fi

# Made with Bob
