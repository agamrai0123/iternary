-- Phase A Week 2: Multi-Currency Schema
-- Purpose: Add multi-currency support and performance monitoring
-- Created: 2026-03-24
-- Status: Ready to apply

-- ============================================================================
-- MULTI-CURRENCY SUPPORT TABLES
-- ============================================================================

-- User Preferences: Store nationality, preferred currency, language, timezone
CREATE TABLE IF NOT EXISTS user_preferences (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL UNIQUE,
    nationality TEXT NOT NULL,                -- ISO 3166-1 (e.g., 'US', 'IN', 'GB', 'JP', 'DE')
    preferred_currency TEXT NOT NULL,         -- ISO 4217 (e.g., 'USD', 'INR', 'GBP', 'JPY', 'EUR')
    preferred_language TEXT NOT NULL,         -- ISO 639-1 (e.g., 'en', 'hi', 'es', 'fr', 'de', 'ja', 'pt', 'zh')
    timezone TEXT NOT NULL,                   -- IANA timezone (e.g., 'America/New_York', 'Asia/Kolkata', 'Europe/London')
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Supported Currencies: Reference list of all supported currencies
CREATE TABLE IF NOT EXISTS supported_currencies (
    code TEXT PRIMARY KEY,                    -- ISO 4217 code (e.g., 'USD', 'INR', 'GBP', 'JPY', 'EUR')
    name TEXT NOT NULL,                       -- Full currency name
    symbol TEXT NOT NULL,                     -- Currency symbol (e.g., '$', '₹', '£', '¥', '€')
    decimal_places INTEGER NOT NULL DEFAULT 2, -- Number of decimal places (e.g., JPY = 0, others = 2)
    exchange_rate_to_usd DECIMAL(18,6) NOT NULL DEFAULT 1.0, -- Exchange rate to USD
    is_active BOOLEAN NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Supported Languages: Reference list of supported languages
CREATE TABLE IF NOT EXISTS supported_languages (
    code TEXT PRIMARY KEY,                    -- ISO 639-1 code (e.g., 'en', 'hi', 'es', 'fr', 'de', 'ja', 'pt', 'zh')
    name TEXT NOT NULL,                       -- Full language name
    native_name TEXT NOT NULL,                -- Language name in native script
    is_active BOOLEAN NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- MULTI-CURRENCY TRANSACTION TRACKING
-- ============================================================================

-- Expense Conversions: Track currency conversion details for each expense
CREATE TABLE IF NOT EXISTS expense_conversions (
    id TEXT PRIMARY KEY,
    expense_id TEXT NOT NULL,
    original_currency TEXT NOT NULL,          -- Currency in which expense was recorded
    original_amount DECIMAL(18,2) NOT NULL,   -- Amount in original currency
    converted_currency TEXT NOT NULL,         -- Currency after conversion (usually USD)
    converted_amount DECIMAL(18,2) NOT NULL,  -- Amount after conversion
    exchange_rate DECIMAL(18,6) NOT NULL,     -- Exchange rate used
    conversion_date TIMESTAMP NOT NULL,       -- When conversion happened
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (original_currency) REFERENCES supported_currencies(code),
    FOREIGN KEY (converted_currency) REFERENCES supported_currencies(code)
);

-- Settlement Details: Track settlement calculations with multi-currency
CREATE TABLE IF NOT EXISTS settlement_details (
    id TEXT PRIMARY KEY,
    settlement_id TEXT NOT NULL,
    from_user TEXT NOT NULL,
    to_user TEXT NOT NULL,
    from_user_currency TEXT NOT NULL,        -- Original currency of from_user
    to_user_currency TEXT NOT NULL,          -- Target currency of to_user
    original_amount DECIMAL(18,2) NOT NULL,  -- Amount in from_user's currency
    converted_amount DECIMAL(18,2) NOT NULL, -- Amount converted to to_user's currency
    exchange_rate DECIMAL(18,6) NOT NULL,    -- Exchange rate applied
    conversion_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (from_user_currency) REFERENCES supported_currencies(code),
    FOREIGN KEY (to_user_currency) REFERENCES supported_currencies(code)
);

-- ============================================================================
-- PERFORMANCE MONITORING TABLES
-- ============================================================================

-- Performance Metrics: Real-time endpoint performance data
CREATE TABLE IF NOT EXISTS performance_metrics (
    id TEXT PRIMARY KEY,
    endpoint_path TEXT NOT NULL,              -- API endpoint path (e.g., '/api/v1/group-trips')
    method TEXT NOT NULL,                     -- HTTP method (GET, POST, etc.)
    response_time_ms INTEGER NOT NULL,        -- Response time in milliseconds
    status_code INTEGER NOT NULL,             -- HTTP status code
    memory_delta_mb INTEGER,                  -- Memory delta in MB
    goroutines_count INTEGER,
    error_message TEXT,                       -- If status_code >= 400
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Performance Alerts: Alert records for threshold violations
CREATE TABLE IF NOT EXISTS performance_alerts (
    id TEXT PRIMARY KEY,
    alert_type TEXT NOT NULL,                 -- Type: high_response_time, error_rate, memory_spike, etc.
    endpoint_path TEXT NOT NULL,
    method TEXT NOT NULL,
    severity TEXT NOT NULL,                   -- info, warning, critical
    threshold_value DECIMAL(18,4) NOT NULL,
    current_value DECIMAL(18,4) NOT NULL,
    message TEXT NOT NULL,
    recommendation TEXT,
    is_acknowledged BOOLEAN DEFAULT 0,
    acknowledged_by TEXT,
    acknowledged_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Alert Rules: Configurable thresholds for triggering alerts
CREATE TABLE IF NOT EXISTS alert_rules (
    id TEXT PRIMARY KEY,
    alert_type TEXT NOT NULL UNIQUE,
    description TEXT,
    threshold_value DECIMAL(18,4) NOT NULL,
    severity TEXT NOT NULL,
    is_active BOOLEAN DEFAULT 1,
    cooldown_seconds INTEGER DEFAULT 300,     -- Don't repeat same alert for X seconds
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Monitoring Settings: Global performance monitoring configuration
CREATE TABLE IF NOT EXISTS monitoring_settings (
    id TEXT PRIMARY KEY,
    key TEXT NOT NULL UNIQUE,
    value TEXT NOT NULL,
    description TEXT,
    setting_type TEXT,                        -- integer, decimal, text, boolean
    is_active BOOLEAN DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Hourly Performance Stats: Aggregated statistics per hour
CREATE TABLE IF NOT EXISTS hourly_performance_stats (
    id TEXT PRIMARY KEY,
    hour_timestamp TIMESTAMP NOT NULL,
    endpoint_path TEXT NOT NULL,
    method TEXT NOT NULL,
    total_requests INTEGER,
    successful_requests INTEGER,
    failed_requests INTEGER,
    avg_response_time_ms DECIMAL(18,2),
    p50_response_time_ms INTEGER,
    p95_response_time_ms INTEGER,
    p99_response_time_ms INTEGER,
    max_response_time_ms INTEGER,
    error_rate DECIMAL(18,4),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- INDEXES FOR PERFORMANCE
-- ============================================================================

-- User preferences indexes
CREATE INDEX IF NOT EXISTS idx_user_preferences_user_id ON user_preferences(user_id);
CREATE INDEX IF NOT EXISTS idx_user_preferences_currency ON user_preferences(preferred_currency);
CREATE INDEX IF NOT EXISTS idx_user_preferences_language ON user_preferences(preferred_language);

-- Performance metrics indexes
CREATE INDEX IF NOT EXISTS idx_performance_metrics_endpoint ON performance_metrics(endpoint_path, method);
CREATE INDEX IF NOT EXISTS idx_performance_metrics_status ON performance_metrics(status_code);
CREATE INDEX IF NOT EXISTS idx_performance_metrics_timestamp ON performance_metrics(timestamp);

-- Performance alerts indexes
CREATE INDEX IF NOT EXISTS idx_performance_alerts_type ON performance_alerts(alert_type);
CREATE INDEX IF NOT EXISTS idx_performance_alerts_severity ON performance_alerts(severity);
CREATE INDEX IF NOT EXISTS idx_performance_alerts_endpoint ON performance_alerts(endpoint_path, method);
CREATE INDEX IF NOT EXISTS idx_performance_alerts_timestamp ON performance_alerts(created_at);

-- ============================================================================
-- REPORTING VIEWS
-- ============================================================================

-- Current Performance Status View
CREATE VIEW IF NOT EXISTS vw_current_performance_status AS
SELECT 
    endpoint_path,
    method,
    COUNT(*) as total_requests,
    SUM(CASE WHEN status_code < 400 THEN 1 ELSE 0 END) as successful_requests,
    SUM(CASE WHEN status_code >= 400 THEN 1 ELSE 0 END) as failed_requests,
    CAST(AVG(response_time_ms) AS DECIMAL(18,2)) as avg_response_time_ms,
    MAX(response_time_ms) as max_response_time_ms,
    MIN(response_time_ms) as min_response_time_ms,
    CAST(SUM(CASE WHEN status_code >= 400 THEN 1 ELSE 0 END) * 100.0 / COUNT(*) AS DECIMAL(18,2)) as error_rate_percent,
    MAX(timestamp) as last_request
FROM performance_metrics
WHERE timestamp > datetime('now', '-1 hour')
GROUP BY endpoint_path, method
ORDER BY last_request DESC;

-- Active Alerts View  
CREATE VIEW IF NOT EXISTS vw_active_alerts AS
SELECT 
    id,
    alert_type,
    endpoint_path,
    method,
    severity,
    current_value,
    threshold_value,
    message,
    recommendation,
    CASE 
        WHEN severity = 'critical' THEN 'URGENT'
        WHEN severity = 'warning' THEN 'CHECK'
        ELSE 'INFO'
    END as priority,
    created_at
FROM performance_alerts
WHERE is_acknowledged = 0
    AND created_at > datetime('now', '-24 hours')
ORDER BY 
    CASE WHEN severity = 'critical' THEN 1 WHEN severity = 'warning' THEN 2 ELSE 3 END,
    created_at DESC;

-- 24-Hour Performance Trends View
CREATE VIEW IF NOT EXISTS vw_performance_trends_24h AS
SELECT 
    endpoint_path,
    method,
    strftime('%Y-%m-%d %H:00', timestamp) as hour,
    COUNT(*) as requests_count,
    CAST(AVG(response_time_ms) AS DECIMAL(18,2)) as avg_response_time_ms,
    MAX(response_time_ms) as max_response_time_ms,
    SUM(CASE WHEN status_code >= 400 THEN 1 ELSE 0 END) as error_count
FROM performance_metrics
WHERE timestamp > datetime('now', '-24 hours')
GROUP BY endpoint_path, method, hour
ORDER BY hour DESC, endpoint_path, method;

-- ============================================================================
-- SAMPLE DATA INSERTION
-- ============================================================================

-- Insert Supported Currencies
INSERT OR IGNORE INTO supported_currencies (code, name, symbol, decimal_places, exchange_rate_to_usd) VALUES
('USD', 'United States Dollar', '$', 2, 1.000000),
('EUR', 'Euro', '€', 2, 1.084000),
('INR', 'Indian Rupee', '₹', 2, 0.012038),
('GBP', 'British Pound', '£', 2, 1.255000),
('JPY', 'Japanese Yen', '¥', 0, 0.006700),
('SGD', 'Singapore Dollar', 'SGD', 2, 0.744000),
('CAD', 'Canadian Dollar', 'C$', 2, 0.737000),
('MXN', 'Mexican Peso', '$', 2, 0.058000);

-- Insert Supported Languages
INSERT OR IGNORE INTO supported_languages (code, name, native_name) VALUES
('en', 'English', 'English'),
('es', 'Spanish', 'Español'),
('fr', 'French', 'Français'),
('de', 'German', 'Deutsch'),
('hi', 'Hindi', 'हिन्दी'),
('ja', 'Japanese', '日本語'),
('pt', 'Portuguese', 'Português'),
('zh', 'Chinese', '中文');

-- Insert Alert Rules
INSERT OR IGNORE INTO alert_rules (id, alert_type, description, threshold_value, severity, cooldown_seconds) VALUES
('rule-001', 'high_response_time_p95', 'P95 response time exceeds 500ms', 500, 'warning', 300),
('rule-002', 'high_response_time_p99', 'P99 response time exceeds 1000ms', 1000, 'critical', 300),
('rule-003', 'high_error_rate', 'Error rate exceeds 1%', 1.0, 'warning', 600),
('rule-004', 'memory_spike', 'Memory delta exceeds 200MB', 200, 'warning', 900),
('rule-005', 'db_connection_pool', 'DB connections exceed 80% of pool', 80, 'warning', 300),
('rule-006', 'slow_query', 'Individual query exceeds 5000ms', 5000, 'critical', 600);

-- Insert Monitoring Settings
INSERT OR IGNORE INTO monitoring_settings (id, key, value, description, setting_type) VALUES
('setting-001', 'p95_threshold_ms', '500', 'P95 response time threshold in milliseconds', 'integer'),
('setting-002', 'p99_threshold_ms', '1000', 'P99 response time threshold in milliseconds', 'integer'),
('setting-003', 'error_rate_threshold', '1.0', 'Error rate threshold in percent', 'decimal'),
('setting-004', 'memory_threshold_mb', '200', 'Memory delta threshold in MB', 'integer'),
('setting-005', 'alert_cooldown_seconds', '300', 'Cooldown period between same alerts in seconds', 'integer'),
('setting-006', 'metrics_retention_days', '30', 'Number of days to retain metrics data', 'integer'),
('setting-007', 'sampling_rate', '100', 'Percentage of requests to sample for metrics', 'integer'),
('setting-008', 'enable_real_time_alerts', 'true', 'Enable real-time alerting system', 'boolean'),
('setting-009', 'enable_performance_dashboard', 'true', 'Enable performance dashboard endpoint', 'boolean'),
('setting-010', 'timezone_utc_offset', '0', 'UTC offset for this deployment in hours', 'integer');

-- ============================================================================
-- MIGRATION NOTES FOR OTHER DATABASES
-- ============================================================================
-- Oracle:
--   - Replace TIMESTAMP DEFAULT CURRENT_TIMESTAMP with TIMESTAMP DEFAULT SYSTIMESTAMP
--   - Replace INTEGER with NUMBER
--   - Replace TEXT with VARCHAR2(max_length)
--   - Use CREATE TABLE statement with CREATE TABLE IF NOT EXISTS replacement
--
-- PostgreSQL:
--   - Replace TIMESTAMP DEFAULT CURRENT_TIMESTAMP with NOW()
--   - Replace INTEGER with INTEGER
--   - Replace BOOLEAN DEFAULT 1 with BOOLEAN DEFAULT TRUE
--   - Replace BOOLEAN DEFAULT 0 with BOOLEAN DEFAULT FALSE
--   - Use CREATE TABLE IF NOT EXISTS (already compatible)
--   - For JSON support, consider JSONB columns instead of TEXT fields
-- ============================================================================
