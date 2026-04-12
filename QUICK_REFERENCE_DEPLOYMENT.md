# ⚡ QUICK REFERENCE - Deployment in 4 Commands

## 1️⃣ Test Locally

**Windows:**
```bash
.\run-local.ps1
```

**Linux/Mac:**
```bash
./run-local.sh
```

✅ Result: App at **http://localhost:8080**

---

## 2️⃣ Push to GitHub

**Windows:**
```bash
.\push-to-github.bat
```

**Linux/Mac:**
```bash
./push-to-github.sh
```

✅ Result: Files committed & pushed to GitHub

---

## 3️⃣ Deploy on Render

Go to: **https://render.com**

1. Sign up with GitHub
2. Click "New Web Service"
3. Select your repository
4. Click "Create"

✅ Result: Live in 5-10 minutes

---

## 4️⃣ Share Your App

Your app is now at: `https://itinerary-backend.onrender.com` (or similar)

**Share the URL!** 🎉

---

## 📊 Cost

**$0/month** ✨

---

## 🔄 Future Updates

```bash
git push origin main
```

✅ Renders automatic deploy (no extra steps!)

---

## 🆘 If Something Goes Wrong

Check logs:
1. Render Dashboard → Select service
2. View Logs tab
3. Look for red errors

Most common:
- ❌ `go build` fails → Check go.mod locally
- ❌ Database error → Auto-created by Render
- ❌ Port error → Already configured

See **DEPLOYMENT_GUIDE.md** for detailed fixes.

---

## 📚 Full Guides

- **START_HERE_DEPLOYMENT.md** - Overview
- **VISUAL_DEPLOYMENT_GUIDE.md** - Diagrams
- **DEPLOYMENT_CHECKLIST.md** - Detailed steps
- **DEPLOYMENT_GUIDE.md** - Troubleshooting

---

## ✨ You're Done!

Everything is set up. Just follow the 4 commands above! 🚀
