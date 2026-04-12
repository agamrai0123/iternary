@echo off
REM Push deployment configuration to GitHub
REM Windows Batch Script

echo.
echo 🚀 Preparing deployment configuration for GitHub...
echo.

REM Check if git is installed
where git >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo ❌ Git not found! Install from https://git-scm.com
    pause
    exit /b 1
)

echo ✅ Git found
echo.

REM Add all deployment files
echo 📦 Adding deployment files...
git add render.yaml
git add .github/workflows/deploy.yml
git add DEPLOYMENT_GUIDE.md
git add DEPLOYMENT_CHECKLIST.md
git add START_HERE_DEPLOYMENT.md
git add run-local.ps1
git add run-local.sh
git add itinerary-backend\.env.production

if %ERRORLEVEL% NEQ 0 (
    echo ❌ Failed to add files!
    pause
    exit /b 1
)

echo ✅ Files added
echo.

REM Show what will be committed
echo 📋 Files ready to commit:
git diff --cached --name-only

echo.
echo 📝 Committing to git...
git commit -m "Add free deployment configuration (Render.com + GitHub Actions)"

if %ERRORLEVEL% NEQ 0 (
    echo ⚠️  Commit failed (maybe nothing changed?)
    pause
    exit /b 1
)

echo ✅ Committed successfully
echo.

REM Ask to push
set /p push="Push to GitHub now? (y/n): "
if /i "%push%"=="y" (
    echo 📤 Pushing to GitHub...
    git push origin main
    
    if %ERRORLEVEL% EQU 0 (
        echo.
        echo ✅ ✅ ✅ DEPLOYMENT CONFIGURATION PUSHED ✅ ✅ ✅
        echo.
        echo 🎉 Next steps:
        echo   1. Visit https://render.com
        echo   2. Sign up with GitHub
        echo   3. Click "New Web Service"
        echo   4. Select your repository
        echo   5. Deploy!
        echo.
    ) else (
        echo ❌ Push failed!
    )
) else (
    echo ⏭️  Skipping push. Commit files locally when ready.
)

pause
