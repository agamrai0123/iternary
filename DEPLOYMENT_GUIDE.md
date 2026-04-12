# 🚀 Complete Deployment Guide - Render.com (Free)

## Part 1: Run Locally (Testing)

### Quick Start with Docker Compose
```bash
cd itinerary-backend
docker-compose up --build
```

This starts:
- ✅ PostgreSQL (port 5432)
- ✅ Redis (port 6379)  
- ✅ Go Backend (port 8080)

Access: http://localhost:8080

### Or Run Without Docker
```bash
# Install Go 1.21+
cd itinerary-backend
cp .env.example .env
go mod download
go run main.go
```

---

## Part 2: Deploy Free to Render.com

### Why Render?
- ✅ **Free PostgreSQL database** (5GB storage)
- ✅ **Free Redis cache** (included)
- ✅ **Automatic GitHub deployments** (push to deploy)
- ✅ **Free SSL certificates**
- ✅ **Easy environment variables**
- ✅ No credit card needed for small projects

### Step-by-Step Deployment

#### 1. **Create `render.yaml` in root directory**

This file tells Render how to deploy your app:

```yaml
services:
  - type: web
    name: itinerary-backend
    env: go
    buildCommand: "cd itinerary-backend && go mod download && go build -o itinerary-backend ."
    startCommand: "cd itinerary-backend && ./itinerary-backend"
    envVars:
      - key: GIN_MODE
        value: release
      - key: DB_HOST
        fromService:
          type: postgres
          name: itinerary-db
          envVarKey: PGHOST
      - key: DB_PORT
        value: 5432
      - key: DB_USER
        fromService:
          type: postgres
          name: itinerary-db
          envVarKey: PGUSER
      - key: DB_PASSWORD
        fromService:
          type: postgres
          name: itinerary-db
          envVarKey: PGPASSWORD
      - key: DB_NAME
        fromService:
          type: postgres
          name: itinerary-db
          envVarKey: PGDATABASE
      - key: REDIS_HOST
        value: redis
      - key: REDIS_PORT
        value: 6379
      - key: PORT
        value: 8080

  - type: postgres
    name: itinerary-db
    ipAllowList: [] # Allow all IP addresses

  - type: redis
    name: redis
    ipAllowList: [] # Allow all IP addresses
```

#### 2. **Create `.env.production` in `itinerary-backend/`**

```env
GIN_MODE=release
PORT=8080
DB_HOST=itinerary-db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=<will-be-set-by-render>
DB_NAME=itinerary_production
REDIS_HOST=redis
REDIS_PORT=6379
```

#### 3. **Push to GitHub**

```bash
git add .
git commit -m "Add Render deployment configuration"
git push origin main
```

#### 4. **Deploy on Render.com**

1. Go to https://render.com (Sign up with GitHub)
2. Click **"New +"** → **"Web Service"**
3. Connect your GitHub repository
4. Select the `main` branch
5. Name: `itinerary-backend`
6. Build Command: `cd itinerary-backend && go mod download && go build -o itinerary-backend .`
7. Start Command: `cd itinerary-backend && ./itinerary-backend`
8. Plan: **Free**
9. Click **"Create Web Service"**

Render will automatically:
- ✅ Create PostgreSQL database
- ✅ Create Redis cache
- ✅ Build & deploy your app
- ✅ Give you a free `.onrender.com` URL

#### 5. **View Deployment**

After ~5-10 minutes:
- Visit: `https://itinerary-backend.onrender.com`
- Logs visible in Render dashboard
- Database automatically migrated on first run

---

## Part 3: Automatic Deployments (CI/CD)

Every time you push to `main`, Render automatically:
1. Pulls latest code from GitHub
2. Builds the binary
3. Redeploys to production

**No additional setup needed!** ✨

---

## Part 4: Custom Domain (Optional - Free)

If you want a custom domain like `myitinerary.com`:
- Render gives you CNAME settings
- Register domain on Namecheap (~$8.88/year) or Freenom (free tier)
- Point CNAME to Render

---

## Environment Variables in Render

Visit Dashboard → **Itinerary Settings** → **Environment**

Defaults auto-configured. You can override:
- Database credentials
- Redis settings
- API keys

---

## Troubleshooting

### App crashes on deploy?
1. Check Render logs: Dashboard → Logs tab
2. Verify `go mod download` works locally
3. Check database URL format

### Database not connecting?
1. Render creates DB_HOST automatically
2. Check DATABASE_SETUP.md for migration scripts
3. Run migrations in Render shell:
   ```bash
   psql $DATABASE_URL < migrations/schema.sql
   ```

### Static files 404?
Ensure `static/` and `templates/` directories are in your Go code properly served.

---

## Pricing (Estimate)

- **Total cost: $0/month** (for small projects)
- Free tier includes:
  - 750 free hours/month (~22 hours/day)
  - 512MB free PostgreSQL
  - 256MB free Redis
  - 1 concurrent request

Upgrades only if you exceed free tier (about $7-12/month each).

---

## Quick Reference

| Component | Local | Render |
|-----------|-------|--------|
| Backend | Go 1.21+ | Automatic |
| Database | PostgreSQL 18 | PostgreSQL (free) |
| Cache | Redis 7 | Redis (free) |
| HTTPS | No | Yes |
| URL | localhost:8080 | *.onrender.com |
| Cost | Free | Free |

