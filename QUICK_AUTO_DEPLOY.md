# ⚡ Auto-Deploy Quick Start

## 1 Command to Enable Auto-Deploy

**Windows PowerShell:**
```bash
cd d:\Learn\iternary && powershell -ExecutionPolicy Bypass -File auto-commit-deploy.ps1
```

**Linux/macOS:**
```bash
cd /d/Learn/iternary && bash auto-commit-deploy.sh
```

Then save any file and watch it auto-deploy! ✨

---

## What Happens

```
You save a file
    ↓
Script detects change (5s)
    ↓
Auto-commits & pushes
    ↓
GitHub Actions tests it
    ↓
Render deploys (1-2 min)
    ↓
✅ LIVE!
```

---

## Deployment Pipeline

| Component | Status | How It Works |
|-----------|--------|-------------|
| **File Watcher** | Optional | Watches for changes → auto-commits → pushes |
| **GitHub Actions** | ✅ Active | Tests & builds on every push |
| **Render Deploy** | ✅ Active | Auto-deploys on successful GitHub Actions |

---

## Monitor Deployments

1. **GitHub Actions:** https://github.com/agamrai0123/iternary/actions
2. **Render Dashboard:** https://dashboard.render.com
3. **Live App:** https://itinerary-backend-xxxxx.onrender.com

---

## Files Created

| File | Purpose |
|------|---------|
| `auto-commit-deploy.ps1` | Windows file watcher |
| `auto-commit-deploy.sh` | Linux/macOS file watcher |
| `.github/workflows/deploy.yml` | GitHub Actions (already active) |
| `AUTO_DEPLOYMENT_SETUP.md` | Detailed setup guide |

---

## Stop Auto-Commit

Press `Ctrl+C` in the terminal where the script runs.

---

**That's it!** 🚀 Run the script and enjoy automated deployments!
