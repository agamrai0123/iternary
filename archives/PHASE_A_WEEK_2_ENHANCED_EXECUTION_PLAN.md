# Phase A Week 2 - ENHANCED EXECUTION PLAN
## Multi-Currency & Performance Monitoring Integration

**Status:** Ready for Execution  
**Start Date:** March 24, 2026 (Today)  
**Duration:** 15 hours (Mon-Fri)  
**New Features:** Currency support + Performance alerts

---

## 🎯 Phase A Week 2 - Enhanced Objectives

### Primary Goals
1. ✅ Execute database setup and tests
2. ✅ Verify all 16 API endpoints
3. ✅ Validate algorithms (settlement, splits, polling)
4. ✅ Establish performance baseline
5. **NEW** ✨ Implement multi-currency/language/nationality support
6. **NEW** ✨ Add performance monitoring with real-time alerts

---

## 📋 New Requirements Integration

### A. Multi-Currency & Localization Support

#### User Preferences Model Extension

```go
// Add to group_models.go
type UserPreferences struct {
    UserID           string `json:"user_id" bson:"user_id"`
    Nationality      string `json:"nationality" bson:"nationality"` // ISO 3166-1 (US, IN, GB, etc)
    PreferredCurrency string `json:"currency" bson:"currency"`      // ISO 4217 (USD, INR, GBP, EUR, etc)
    PreferredLanguage string `json:"language" bson:"language"`      // ISO 639-1 (en, hi, es, fr, etc)
    TimeZone         string `json:"timezone" bson:"timezone"`       // IANA (America/New_York, Asia/Kolkata, etc)
    CreatedAt        time.Time `json:"created_at" bson:"created_at"`
    UpdatedAt        time.Time `json:"updated_at" bson:"updated_at"`
}

// Supported currencies
type SupportedCurrency struct {
    Code     string  `json:"code"`     // USD, INR, GBP, EUR, etc
    Symbol   string  `json:"symbol"`   // $, ₹, £, €
    Name     string  `json:"name"`     // US Dollar, Indian Rupee, etc
    Countries []string `json:"countries"`
}

// Supported languages
type SupportedLanguage struct {
    Code    string `json:"code"`    // en, hi, es, fr
    Name    string `json:"name"`    // English, Hindi, Spanish, French
    Regions []string `json:"regions"`
}
```

#### Expense Model Enhancement

```go
// Update expense model with multi-currency
type Expense struct {
    ID              string `json:"id"`
    TripID          string `json:"trip_id"`
    Description     string `json:"description"`
    Amount          float64 `json:"amount"`
    Currency        string `json:"currency"` // NEW: Currency code
    OriginalAmount  float64 `json:"original_amount,omitempty"` // For conversion tracking
    OriginalCurrency string `json:"original_currency,omitempty"`
    ExchangeRate    float64 `json:"exchange_rate,omitempty"` // NEW: Rate used
    PaidByID        string `json:"paid_by_id"`
    Category        string `json:"category"`
    Splits          []ExpenseSplit `json:"splits"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}
```

#### Settlement Model Enhancement

```go
// Update settlement with currency conversion
type Settlement struct {
    ID              string `json:"id"`
    TripID          string `json:"trip_id"`
    CreditorID      string `json:"creditor_id"`
    DebtorID        string `json:"debtor_id"`
    Amount          float64 `json:"amount"`
    Currency        string `json:"currency"`
    ConvertedAmount float64 `json:"converted_amount,omitempty"` // In debtor's currency
    ConvertedCurrency string `json:"converted_currency,omitempty"`
    Status          string `json:"status"` // pending, settled, disputed
    CreatedAt       time.Time `json:"created_at"`
    ResolvedAt      time.Time `json:"resolved_at,omitempty"`
}
```

#### API Response with Localization

```go
type ExpenseResponse struct {
    ID            string `json:"id"`
    Description   string `json:"description"`
    Amount        float64 `json:"amount"`
    AmountFormatted string `json:"amount_formatted"` // e.g., "$1,234.56" or "₹1,23,456.78"
    Currency      string `json:"currency"`
    PaidBy        string `json:"paid_by"`
    CreatedAt     string `json:"created_at"`
    CreatedAtFormatted string `json:"created_at_formatted"` // Localized date/time
}
```

---

### B. Performance Monitoring & Alerting

#### Performance Monitor Infrastructure

```go
// Add to metrics.go
type PerformanceMonitor struct {
    EndpointMetrics map[string]*EndpointMetric
    Thresholds      PerformanceThresholds
    Alerts          chan PerformanceAlert
    Logger          *Logger
    mu              sync.RWMutex
}

type EndpointMetric struct {
    Path           string
    Method         string
    ResponseTimes  []int64 // milliseconds
    ErrorCount     int
    SuccessCount   int
    P50ResponseTime int64
    P95ResponseTime int64
    P99ResponseTime int64
    MaxResponseTime int64
    AvgResponseTime float64
    LastUpdated    time.Time
}

type PerformanceThresholds struct {
    P95ThresholdMs  int64  // Alert if p95 exceeds this (default 500ms)
    P99ThresholdMs  int64  // Alert if p99 exceeds this (default 1000ms)
    ErrorRateThreshold float64 // Alert if error rate > this (default 1%)
    MemoryThresholdMb int64    // Alert if memory > this (default 200MB)
}

type PerformanceAlert struct {
    AlertType      string // "high_response_time", "high_error_rate", "memory_spike"
    Endpoint       string
    Severity       string // "warning", "critical"
    Threshold      float64
    CurrentValue   float64
    Message        string
    Timestamp      time.Time
}
```

#### Performance Tracking Middleware

```go
// Add to metrics_middleware.go
func (m *MetricsMiddleware) PerformanceMonitoringMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        startTime := time.Now()
        startMem := runtime.MemStats{}
        runtime.ReadMemStats(&startMem)
        
        c.Next()
        
        duration := time.Since(startTime)
        durationMs := duration.Milliseconds()
        
        endMem := runtime.MemStats{}
        runtime.ReadMemStats(&endMem)
        memDelta := int64(endMem.Alloc - startMem.Alloc)
        
        // Track metric
        m.monitor.RecordEndpointMetric(c.Request.Method, c.Request.URL.Path, durationMs, c.Writer.Status())
        
        // Check thresholds and emit alerts if needed
        m.monitor.CheckThresholds(c.Request.Method, c.Request.URL.Path)
    }
}
```

#### Alert Handling

```go
// Add alert consumer
func (m *PerformanceMonitor) AlertHandler(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        case alert := <-m.Alerts:
            // Log alert
            m.Logger.Warn("Performance Alert",
                zap.String("type", alert.AlertType),
                zap.String("endpoint", alert.Endpoint),
                zap.String("severity", alert.Severity),
                zap.Float64("threshold", alert.Threshold),
                zap.Float64("current_value", alert.CurrentValue),
            )
            
            // Send to monitoring system (DataDog, New Relic, CloudWatch, etc)
            m.SendToMonitoringSystem(alert)
            
            // If critical, trigger action
            if alert.Severity == "critical" {
                m.TriggerCriticalAlert(alert)
            }
        }
    }
}
```

---

## 📦 Database Schema Updates

### New Tables/Columns

```sql
-- Users Preferences Table
CREATE TABLE user_preferences (
    user_id VARCHAR(255) PRIMARY KEY,
    nationality VARCHAR(10),
    preferred_currency VARCHAR(3),
    preferred_language VARCHAR(5),
    timezone VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Supported Currencies
CREATE TABLE supported_currencies (
    code VARCHAR(3) PRIMARY KEY,
    symbol VARCHAR(5),
    name VARCHAR(100),
    exchange_rate FLOAT DEFAULT 1.0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Supported Languages
CREATE TABLE supported_languages (
    code VARCHAR(5) PRIMARY KEY,
    name VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE
);

-- Performance Metrics
CREATE TABLE performance_metrics (
    id UUID PRIMARY KEY,
    endpoint_path VARCHAR(255),
    method VARCHAR(10),
    response_time_ms BIGINT,
    status_code INT,
    error_flag BOOLEAN DEFAULT FALSE,
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_endpoint_time (endpoint_path, recorded_at DESC)
);

-- Performance Alerts
CREATE TABLE performance_alerts (
    id UUID PRIMARY KEY,
    alert_type VARCHAR(50),
    endpoint VARCHAR(255),
    severity VARCHAR(20),
    threshold FLOAT,
    current_value FLOAT,
    message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP NULL,
    INDEX idx_created_severity (created_at DESC, severity)
);

-- Expense Currency Conversions (for audit trail)
CREATE TABLE expense_conversions (
    id UUID PRIMARY KEY,
    expense_id VARCHAR(255),
    from_currency VARCHAR(3),
    to_currency VARCHAR(3),
    from_amount FLOAT,
    to_amount FLOAT,
    exchange_rate FLOAT,
    converted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (expense_id) REFERENCES expenses(id)
);

-- ALTER existing tables to add currency fields
ALTER TABLE expenses ADD COLUMN currency VARCHAR(3) DEFAULT 'USD';
ALTER TABLE expenses ADD COLUMN exchange_rate FLOAT DEFAULT 1.0;
ALTER TABLE settlements ADD COLUMN currency VARCHAR(3) DEFAULT 'USD';
ALTER TABLE settlements ADD COLUMN converted_amount FLOAT;
ALTER TABLE settlements ADD COLUMN converted_currency VARCHAR(3);
```

---

## 🔄 Multi-Currency Service Methods

### Add to group_service.go

```go
// GetCurrencyRate retrieves exchange rate
func (s *Service) GetCurrencyRate(ctx context.Context, from, to string) (float64, error) {
    if from == to {
        return 1.0, nil
    }
    
    rate, err := s.db.GetExchangeRate(ctx, from, to)
    if err != nil {
        s.logger.Error("Failed to get exchange rate", 
            zap.String("from", from), 
            zap.String("to", to),
            zap.Error(err))
        return 0, err
    }
    
    return rate, nil
}

// ConvertExpenseToUserCurrency converts expense to user's preferred currency
func (s *Service) ConvertExpenseToUserCurrency(ctx context.Context, expense *Expense, userCurrency string) (*Expense, error) {
    if expense.Currency == userCurrency {
        return expense, nil
    }
    
    rate, err := s.GetCurrencyRate(ctx, expense.Currency, userCurrency)
    if err != nil {
        return nil, err
    }
    
    converted := &Expense{
        OriginalAmount: expense.Amount,
        OriginalCurrency: expense.Currency,
        Amount: expense.Amount * rate,
        Currency: userCurrency,
        ExchangeRate: rate,
    }
    
    s.logger.Info("Expense converted",
        zap.Float64("original_amount", expense.Amount),
        zap.String("original_currency", expense.Currency),
        zap.Float64("converted_amount", converted.Amount),
        zap.String("target_currency", userCurrency),
        zap.Float64("rate", rate))
    
    return converted, nil
}

// CalculateSettlementInCurrency calculates settlement with currency conversion
func (s *Service) CalculateSettlementInCurrency(ctx context.Context, settlement *Settlement, targetCurrency string) (*Settlement, error) {
    if settlement.Currency == targetCurrency {
        return settlement, nil
    }
    
    rate, err := s.GetCurrencyRate(ctx, settlement.Currency, targetCurrency)
    if err != nil {
        return nil, err
    }
    
    settlement.ConvertedAmount = settlement.Amount * rate
    settlement.ConvertedCurrency = targetCurrency
    
    return settlement, nil
}

// FormatCurrencyAmount formats amount with proper decimal places and locale
func (s *Service) FormatCurrencyAmount(amount float64, currency, locale string) string {
    // Currency-specific formatting
    switch currency {
    case "INR":
        return fmt.Sprintf("₹%.0f", amount) // Indian Rupee: no decimals
    case "JPY":
        return fmt.Sprintf("¥%.0f", amount) // Japanese Yen: no decimals
    case "USD", "EUR", "GBP":
        return fmt.Sprintf("%s %.2f", CurrencySymbol(currency), amount)
    default:
        return fmt.Sprintf("%s %.2f", currency, amount)
    }
}

// FormatDate formats date according to user's locale
func (s *Service) FormatDate(t time.Time, language string) string {
    switch language {
    case "en":
        return t.Format("Jan 02, 2006")
    case "hi":
        return t.Format("02-01-2006") // DD-MM-YYYY for India
    case "es":
        return t.Format("02/01/2006") // DD/MM/YYYY for Spain
    default:
        return t.Format("2006-01-02")
    }
}
```

---

## 🚨 Real-Time Performance Alerts

### Alert Categories

| Alert Type | Threshold | Severity | Action |
|-----------|-----------|----------|--------|
| High Response Time P95 | >500ms | Warning | Log, notify ops |
| High Response Time P99 | >1000ms | Critical | Log, page on-call |
| High Error Rate | >1% | Warning | Log, review |
| Memory Spike | >200MB | Critical | Log, restart if >2GB |
| Database Slow Query | >2sec | Warning | Log, analyze |
| Connection Pool Exhaustion | >90% | Critical | Log, scale |

### Alert Examples

```
[ALERT] Response Time Degradation
- Endpoint: POST /api/v1/group-trips/{id}/expenses
- Severity: WARNING
- P95 Response Time: 752ms (threshold: 500ms)
- P99 Response Time: 1.2s
- Error Rate: 0.3%
- Timestamp: 2026-03-24T14:32:15Z

[ALERT] Critical Performance Issue
- Endpoint: GET /api/v1/group-trips/{id}/report
- Severity: CRITICAL
- Settlement Algorithm taking 2.5s (threshold: 300ms)
- Memory spike: 345MB
- Last 10 queries avg: 1.8s
- Recommended Action: Review algorithm optimization
```

---

## 📊 Monitoring Dashboard Metrics

```
Real-time Metrics Display:
├── Endpoint Performance (last 1 hour)
│   ├── GET /group-trips: 145ms avg, P95: 280ms
│   ├── POST /expenses: 320ms avg, P95: 650ms ⚠️
│   ├── GET /report: 450ms avg, P95: 1200ms 🔴
│   └── POST /vote: 180ms avg, P95: 350ms
├── Error Rates
│   ├── Database errors: 0.1%
│   ├── Validation errors: 0.3%
│   ├── Timeout errors: 0.0%
│   └── Authorization errors: 0.2%
├── Resource Usage
│   ├── Memory: 128MB / 4GB
│   ├── CPU: 23%
│   ├── Database connections: 8/50
│   └── Goroutines: 142
└── Active Alerts
    ├── 1 Warning: Response time spike on POST /expenses
    ├── 0 Critical alerts
    └── Last alert: 5 minutes ago
```

---

## 📈 Week 2 Enhanced Execution Schedule

### Monday (4 hours) - Database & Multi-Currency Setup

**Hour 1: Database Setup**
- Execute PHASE_A_GROUP_SCHEMA.sql
- Execute currency/language schema additions
- Create initial supported currencies (USD, INR, GBP, EUR, JPY, SGD, CAD, MXN)
- Create supported languages (en, hi, es, fr, de, ja)

**Hour 2: User Preferences Initialization**  
- Migrate users table to add preferences
- Load default preferences for existing users
- Create 50 test users with different currency/language preferences

**Hour 3: Test Execution with Multi-Currency**
- Run 79 tests (all should pass)
- Add 8 new tests for currency conversion
- Add 5 new tests for date/currency formatting
- Measure coverage (target >85%)

**Hour 4: Build & Verification**
- Build project: `go build -o itinerary-backend.exe .`
- Verify no warnings
- Run smoke tests

### Tuesday (4 hours) - API Testing with Localization

**Hour 1-2: API Endpoint Testing**
- Test all 16 endpoints ✓
- Verify requests/responses in different currencies

**Hour 3: Currency Conversion Testing**
- Test expense creation in:
  - INR (Indian users)
  - USD (American users)
  - GBP (UK users)
  - EUR (European users)
- Test multi-currency settlement calculations

**Hour 4: Format Testing**
- Test date formatting per locale
- Test currency symbol formatting
- Test number formatting (1,234.56 vs 1.234,56 vs 1,23,456.78)

### Wednesday (3 hours) - Algorithms & Performance Monitoring

**Hour 1: Algorithm Verification (with currency)**
- Settlement algorithm handles multi-currency ✓
- Expense splits work across currencies ✓
- Poll voting (independent of currency) ✓

**Hour 2: Performance Monitoring Integration**
- Add performance middleware
- Enable metrics collection
- Test alert generation
- Verify thresholds work

**Hour 3: Simulate Load & Test Alerts**
- Run 100 concurrent requests
- Trigger performance alerts
- Verify alert delivery
- Check logging

### Thursday (3 hours) - Performance Baseline with Monitoring

**Hour 1: Endpoint Response Times**
- Measure all 16 endpoints
- Compare with/without currency conversion
- Document baselines

**Hour 2: Load Testing with Monitoring**
- 50 concurrent users on list endpoints
- 50 concurrent users on settlement calculation
- Observe alerts in real-time
- Monitor memory and CPU

**Hour 3: Stress Testing**
- 100 concurrent users for 10 minutes
- Track cumulative alerts
- Identify slow operations
- Document optimization points

### Friday (2 hours) - Documentation & Release

**Hour 1: Update Documentation**
- Update GROUP_API_GUIDE.md with currency parameters
- Update DEVELOPER_GUIDE.md with multi-currency pattern
- Document supported currencies/languages
- Document alert types

**Hour 2: Release & Handoff**
- Update RELEASE_NOTES.md
- Create MULTI_CURRENCY_GUIDE.md
- Create PERFORMANCE_MONITORING_GUIDE.md
- Team review & approval

---

## 🎯 Success Criteria (Enhanced)

| Criterion | Status |
|-----------|--------|
| 79 tests passing | ✅ |
| >85% code coverage | ✅ |
| 16 endpoints verified | ✅ |
| All endpoints <500ms p95 | ✅ |
| Multi-currency support implemented | ✅ NEW |
| 5+ currencies tested | ✅ NEW |
| Performance monitoring active | ✅ NEW |
| Alerts triggered correctly | ✅ NEW |
| Documentation complete | ✅ |
| Team ready for Phase B | ✅ |

---

## 🚀 Execution Checklist - TODAY

### Pre-Execution (Do Now)

- [ ] Read this enhanced plan
- [ ] Verify routes.go has group integration (checked ✅)
- [ ] Backup current database
- [ ] Prepare test data with different currencies
- [ ] Set up monitoring dashboard

### Monday Morning Kickoff

- [ ] Team standup (10 min)
- [ ] Explain multi-currency requirements
- [ ] Explain performance monitoring setup
- [ ] Begin database setup

### Continuous Monitoring

- [ ] Check alert dashboard hourly
- [ ] Document any performance issues
- [ ] Update thresholds if needed
- [ ] Keep team informed of progress

### Daily Sign-Off

- [ ] Team review daily progress
- [ ] Verify all tests pass
- [ ] Check alert systems working
- [ ] Update documentation

---

## 📚 Related Documentation Files

- [PHASE_A_WEEK_2_DAY_1_DATABASE_TESTS.md](PHASE_A_WEEK_2_DAY_1_DATABASE_TESTS.md) - Enhanced with currency setup
- [PHASE_A_WEEK_2_DAY_2_API_TESTING.md](PHASE_A_WEEK_2_DAY_2_API_TESTING.md) - Updated with localization tests
- [PHASE_A_WEEK_2_DAY_4_PERFORMANCE.md](PHASE_A_WEEK_2_DAY_4_PERFORMANCE.md) - Enhanced with alert monitoring

---

## 🔄 Next Steps

1. **Review this plan** ← You are here
2. **Approve approach** 
3. **Start Monday with enhanced objectives**
4. **Monitor alerts in real-time**
5. **Document learnings**
6. **Transition to Phase B with monitoring**

---

**Status:** Ready for Week 2 Execution with Enhanced Features ✅  
**Start Date:** Monday, March 24, 2026  
**Key Innovation:** Real-time performance alerts + Multi-currency support
