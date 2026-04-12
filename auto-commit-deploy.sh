#!/bin/bash
# Auto-Commit and Deploy Script for Linux/macOS
# Watches for changes, commits them, and pushes to GitHub (which auto-deploys on Render)

WATCH_PATH="/d/Learn/iternary/itinerary-backend"
GIT_PATH="/d/Learn/iternary"
INTERVAL_SECONDS=5
LAST_HASH=""

echo "🚀 Auto-Commit & Deploy Monitor Started"
echo "Watching: $WATCH_PATH"
echo "Check interval: ${INTERVAL_SECONDS}s"
echo "Auto-deploys to Render on every push!"
echo ""
echo "Press Ctrl+C to stop"
echo ""

while true; do
    # Get all changes
    STATUS=$(cd "$GIT_PATH" && git status --porcelain)
    
    # Create hash of current changes
    CURRENT_HASH=$(echo "$STATUS" | md5sum | cut -d' ' -f1)
    
    if [ "$CURRENT_HASH" != "$LAST_HASH" ] && [ -n "$STATUS" ]; then
        echo "[$(date +'%H:%M:%S')] ✨ Changes detected!"
        
        # Stage all changes
        echo "  📦 Staging changes..."
        cd "$GIT_PATH" && git add -A
        
        # Show changed files
        echo "  📝 Changed files:"
        cd "$GIT_PATH" && git diff --cached --name-only | sed 's/^/     - /'
        
        # Commit changes
        echo "  📤 Committing..."
        COMMIT_MSG="🔄 Auto-commit: $(date +'%Y-%m-%d %H:%M:%S')

Automatic commit from file watcher"
        
        cd "$GIT_PATH" && git commit -m "$COMMIT_MSG" &>/dev/null
        
        if [ $? -eq 0 ]; then
            echo "  ✅ Committed successfully"
            
            # Push to GitHub
            echo "  🚀 Pushing to GitHub..."
            cd "$GIT_PATH" && git push origin main 2>&1 | grep -E "^(To|.*\[|.*\])" || echo "  ✅ Pushed to GitHub"
            
            if [ $? -eq 0 ]; then
                echo "  ✅ Pushed successfully"
                echo "  ⏳ Render will auto-deploy in 1-2 minutes..."
                echo ""
            fi
        fi
        
        LAST_HASH="$CURRENT_HASH"
    fi
    
    sleep $INTERVAL_SECONDS
done
