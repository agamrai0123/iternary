# Itinerary Project 🧳

A comprehensive travel itinerary planning platform built with Go, allowing users to discover, share, and customize travel plans with transparent pricing.

## 🚀 Quick Start

1. **Browse Documentation** → See `docs/DOCUMENTATION_NAVIGATOR.md`
2. **Get Started** → See `docs/getting-started/GETTING_STARTED.md`
3. **Run the App** → See `itinerary-backend/README.md`
4. **Deploy** → See `docs/deployment/DEPLOYMENT_GUIDE.md`

---

## 📁 Clean Project Structure

```
iternary/
├── docs/                        # 📚 All Documentation (organized by topic)
│   ├── DOCUMENTATION_NAVIGATOR.md # ← Master index (start here!)
│   ├── getting-started/         # Entry point guides
│   ├── guides/                  # Developer guides & API docs
│   ├── architecture/            # System design & strategy
│   ├── deployment/              # Deployment documentation
│   ├── phases/                  # Sprint & execution reports
│   ├── api/                     # API & integration docs
│   ├── reference/               # Requirements & roadmap
│   └── archived/                # Historical documentation
│
├── itinerary-backend/           # 💻 Go application
│   ├── itinerary/               # Business logic
│   ├── config/                  # Configuration
│   ├── static/                  # Frontend assets
│   ├── templates/               # HTML templates
│   └── README.md               # Backend guide
│
├── docker/                      # 🐳 Docker configuration
│   ├── Dockerfile
│   └── docker-compose.yml
│
├── scripts/                     # 🛠️ Automation scripts
├── dev/                         # 🔧 Dev utilities & migrations
├── archives/                    # Old phase documentation
│
└── Essential Files
    ├── README.md (you are here)
    ├── .gitignore
    ├── .env.development
    └── .github/
```

## 📖 Documentation Navigator

**👉 START HERE: `docs/DOCUMENTATION_NAVIGATOR.md`**

Complete documentation is organized into 8 categories:
- **Getting Started** - Setup and first steps
- **Guides** - Developer guides, API docs
- **Architecture** - System design
- **Deployment** - Production guides
- **Phases** - Sprint reports & summaries
- **API** - API documentation
- **Reference** - Requirements, roadmap
- **Archived** - Historical docs

## 🎯 What to Read First

| You Are... | Start With | 
|-----------|-----------|
| 🆕 **New to project** | `docs/getting-started/GETTING_STARTED.md` |
| 🧑‍💻 **A developer** | `docs/guides/BEGINNER_DEVELOPER_GUIDE.md` |
| 🚀 **Deploying** | `docs/deployment/DEPLOYMENT_GUIDE.md` |
| 📋 **Project manager** | `docs/reference/PROJECT_REQUIREMENTS.md` |
| 🤔 **Looking for docs** | `docs/DOCUMENTATION_NAVIGATOR.md` |

## ✨ Key Features

- **User Authentication** - Secure login & registration
- **Itinerary Browsing** - Discover shared travel plans
- **Social Features** - Like, comment & rate itineraries
- **Copy & Customize** - Duplicate and modify itineraries
- **Group Management** - Create and manage travel groups
- **Multi-Currency** - Multiple currency support
- **Transparent Pricing** - See all costs upfront

## 🧯 5-Minute Quick Start

```bash
# Navigate to backend
cd itinerary-backend

# Install dependencies
go mod download

# Build
go build -o itinerary-backend

# Run
./itinerary-backend

# Open browser: http://localhost:8080
```

For detailed setup: `docs/getting-started/GETTING_STARTED.md`

## 🛠️ Tech Stack

- **Backend** - Go + Gorilla/mux
- **Frontend** - HTML5/CSS3/JavaScript
- **Database** - SQLite3
- **Containers** - Docker & Docker Compose
- **Testing** - Go testing package

## 📊 Documentation Structure

```
docs/
├── DOCUMENTATION_NAVIGATOR.md   # ← Master index
├── getting-started/             # Quick start guides
├── guides/                       # How-to & API docs
├── architecture/                # Design & strategy
├── deployment/                  # Deployment guides
├── phases/                       # Sprint reports
├── api/                          # API documentation
├── reference/                    # Requirements
└── archived/                     # Historical
```

**All docs are organized and grouped by topic.**

## 🧪 Testing

```bash
cd itinerary-backend
go test ./...
go test -v ./itinerary
```

## 🐳 Docker

```bash
cd docker
docker-compose up --build
```

## 📁 Why This Structure?

✅ **Clear organization** - Docs grouped by topic  
✅ **Easy navigation** - Master index in `DOCUMENTATION_NAVIGATOR.md`  
✅ **Clean root** - Only essential files  
✅ **Scalable** - Room to grow without clutter  
✅ **Professional** - Industry-standard layout  

## 🔒 Git Configuration

`.gitignore` excludes:
- Build artifacts & binaries
- Databases & logs
- Environment files
- Development utilities
- IDE settings & OS files

## 📞 Common Tasks

**Getting started?**
→ `docs/getting-started/GETTING_STARTED.md`

**Understanding API?**
→ `docs/guides/API_REFERENCE.md`

**Need to deploy?**
→ `docs/deployment/DEPLOYMENT_GUIDE.md`

**Want system overview?**
→ `docs/architecture/IMPLEMENTATION_STRATEGY.md`

**See all requirements?**
→ `docs/reference/PROJECT_REQUIREMENTS.md`

**Can't find something?**
→ `docs/DOCUMENTATION_NAVIGATOR.md` (complete index)

---

**Status:** Production Ready ✅  
**Last Updated:** April 12, 2026

**Next Step:** Open `docs/DOCUMENTATION_NAVIGATOR.md` for complete documentation map
