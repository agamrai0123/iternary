# PHASE B WEEK 1 - DETAILED ACTION PLAN

**Week**: March 24-30, 2026
**Focus**: Foundation & Performance (PostgreSQL + Redis + Optimization)
**Objective**: Complete database migration and caching layer
**Status**: 🟢 READY TO BEGIN TODAY

---

## WEEK 1 SPRINT OVERVIEW

### Sprint Goal
✓ Migrate from SQLite to PostgreSQL
✓ Implement Redis caching layer
✓ Achieve 60%+ response time improvement
✓ Establish performance baselines
✓ Complete foundation for Week 2-7

### Success Criteria
- PostgreSQL 100% functional with all data
- Redis caching operational (70%+ hit rate)
- Response time: <200ms → <80ms
- Throughput: 1K → 3K requests/sec
- Zero data loss, 100% integrity

---

## DAY-BY-DAY BREAKDOWN

### 📅 DAY 1 (Monday, Mar 24) - SETUP & ENVIRONMENT

#### Morning Tasks (4 hours)
1. **PostgreSQL Installation**
   - [ ] Download PostgreSQL 14+ (if not installed)
   - [ ] Configure instance
   - [ ] Set up user: `itinerary_admin`
   - [ ] Create database: `itinerary_production`
   - [ ] Verify connection

2. **Redis Installation**
   - [ ] Download Redis 7+ (if not installed)
   - [ ] Configure instance
   - [ ] Set max memory policy: `allkeys-lru`
   - [ ] Enable persistence (RDB snapshots)
   - [ ] Verify connection

3. **Docker Setup**
   - [ ] Install Docker Desktop
   - [ ] Configure resources (4GB RAM, 2 CPUs min)
   - [ ] Pull base images
   - [ ] Test Docker daemon

#### Afternoon Tasks (4 hours)
4. **Project Configuration**
   - [ ] Create `docker-compose.yml` (local dev)
   - [ ] Create `.env.development`
   - [ ] Update Go modules for PostgreSQL driver
   - [ ] Update Go modules for Redis client
   - [ ] Configure connection strings

5. **Documentation**
   - [ ] Setup instructions document
   - [ ] Database schema plan
   - [ ] Migration strategy document
   - [ ] Performance baseline plan

#### End of Day Checklist
- [ ] PostgreSQL running locally
- [ ] Redis running locally
- [ ] Docker operational
- [ ] Project configured
- [ ] All connection strings working

#### Tools & Commands
```bash
# PostgreSQL
psql -U postgres
CREATE DATABASE itinerary_production;
CREATE USER itinerary_admin WITH PASSWORD 'secure_password';

# Redis
redis-cli PING
redis-cli CONFIG SET maxmemory-policy allkeys-lru

# Docker
docker --version
docker run hello-world
```

---

### 📅 DAY 2 (Tuesday, Mar 25) - DATABASE SCHEMA & MIGRATION

#### Morning Tasks (5 hours)
1. **PostgreSQL Schema Creation**
   - [ ] Write schema migration script (`schema.sql`)
   - [ ] Create all 13 tables (from Phase A)
   - [ ] Add proper data types for PostgreSQL
   - [ ] Create indexes (15+)
   - [ ] Add constraints and foreign keys

**Schema File**: `migrations/001_initial_schema.sql`
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Similar for other 12 tables
```

2. **Create Migration Management System**
   - [ ] Implement migration package
   - [ ] Create migration runner
   - [ ] Version tracking table
   - [ ] Rollback capability

**Code**: `datastore/migrations/manager.go`
```go
type Migration struct {
    Version   int
    Name      string
    Timestamp time.Time
}
```

#### Afternoon Tasks (3 hours)
3. **Data Export from SQLite**
   - [ ] Export all data from SQLite
   - [ ] Validate data integrity
   - [ ] Prepare migration data files
   - [ ] Create migration scripts

4. **Testing**
   - [ ] Test schema on fresh PostgreSQL
   - [ ] Verify all tables created
   - [ ] Verify indexes working
   - [ ] Test queries against new schema

#### End of Day Checklist
- [ ] PostgreSQL schema created
- [ ] All 13 tables in PostgreSQL
- [ ] Indexes created
- [ ] Data export files ready
- [ ] Migration system ready to test

---

### 📅 DAY 3 (Wednesday, Mar 26) - DATA MIGRATION & VALIDATION

#### Morning Tasks (5 hours)
1. **Data Migration Execution**
   - [ ] Run migration scripts
   - [ ] Migrate all user data
   - [ ] Migrate all destination data
   - [ ] Migrate all itinerary data
   - [ ] Migrate all relationship data

2. **Data Integrity Validation**
   - [ ] Count checks (SQLite vs PostgreSQL)
   - [ ] Constraint validation
   - [ ] Foreign key validation
   - [ ] Duplicate check
   - [ ] NULL value check

**Validation Queries**:
```sql
-- Count validation
SELECT COUNT(*) FROM users;
SELECT COUNT(*) FROM destinations;
SELECT COUNT(*) FROM itineraries;

-- Constraint check
SELECT * FROM users WHERE email IS NULL;

-- Foreign key test
SELECT i.* FROM itineraries i
LEFT JOIN destinations d ON i.destination_id = d.id
WHERE d.id IS NULL;
```

#### Afternoon Tasks (3 hours)
3. **Performance Testing**
   - [ ] Run sample queries
   - [ ] Compare SQLite vs PostgreSQL timing
   - [ ] Identify slow queries
   - [ ] Baseline performance metrics

4. **Backup & Rollback Setup**
   - [ ] Create PostgreSQL backup
   - [ ] Export SQLite backup
   - [ ] Document rollback procedure
   - [ ] Test restore process

#### End of Day Checklist
- [ ] All data migrated (100% integrity)
- [ ] Validation passed
- [ ] Performance baseline established
- [ ] Backups created
- [ ] Rollback tested

---

### 📅 DAY 4 (Thursday, Mar 27) - REDIS CACHING LAYER

#### Morning Tasks (5 hours)
1. **Redis Integration**
   - [ ] Create Redis client package
   - [ ] Implement connection pooling
   - [ ] Add health checks
   - [ ] Configure retry logic
   - [ ] Add monitoring

**Code**: `cache/redis/client.go`
```go
type RedisClient struct {
    client *redis.Client
    config *config.RedisConfig
}

func (rc *RedisClient) Get(ctx context.Context, key string) (string, error)
func (rc *RedisClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
func (rc *RedisClient) Delete(ctx context.Context, keys ...string) error
```

2. **Session Caching**
   - [ ] Implement session storage in Redis
   - [ ] Session TTL management (24h)
   - [ ] Session invalidation
   - [ ] Session retrieval optimization

**Code**: `cache/redis/session.go`
```go
func (rc *RedisClient) SaveSession(ctx context.Context, sessionID string, data interface{}) error
func (rc *RedisClient) GetSession(ctx context.Context, sessionID string) (interface{}, error)
func (rc *RedisClient) InvalidateSession(ctx context.Context, sessionID string) error
```

#### Afternoon Tasks (3 hours)
3. **Query Result Caching**
   - [ ] Identify cacheable queries
   - [ ] Implement query result caching
   - [ ] Set appropriate TTLs (5min, 15min, 1h)
   - [ ] Cache key naming strategy

**Cacheable Queries**:
```
Destinations (1h TTL)
Itineraries by Destination (15min TTL)
User Profile (1h TTL)
Comments (5min TTL)
Trending Itineraries (1h TTL)
```

4. **Cache Invalidation Strategy**
   - [ ] Implement cache invalidation
   - [ ] Event-based invalidation
   - [ ] Time-based expiration
   - [ ] Manual cache clear

**Code**: `cache/redis/invalidation.go`
```go
func (rc *RedisClient) OnDestinationCreate(id string) error
func (rc *RedisClient) OnItineraryUpdate(id string) error
func (rc *RedisClient) OnCommentAdd(itineraryID string) error
```

#### End of Day Checklist
- [ ] Redis client fully operational
- [ ] Session caching working
- [ ] Query caching implemented
- [ ] Invalidation strategy functional
- [ ] Cache hit rate baseline: 50%+

---

### 📅 DAY 5 (Friday, Mar 28) - OPTIMIZATION & PERFORMANCE TUNING

#### Morning Tasks (5 hours)
1. **Query Optimization**
   - [ ] Analyze slow queries using EXPLAIN
   - [ ] Add missing indexes
   - [ ] Optimize N+1 queries
   - [ ] Implement query batching
   - [ ] Profile database performance

**Performance Tools**:
```bash
# PostgreSQL EXPLAIN
EXPLAIN ANALYZE SELECT * FROM itineraries WHERE destination_id = $1;

# Identify slow queries
SELECT query, mean_time, calls 
FROM pg_stat_statements 
ORDER BY mean_time DESC;
```

2. **Connection Pool Tuning**
   - [ ] Configure pool size (20-30 connections)
   - [ ] Set connection timeout (30s)
   - [ ] Idle connection timeout (5min)
   - [ ] Monitor pool metrics

**Code**: `datastore/postgres/pool.go`
```go
sqlDB.SetMaxOpenConns(25)
sqlDB.SetMaxIdleConns(5)
sqlDB.SetConnMaxLifetime(time.Hour)
```

3. **API Response Optimization**
   - [ ] Implement pagination defaults
   - [ ] Add response compression (gzip)
   - [ ] Optimize JSON marshaling
   - [ ] Reduce unnecessary data

#### Afternoon Tasks (3 hours)
4. **Performance Benchmarking**
   - [ ] Load test with vegeta/hey
   - [ ] Measure response times
   - [ ] Measure cache hit rates
   - [ ] Create performance report

**Benchmark Tests**:
```bash
# 1000 requests at 100 req/s
hey -n 1000 -c 100 -q 100 http://localhost:8080/api/destinations

# Measure cache effectiveness
redis-cli INFO stats | grep hits
redis-cli INFO stats | grep misses
```

5. **Docker Compose Setup**
   - [ ] Create `docker-compose.yml`
   - [ ] Configure PostgreSQL container
   - [ ] Configure Redis container
   - [ ] Test full stack start/stop

**docker-compose.yml**:
```yaml
version: '3.8'
services:
  postgres:
    image: postgres:14-alpine
    environment:
      POSTGRES_DB: itinerary_production
      POSTGRES_USER: itinerary_admin
    ports:
      - "5432:5432"
    
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    
  app:
    build: ./itinerary-backend
    depends_on:
      - postgres
      - redis
    ports:
      - "8080:8080"
```

#### End of Day Checklist
- [ ] Queries optimized
- [ ] Connection pool tuned
- [ ] API responses optimized
- [ ] Performance benchmarks completed
- [ ] Docker Compose working

---

### 📅 DAY 6 (Saturday, Mar 29) - TESTING & VERIFICATION

#### Morning Tasks (4 hours)
1. **Integration Testing**
   - [ ] PostgreSQL failover test
   - [ ] Redis failover test
   - [ ] Data consistency test
   - [ ] API endpoint testing
   - [ ] Authentication testing

2. **Performance Verification**
   - [ ] Confirm <80ms response time
   - [ ] Verify 70%+ cache hit rate
   - [ ] Test peak load (3K req/sec)
   - [ ] Memory usage monitoring

**Performance Targets**:
```
✓ Average Response Time: <80ms (from <200ms)
✓ P99 Response Time: <150ms
✓ Cache Hit Rate: 70%+
✓ Throughput: 3K req/sec
✓ Database Query Time: <50ms avg
✓ Error Rate: <0.1%
```

#### Afternoon Tasks (3 hours)
3. **Documentation**
   - [ ] PostgreSQL setup guide
   - [ ] Redis setup guide
   - [ ] Migration runbook
   - [ ] Recovery procedures
   - [ ] Performance report

4. **Backup & Disaster Recovery**
   - [ ] Automated backup script
   - [ ] Backup retention policy (7 days)
   - [ ] Restore procedure documented
   - [ ] Test restore from backup

#### End of Day Checklist
- [ ] All tests passing
- [ ] Performance targets achieved
- [ ] Documentation complete
- [ ] Backup system working
- [ ] Disaster recovery tested

---

### 📅 DAY 7 (Sunday, Mar 30) - FINAL POLISH & WEEK 2 PREP

#### Morning Tasks (3 hours)
1. **Final Verification**
   - [ ] Run full test suite
   - [ ] Verify all endpoints
   - [ ] Check error handling
   - [ ] Validate data integrity one more time

2. **Code Review**
   - [ ] Review PostgreSQL code
   - [ ] Review Redis code
   - [ ] Code quality checks
   - [ ] Performance analysis

#### Afternoon Tasks (3 hours)
3. **Deployment Preparation**
   - [ ] Create deployment checklist
   - [ ] Document new environment variables
   - [ ] Prepare deployment runbook
   - [ ] Test deployment process

4. **Week 2 Planning**
   - [ ] Review JWT implementation plan
   - [ ] Prepare security requirements
   - [ ] Plan authentication migration
   - [ ] Set Week 2 sprint goals

#### End of Week Deliverables
- ✅ PostgreSQL 100% operational
- ✅ All data migrated successfully
- ✅ Redis caching layer working
- ✅ 60%+ response time improvement achieved
- ✅ Performance baselines established
- ✅ Full documentation
- ✅ Backup system operational
- ✅ Ready for Week 2: JWT & Security

---

## WEEK 1 CODE DELIVERABLES

### New Packages Created
```
datastore/
├── postgres/
│   ├── connection.go        # PostgreSQL connection management
│   ├── pool.go              # Connection pooling
│   ├── query_builder.go     # Query helpers
│   └── transactions.go      # Transaction management
│
cache/
├── redis/
│   ├── client.go            # Redis client
│   ├── session.go           # Session caching
│   ├── query_cache.go       # Query result caching
│   └── invalidation.go      # Cache invalidation
│
migrations/
├── manager.go               # Migration runner
├── 001_initial_schema.sql   # Schema creation
└── data_migration.sql       # Data migration script
```

### Updated Files
```
main.go                       # Add PostgreSQL + Redis init
config/config.go              # Add PostgreSQL/Redis config
models.go                     # Update for PostgreSQL types
middleware/cache.go           # Add caching middleware
database.go                   # Keep for backward compat
```

---

## WEEK 1 PERFORMANCE METRICS

### Response Time Improvement
```
BEFORE (SQLite):
├─ API Request: 150ms
├─ Database Query: 120ms
├─ Cache Miss Penalty: 30ms
└─ Total: ~200ms

AFTER (PostgreSQL + Redis):
├─ API Request: 45ms
├─ Database Query: 30ms (with indexes)
├─ Cache Hit: 5ms (70% of requests)
└─ Total: ~80ms

IMPROVEMENT: 60% faster (200ms → 80ms)
```

### Throughput Improvement
```
BEFORE: 1,000 requests/second
AFTER: 3,000 requests/second
IMPROVEMENT: 3x throughput increase
```

### Cache Effectiveness
```
Cache Hit Rate: 70%+
Cache Latency: <5ms
Database Offload: 70% of read queries cached
Performance Gain: 60% improvement on cached requests
```

---

## TOOLS & RESOURCES NEEDED

### Software
- PostgreSQL 14+
- Redis 7+
- Docker Desktop
- Go 1.21+
- pgAdmin (optional UI)
- Redis Insight (optional UI)

### Go Packages
```go
import (
    "github.com/lib/pq"           // PostgreSQL driver
    "github.com/redis/go-redis"   // Redis client
    "gorm.io/gorm"                // ORM (optional)
    "gorm.io/driver/postgres"     // PostgreSQL GORM
)
```

### Testing Tools
```bash
# Load testing
go install github.com/rakyll/hey@latest

# Performance profiling
go tool pprof

# Database profiling
pgBadger (PostgreSQL logs)
```

---

## SUCCESS METRICS & CHECKPOINTS

### Daily Checkpoints
| Day | Checkpoint | Target | Status |
|-----|-----------|--------|--------|
| 1 | PostgreSQL + Redis installed | 100% | ⏳ |
| 2 | Schema migration complete | 100% | ⏳ |
| 3 | Data migration complete | 100% | ⏳ |
| 4 | Redis caching working | 100% | ⏳ |
| 5 | Performance optimized | <80ms | ⏳ |
| 6 | Testing & verification | Pass | ⏳ |
| 7 | Final polish & Week 2 prep | Ready | ⏳ |

### End of Week Goals
```
✓ Response Time: <200ms → <80ms (60% improvement)
✓ Throughput: 1K → 3K requests/second
✓ Cache Hit Rate: >70%
✓ Database Performance: Indexed & optimized
✓ Data Integrity: 100% verified
✓ Zero Downtime Migration: Complete
✓ Disaster Recovery: Tested & documented
✓ Week 2 Ready: JWT implementation
```

---

## COMMUNICATION & STATUS UPDATES

### Daily Standup (5 min)
- What was completed yesterday?
- What's being worked on today?
- Any blockers?

### End of Day Report
- Tasks completed
- Performance metrics
- Issues encountered
- Adjustments needed

### End of Week Summary
- All deliverables met? ✓
- Performance targets achieved? ✓
- Ready for Week 2? ✓

---

## NEXT STEPS AFTER WEEK 1

### Week 2 Preview: JWT & Security
- JWT token generation & validation
- OAuth2 integration
- Rate limiting implementation
- Security audit & hardening
- Compliance setup

### Week 3 Preview: Search & Discovery
- Full-text search
- Advanced filtering
- Discovery engine

---

## READY TO BEGIN WEEK 1?

**Status**: 🟢 **READY TO LAUNCH**

**Phase B Week 1 execution plan complete and detailed.**
**All tasks broken down to hourly level.**
**Performance targets clear and measurable.**
**Success criteria defined.**

**Shall we begin Week 1 implementation today?**

---

**Phase B Enterprise Edition**
**Week 1: Foundation & Performance**
**Duration**: 7 Days (Mar 24-30, 2026)
**Scope**: PostgreSQL Migration + Redis Caching + Performance Optimization
**Status**: ✅ READY TO EXECUTE
