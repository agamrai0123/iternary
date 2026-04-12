#!/usr/bin/env powershell
# Auto-Commit and Deploy Script for Windows
# Watches for changes, commits, and pushes to GitHub (which auto-deploys on Render)
# Run this script in PowerShell

param(
    [string]$WatchPath = "d:\Learn\iternary\itinerary-backend",
    [int]$IntervalSeconds = 5
)

$GitPath = "d:\Learn\iternary"
$LastHash = ""

Write-Host "🚀 Auto-Commit & Deploy Monitor Started" -ForegroundColor Cyan
Write-Host "Watching: $WatchPath" -ForegroundColor Cyan
Write-Host "Check interval: ${IntervalSeconds}s" -ForegroundColor Cyan
Write-Host "Auto-deploys to Render on every push!" -ForegroundColor Green
Write-Host ""
Write-Host "Press Ctrl+C to stop" -ForegroundColor Yellow
Write-Host ""

while ($true) {
    try {
        # Get all staged and unstaged changes
        $Status = git -C $GitPath status --porcelain
        
        # Create hash of current changes
        $CurrentHash = ($Status | ConvertTo-String) | Get-Hash -Algorithm MD5
        
        if ($CurrentHash -ne $LastHash -and $Status) {
            Write-Host "[$(Get-Date -Format 'HH:mm:ss')] ✨ Changes detected!" -ForegroundColor Green
            
            # Stage all changes
            Write-Host "  📦 Staging changes..." -ForegroundColor Cyan
            git -C $GitPath add -A
            
            # Get list of changed files
            $ChangedFiles = git -C $GitPath diff --cached --name-only
            Write-Host "  📝 Changed files:" -ForegroundColor Cyan
            foreach ($file in $ChangedFiles) {
                Write-Host "     - $file" -ForegroundColor Gray
            }
            
            # Create commit message
            $CommitMessage = "🔄 Auto-commit: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')`n`nAutomatic commit from file watcher`nFiles changed: $($ChangedFiles.Count)"
            
            # Commit changes
            Write-Host "  📤 Committing..." -ForegroundColor Cyan
            git -C $GitPath commit -m $CommitMessage | Out-Null
            
            if ($LASTEXITCODE -eq 0) {
                Write-Host "  ✅ Committed successfully" -ForegroundColor Green
                
                # Push to GitHub
                Write-Host "  🚀 Pushing to GitHub..." -ForegroundColor Cyan
                $PushOutput = git -C $GitPath push origin main 2>&1
                
                if ($LASTEXITCODE -eq 0) {
                    Write-Host "  ✅ Pushed to GitHub" -ForegroundColor Green
                    Write-Host "  ⏳ Render will auto-deploy in 1-2 minutes..." -ForegroundColor Yellow
                    Write-Host ""
                } else {
                    Write-Host "  ❌ Push failed: $PushOutput" -ForegroundColor Red
                }
            }
            
            $LastHash = $CurrentHash
        }
        
        Start-Sleep -Seconds $IntervalSeconds
    } catch {
        Write-Host "[$(Get-Date -Format 'HH:mm:ss')] Error: $_" -ForegroundColor Red
        Start-Sleep -Seconds $IntervalSeconds
    }
}
