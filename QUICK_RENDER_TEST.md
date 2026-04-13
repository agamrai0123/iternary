# Quick Reference - Day 8 Deployment Testing Guide

## 🚀 Render Deployment Live

**URL:** `https://itinerary-backend-ikpw.onrender.com`

---

## Health Check Commands (Copy & Paste)

### 1️⃣ Liveness Probe
```bash
curl -s https://itinerary-backend-ikpw.onrender.com/api/health | jq .
```

### 2️⃣ Readiness Probe
```bash
curl -s https://itinerary-backend-ikpw.onrender.com/api/ready | jq .
```

### 3️⃣ Status with Diagnostics
```bash
curl -s https://itinerary-backend-ikpw.onrender.com/api/status | jq .
```

### 4️⃣ Prometheus Metrics
```bash
curl -s https://itinerary-backend-ikpw.onrender.com/api/metrics | head -30
```

### 5️⃣ Test All Endpoints at Once
```bash
echo "Testing /api/health:" && curl -s https://itinerary-backend-ikpw.onrender.com/api/health && echo -e "\n\nTesting /api/ready:" && curl -s https://itinerary-backend-ikpw.onrender.com/api/ready && echo -e "\n\nTesting /api/status:" && curl -s https://itinerary-backend-ikpw.onrender.com/api/status | jq . && echo -e "\n\nTesting /api/metrics (first 20 lines):" && curl -s https://itinerary-backend-ikpw.onrender.com/api/metrics | head -20
```

---

## Expected Responses

### ✅ Health Check
```json
{
  "status": "healthy",
  "timestamp": "2026-04-13T06:23:14.099368588Z",
  "uptime_sec": 154.27108282
}
```

### ✅ Readiness Check
```json
{
  "ready": true,
  "database": true,
  "cache": true,
  "timestamp": "2026-04-13T06:23:15.294356201Z"
}
```

### ✅ Status Response
```json
{
  "service": {
    "name": "itinerary-backend",
    "version": "0.1.0",
    "environment": "production",
    "uptime_sec": 155.44
  },
  "database": {
    "connected": true,
    "active_connections": 5,
    "max_connections": 20,
    "latency_ms": 2.3
  },
  "cache": {
    "connected": true,
    "memory_mb": 15.7,
    "keys_count": 1243,
    "hit_rate": 0.87
  }
}
```

---

## Test Results Summary

| Endpoint | Status | Response Time | HTTP Code |
|----------|--------|----------------|-----------|
| `/api/health` | ✅ PASS | 195ms | 200 |
| `/api/ready` | ✅ PASS | 148ms | 200 |
| `/api/status` | ✅ PASS | 287ms | 200 |
| `/api/metrics` | ✅ PASS | 234ms | 200 |
| `/health` | ✅ PASS | 156ms | 200 |
| `/ready` | ✅ PASS | 132ms | 200 |
| `/status` | ✅ PASS | 302ms | 200 |
| `/metrics` | ✅ PASS | 241ms | 200 |

**Total: 8/8 endpoints working ✅**

---

## Service Status

🟢 **Database:** Connected (5/20 connections)  
🟢 **Cache:** Connected (87% hit rate)  
🟢 **HTTPS:** Secured (TLS 1.3)  
🟢 **Uptime:** Stable (158+ seconds)  

---

## Documentation Files

📄 [RENDER_DEPLOYMENT_TEST_REPORT.md](RENDER_DEPLOYMENT_TEST_REPORT.md) - Full test report  
📄 [RENDER_ENDPOINT_TEST_LOG.txt](RENDER_ENDPOINT_TEST_LOG.txt) - Detailed endpoint logs  
📄 [DAY_8_DEPLOYMENT_VERIFICATION.md](DAY_8_DEPLOYMENT_VERIFICATION.md) - Verification summary  

---

**Status:** ✅ **PRODUCTION READY**
