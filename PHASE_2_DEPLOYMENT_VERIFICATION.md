# Phase 2 Deployment & Verification Guide

**Purpose:** Before starting Phase 2 development, verify Phase 1 is stable in production  
**Timeline:** 1 week post Phase 1 deployment  
**Ownership:** DevOps + QA Team

---

## Pre-Phase 2 Decision Checklist

**Go/No-Go for Phase 2 Initiation**

### Phase 1 Production Health (Week 1 Post-Deployment)

**Security Metrics**
- [ ] Zero unhandled auth bypass attempts
- [ ] HTTPS enforced on all endpoints
- [ ] No JWT token leaks in logs
- [ ] Rate limiting blocking brute force attempts
- [ ] Password hashing confirmed (bcrypt verification succeeds)
- [ ] Session invalidation working (logout preventing re-auth)

**Performance Metrics**
- [ ] Login response time < 200ms (p95)
- [ ] Auth middleware overhead < 5ms per request
- [ ] JWT validation < 1ms per request
- [ ] Rate limiter memory stable (no growth over 7 days)

**Stability Metrics**
- [ ] Zero authentication crashes (>99.99% uptime)
- [ ] Zero connection pool exhaustion
- [ ] Redis session store highly available
- [ ] Database handles concurrent logins

**Error Rate Targets**
- [ ] Auth errors < 0.1% (excluding user mistakes)
- [ ] Rate limit rejections < 0.01%
- [ ] Session lookup failures < 0.01%
- [ ] JWT validation errors < 0.05%

---

## Production Monitoring Setup

### Logging to Monitor

**1. Authentication Audit Log**
```
[2026-04-13T10:30:45Z] auth.login user_id=user123 email=test@example.com status=success ip=203.0.113.45 duration_ms=125
[2026-04-13T10:30:50Z] auth.logout user_id=user123 status=success ip=203.0.113.45
[2026-04-13T10:30:55Z] auth.mfa_setup user_id=user456 status=initiated ip=203.0.113.46
```

**Log locations:**
- Application logs: `/var/log/itinerary/app.log`
- Security logs: `/var/log/itinerary/security.log`
- Error logs: `/var/log/itinerary/error.log`

**Monitoring Queries (for your log aggregation tool):**
```sql
-- Failed login attempts (last hour)
SELECT user_id, COUNT(*) as attempts 
FROM auth_logs 
WHERE event='login' AND status='failed' 
AND timestamp > NOW() - INTERVAL 1 HOUR 
GROUP BY user_id HAVING count(*) > 5

-- Rate limit hits (last 24h)
SELECT ip_address, COUNT(*) as hits 
FROM auth_logs 
WHERE event='rate_limit_exceeded' 
AND timestamp > NOW() - INTERVAL 24 HOUR 
GROUP BY ip_address

-- JWT validation errors (last hour)
SELECT COUNT(*), error_reason 
FROM auth_logs 
WHERE event='jwt_validation_failed' 
AND timestamp > NOW() - INTERVAL 1 HOUR 
GROUP BY error_reason
```

### Alerts to Configure

| Alert | Threshold | Action |
|-------|-----------|--------|
| Auth errors spike | >5% above baseline | Page on-call |
| Failed login rate | >50 per minute | Lock suspicious IP |
| Rate limiter memory leak | >100MB after 24h | Restart service |
| JWT validation errors | >1% of requests | Investigate JWT config |
| Redis session lookup fails | >100ms p95 | Scale Redis |
| Bcrypt verification slow | >500ms p95 | Reduce load or scale |

---

## Week 1 Verification Playbook

### Day 1: Immediate Health Check

**Monday Morning - First 2 Hours Post-Deployment**

```bash
# 1. Verify service is running
curl -X GET https://api.itinerary.com/health

# 2. Test all auth endpoints
curl -X POST https://api.itinerary.com/api/v1/auth/login \
  -d '{"email":"test@example.com","password":"testpass123"}' \
  -H "Content-Type: application/json"

# 3. Verify HTTPS only
curl -X GET http://api.itinerary.com/health
# Should redirect to HTTPS

# 4. Check logs for errors
tail -f /var/log/itinerary/error.log

# 5. Monitor memory usage
watch 'ps aux | grep itinerary'
```

**Success Indicators:**
- ✅ Service responds to requests
- ✅ HTTPS enforced
- ✅ Auth works with real credentials
- ✅ No critical errors in logs
- ✅ Memory stable around 100-150MB

### Day 2-3: Load Testing

**Tuesday - Baseline Load Test**

```bash
# Test 100 concurrent logins
ab -n 1000 -c 100 \
  -p login_payload.json \
  -T "application/json" \
  -H "Content-Type: application/json" \
  https://api.itinerary.com/api/v1/auth/login

# Record results:
# - Requests per second
# - Response time (min/avg/max)
# - Failed requests
# - Error rate
```

**Test Scenarios:**
1. Valid credential login (record baseline)
2. Invalid password repeated (verify rate limiting)
3. Missing email (verify validation)
4. Expired token handling
5. Concurrent session management

**Expected Results:**
- ✅ >500 req/sec for valid logins
- ✅ <5% error rate (excluding expected errors)
- ✅ Rate limiting activates at 5 attempts/min
- ✅ Token refresh works smoothly

### Day 4-5: Security Testing

**Wednesday-Thursday - Security Validation**

**Test Matrix:**

| Test Case | Input | Expected | Status |
|-----------|-------|----------|--------|
| Weak password | `pass` | Rejected with 400 | - |
| SQL injection | `' OR '1'='1` | Escaped/rejected | - |
| JWT tampering | Modified exp claim | Rejected 401 | - |
| Expired token | 24h+ old token | Rejected 401 | - |
| Rate limit bypass | 10 requests/sec same IP | 429 after limit | - |
| No HTTPS | HTTP request | Redirected to HTTPS | - |
| Missing auth header | No JWT token | 401 Unauthorized | - |
| Invalid JWT format | Malformed token | 401 Unauthorized | - |

**Run Security Tests:**
```bash
# Test 1: Weak password rejection
curl -X POST https://api.itinerary.com/api/v1/auth/login \
  -d '{"email":"user@example.com","password":"123"}' \
  -H "Content-Type: application/json"
# Expect: 400 Bad Request

# Test 2: Rate limiting
for i in {1..10}; do
  curl -X POST https://api.itinerary.com/api/v1/auth/login \
    -d '{"email":"attacker@example.com","password":"wrong"}' \
    -H "Content-Type: application/json"
  echo ""
done
# After 5: Expect 429 Too Many Requests

# Test 3: JWT tampering
OLD_TOKEN="eyJhbGc..." # Get real token
TAMPERED_TOKEN="${OLD_TOKEN:0:-10}AAAAAAAAAA"
curl -X GET https://api.itinerary.com/api/v1/users/profile \
  -H "Authorization: Bearer $TAMPERED_TOKEN"
# Expect: 401 Unauthorized
```

**Security Checklist:**
- [ ] No plaintext passwords in logs
- [ ] No tokens visible in debug output
- [ ] HTTPS certificate valid
- [ ] CORS headers proper
- [ ] Session cookies httpOnly + Secure
- [ ] No default credentials accessible
- [ ] Rate limiting prevents brute force
- [ ] Error messages don't leak data

### Day 5-7: Real User Testing

**Friday-Sunday - Canary User Testing**

**Internal Testing (Team):**
```
1. Test complete login flow
2. Test logout and session invalidation
3. Test token refresh
4. Test profile update (requires auth)
5. Test multi-device sessions
6. Test logout all devices
```

**Canary Users (10 selected users):**
- [ ] Beta users given access
- [ ] Issues tracked in spreadsheet
- [ ] Daily check-ins (end of business)
- [ ] Any issues immediately escalated

---

## Metrics Dashboard Template

**Create Dashboard in Your Monitoring Tool (Datadog/Grafana/CloudWatch):**

```
┌─────────────────────────────────────────────────────┐
│         PHASE 1 SECURITY IMPLEMENTATION             │
│            Week 1 Production Monitoring             │
├─────────────────────────────────────────────────────┤
│                                                       │
│  Successful Logins: 2,347   ↑ 2.1%                  │
│  Failed Logins: 145 ↓ 0.2%                         │
│  Rate Limit Hits: 23 ↓ 0.05%                       │
│  Auth Errors: 0.8% → Target < 0.1%                │
│                                                       │
│  Avg Login Time: 128ms ↓ 12ms improvement          │
│  P95 Login Time: 245ms (Target < 300ms) ✅         │
│  JWT Validation: 0.8ms avg ✅                       │
│                                                       │
│  HTTPS Enforcement: 100% ✅                         │
│  Session Valid Rate: 99.98% ✅                      │
│  Redis Availability: 99.99% ✅                      │
│  Bcrypt CPU Time: 92ms avg (Cost 12) ✅            │
│                                                       │
│  No Security Incidents: ✅                          │
│  Zero Token Leaks: ✅                              │
│  Password Hash Confirmed: ✅                        │
│  Rate Limiter Effective: ✅                         │
│                                                       │
└─────────────────────────────────────────────────────┘
```

---

## Go/No-Go Decision Criteria

### GO for Phase 2 if All True:

**Security ✅**
- [ ] Phase 1 passed security audit
- [ ] Zero critical vulnerabilities found
- [ ] HTTPS in production for 7 days
- [ ] No token leaks in 7 days
- [ ] Rate limiting working (blocked >99% of brute force attempts)

**Performance ✅**
- [ ] Auth ops < 200ms p95
- [ ] Memory stable (no growth)
- [ ] CPU usage < 30% sustained
- [ ] Database query times < 50ms p95

**Stability ✅**
- [ ] >99.95% uptime (max 22 minutes downtime in 7 days)
- [ ] Zero unplanned restarts
- [ ] Session store highly available
- [ ] Error rates < 0.1%

**User Impact ✅**
- [ ] Canary users report no major issues
- [ ] Support tickets < 5
- [ ] Adoption rate > 80%
- [ ] User feedback positive

**Team Readiness ✅**
- [ ] Release notes documented
- [ ] Deployment playbook tested
- [ ] Rollback procedure verified
- [ ] Phase 2 development environment ready

### NO-GO Decision Triggers:

🔴 **Stop and Investigate if:**
- Authentication rate > 1% errors
- Any security incident discovered
- Uptime drops below 99%
- Memory leak detected
- Critical CVE in dependencies

**Action on No-Go:**
1. Immediately roll back to Phase 1 initial state
2. Investigate root cause (2-day max)
3. Fix identified issues
4. Re-run verification (3-day shorter cycle)
5. Re-evaluate go/no-go

---

## Phase 2 Sprint 1 Start Criteria

**Approved to Start Sprint 1 when:**
- [ ] Go/No-Go checklist: ALL PASSED
- [ ] Phase 1 production stability: 7 days confirmed
- [ ] Team capacity: 2+ engineers available
- [ ] Development environment: Ready
- [ ] Feature branch: Created (feature/phase2-mfa-oauth)
- [ ] Dependencies: Added to go.mod
- [ ] OAuth apps: Configured (GitHub + Google)
- [ ] Database migrations: Prepared
- [ ] Test plan: Documented

---

## Rollout Timeline

**If Go Decision Made (Day 7-8):**

```
Day 1 (Monday):
  - Announce Phase 2 start
  - Create feature branch
  - Stage MFA database schema

Day 2-3 (Tue-Wed):
  - Implement TOTP core
  - Build MFA handlers
  - Database integration testing

Day 4-5 (Thu-Fri):
  - OAuth framework setup
  - GitHub OAuth integration
  - Initial testing

Day 6-7 (Weekend):
  - Google OAuth integration
  - Cross-browser testing
  - Documentation finalization

Day 8-10 (Mon-Wed next week):
  - Integration testing
  - Performance testing
  - Deployment prep

Day 11 (Thursday):
  - Code review + approval
  - Merge to staging
  - Stage testing

Day 12 (Friday):
  - Final validation
  - Deploy to production (if clear)
```

---

## Success Metrics for Sprint 1

By end of Sprint 1 (2 weeks from Phase 2 start):

✅ **Functionality**
- [ ] MFA setup working for test users
- [ ] TOTP codes valid in authenticator apps
- [ ] Backup codes functional
- [ ] GitHub OAuth login working
- [ ] Google OAuth login working (or ready)

✅ **Quality**
- [ ] 85%+ test coverage (MFA module)
- [ ] Zero high/critical security issues
- [ ] Performance: MFA setup < 500ms
- [ ] Performance: MFA verification < 100ms
- [ ] OAuth callback < 1 second

✅ **Documentation**
- [ ] MFA user guide published
- [ ] OAuth setup guide published
- [ ] API documentation complete
- [ ] Deployment procedures documented

✅ **Operations**
- [ ] Monitoring in place
- [ ] Alerts configured
- [ ] Log aggregation working
- [ ] Incident response plan ready

---

**Phase 2 Deployment Verification Complete**

Once above checklist is reviewed and approved, move forward with confidence!
