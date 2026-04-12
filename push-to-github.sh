#!/bin/bash
# Push deployment configuration to GitHub
# Linux/macOS Bash Script

echo ""
echo "🚀 Preparing deployment configuration for GitHub..."
echo ""

# Check if git is installed
if ! command -v git &> /dev/null; then
    echo "❌ Git not found! Install from https://git-scm.com"
    exit 1
fi

echo "✅ Git found"
echo ""

# Add all deployment files
echo "📦 Adding deployment files..."
git add render.yaml
git add .github/workflows/deploy.yml
git add DEPLOYMENT_GUIDE.md
git add DEPLOYMENT_CHECKLIST.md
git add START_HERE_DEPLOYMENT.md
git add run-local.ps1
git add run-local.sh
git add itinerary-backend/.env.production

if [ $? -ne 0 ]; then
    echo "❌ Failed to add files!"
    exit 1
fi

echo "✅ Files added"
echo ""

# Show what will be committed
echo "📋 Files ready to commit:"
git diff --cached --name-only

echo ""
echo "📝 Committing to git..."
git commit -m "Add free deployment configuration (Render.com + GitHub Actions)"

if [ $? -ne 0 ]; then
    echo "⚠️  Commit failed (maybe nothing changed?)"
    exit 1
fi

echo "✅ Committed successfully"
echo ""

# Ask to push
read -p "Push to GitHub now? (y/n): " push
if [ "$push" = "y" ] || [ "$push" = "Y" ]; then
    echo "📤 Pushing to GitHub..."
    git push origin main
    
    if [ $? -eq 0 ]; then
        echo ""
        echo "✅ ✅ ✅ DEPLOYMENT CONFIGURATION PUSHED ✅ ✅ ✅"
        echo ""
        echo "🎉 Next steps:"
        echo "   1. Visit https://render.com"
        echo "   2. Sign up with GitHub"
        echo "   3. Click 'New Web Service'"
        echo "   4. Select your repository"
        echo "   5. Deploy!"
        echo ""
    else
        echo "❌ Push failed!"
    fi
else
    echo "⏭️  Skipping push. Commit files locally when ready."
fi
