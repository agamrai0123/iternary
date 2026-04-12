# 📋 Quick Deployment Checklist

## Before Deployment ✓

- [ ] Code pushed to GitHub (`main` branch)
- [ ] `render.yaml` in root directory (✅ Created)
- [ ] `.env.production` in `itinerary-backend/` directory
- [ ] All tests passing locally
- [ ] Docker builds successfully

## Step 1: Create Render Account (2 mins)

1. Go to https://render.com
2. Click "Sign up"
3. Use your GitHub account (easiest)
4. Authorize Render to access your repos

## Step 2: Connect GitHub (2 mins)

1. In Render Dashboard: Click "New +" → "Web Service"
2. Select "GitHub"
3. Search for your `agamrai0123/iternary` repository
4. Click "Connect"

## Step 3: Configure Deployment (3 mins)

Fill in deployment settings:

```
Name: itinerary-backend
Environment: Go
Build Command: cd itinerary-backend && go mod download && go build -o itinerary-backend .
Start Command: cd itinerary-backend && ./itinerary-backend
```

Click "Create Web Service"

## Step 4: Add Environment Variables (3 mins)

Render will ask for environment variables. Add:

```
GIN_MODE = release
PORT = 8080
```

(Database variables are auto-created by Render)

## Step 5: Wait for Deployment (5-10 mins)

✅ Watch the logs in Render Dashboard
✅ You'll get a `.onrender.com` URL
✅ Service auto-starts

## Step 6: Test Your App (2 mins)

```bash
# Replace with your actual Render URL
curl https://itinerary-backend-xxxx.onrender.com/api/health

# Should return: {"status":"ok"}
```

## Total Time: ~20 minutes! ⏱️

---

## Automatic Deployments

Every time you push to `main`:

```bash
cd itinerary-backend
git add .
git commit -m "Update feature"
git push origin main
```

✅ GitHub notifies Render
✅ Render pulls latest code
✅ Render builds & deploys
✅ Your app updates automatically! 🚀

---

## Troubleshooting

### "Build failed" error?

1. Check logs in Render Dashboard → Logs tab
2. Verify locally: `cd itinerary-backend && go build`
3. Check `go.mod` and `go.sum` are committed

### "Application crashed" error?

1. Database URL malformed - check auto-generated DATABASE_URL
2. Port binding failed - Render uses PORT env var (already set to 8080)
3. Missing environment variables - add in Render Settings

### "503 Service Unavailable"?

1. App might still be starting
2. Wait 30-60 seconds
3. Check logs for crashes

### Want to rebuild without pushing code?

1. Render Dashboard → Select service
2. Click "Manual Deploy" button
3. Render redeploys instantly

---

## Stop/Pause Deployment

Render Dashboard → Service Settings → Suspend Service

This saves free tier hours but prevents access. Resume anytime.

---

## Monitor Your App

View real-time:
- 📊 Metrics: CPU, Memory, Network
- 📝 Logs: All application output
- ⚠️ Alerts: Deployment failures
- 🔄 Activity: Recent deployments

All in Render Dashboard!

---

## Upgrade Path (If Needed)

**Free tier limit reached?** Upgrade anytime:

- **Pro Web Service**: $12/month (500+ hours/month)
- **Pro PostgreSQL**: $15/month (more storage/connections)
- **Pro Redis**: $5/month (more memory)

Or just create new free service (takes 2 mins)

---

## Next Steps

1. ✅ Run locally with `./run-local.sh` or `run-local.ps1`
2. ✅ Push code to GitHub
3. ✅ Follow deployment steps above
4. ✅ Share your `.onrender.com` URL
5. ✅ Updates auto-deploy on every push!

**Questions?** See DEPLOYMENT_GUIDE.md for detailed info.
