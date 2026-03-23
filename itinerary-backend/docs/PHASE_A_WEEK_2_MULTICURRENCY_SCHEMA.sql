-- Phase A Week 2 Enhanced Schema: Multi-Currency & Language Support
-- This schema extends PHASE_A_GROUP_SCHEMA.sql with currency and performance monitoring

-- ==================== MULTI-CURRENCY & LOCALIZATION ====================

-- User Preferences Table
CREATE TABLE user_preferences (
    user_id VARCHAR(255) PRIMARY KEY,
    nationality VARCHAR(10),                 -- ISO 3166-1: US, IN, GB, FR, DE, JP, SG, CA, MX, etc
    preferred_currency VARCHAR(3) NOT NULL DEFAULT 'USD',  -- ISO 4217: USD, INR, GBP, EUR, JPY, SGD, CAD, MXN
    preferred_language VARCHAR(5) NOT NULL DEFAULT 'en',   -- ISO 639-1: en, hi, es, fr, de, ja
    timezone VARCHAR(50),                   -- IANA timezone: America/New_York, Asia/Kolkata, Europe/London
    date_format VARCHAR(20) DEFAULT 'YYYY-MM-DD',
    number_format VARCHAR(20) DEFAULT 'US',  -- US (1,234.56) vs EU (1.234,56) vs IN (1,23,456.78)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_nationality (nationality),
    INDEX idx_currency (preferred_currency),
    INDEX idx_language (preferred_language)
);

-- Supported Currencies Table
CREATE TABLE supported_currencies (
    code VARCHAR(3) PRIMARY KEY,
    symbol VARCHAR(5) NOT NULL,
    name VARCHAR(100) NOT NULL,
    decimal_places INT DEFAULT 2,           -- USD/EUR: 2, JPY: 0, INR: 2
    exchange_rate_to_usd FLOAT NOT NULL DEFAULT 1.0,
    is_active BOOLEAN DEFAULT TRUE,
    supported_countries TEXT,               -- Comma-separated ISO codes
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_active (is_active)
);

-- Supported Languages Table
CREATE TABLE supported_languages (
    code VARCHAR(5) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    native_name VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    region_codes TEXT,                      -- Supported regions
    rtl BOOLEAN DEFAULT FALSE,              -- Right-to-left (Arabic, Hebrew, etc)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_active (is_active)
);

-- ==================== MULTI-CURRENCY EXPENSES ====================

-- Update existing expenses table structure (these are the new columns to ADD)
-- ALTER TABLE expenses ADD COLUMN currency VARCHAR(3) DEFAULT 'USD' AFTER amount;
-- ALTER TABLE expenses ADD COLUMN exchange_rate FLOAT DEFAULT 1.0 AFTER currency;
-- ALTER TABLE expenses ADD COLUMN original_amount FLOAT AFTER exchange_rate;
-- ALTER TABLE expenses ADD COLUMN original_currency VARCHAR(3) AFTER original_amount;

-- Expense Currency Conversion Log (audit trail)
CREATE TABLE expense_conversions (
    id VARCHAR(255) PRIMARY KEY,
    expense_id VARCHAR(255) NOT NULL,
    from_currency VARCHAR(3) NOT NULL,
    to_currency VARCHAR(3) NOT NULL,
    from_amount FLOAT NOT NULL,
    to_amount FLOAT NOT NULL,
    exchange_rate FLOAT NOT NULL,
    converted_by_user_id VARCHAR(255),
    converted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (expense_id) REFERENCES expenses(id) ON DELETE CASCADE,
    FOREIGN KEY (converted_by_user_id) REFERENCES users(id),
    INDEX idx_expense (expense_id),
    INDEX idx_currencies (from_currency, to_currency),
    INDEX idx_converted_at (converted_at DESC)
);

-- ==================== MULTI-CURRENCY SETTLEMENTS ====================

-- Update existing settlements table structure (these are the new columns to ADD)
-- ALTER TABLE settlements ADD COLUMN currency VARCHAR(3) DEFAULT 'USD' AFTER amount;
-- ALTER TABLE settlements ADD COLUMN converted_amount FLOAT AFTER currency;
-- ALTER TABLE settlements ADD COLUMN converted_currency VARCHAR(3) AFTER converted_amount;

-- Settlement Currency Details
CREATE TABLE settlement_details (
    id VARCHAR(255) PRIMARY KEY,
    settlement_id VARCHAR(255) NOT NULL UNIQUE,
    creditor_currency VARCHAR(3) NOT NULL,
    debtor_currency VARCHAR(3) NOT NULL,
    exchange_rate FLOAT NOT NULL,
    amount_in_creditor_currency FLOAT NOT NULL,
    amount_in_debtor_currency FLOAT NOT NULL,
    calculated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (settlement_id) REFERENCES settlements(id) ON DELETE CASCADE,
    INDEX idx_settlement (settlement_id),
    INDEX idx_currencies (creditor_currency, debtor_currency)
);

-- ==================== PERFORMANCE MONITORING ====================

-- Performance Metrics Table (lightweight - stores aggregates)
CREATE TABLE performance_metrics (
    id VARCHAR(255) PRIMARY KEY,
    endpoint_path VARCHAR(255) NOT NULL,
    method VARCHAR(10) NOT NULL,           -- GET, POST, PUT, DELETE
    response_time_ms BIGINT NOT NULL,       -- Response time in milliseconds
    status_code INT NOT NULL,
    user_id VARCHAR(255),
    trip_id VARCHAR(255),
    error_flag BOOLEAN DEFAULT FALSE,
    error_message VARCHAR(500),
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    INDEX idx_endpoint_path (endpoint_path),
    INDEX idx_recorded_at (recorded_at DESC),
    INDEX idx_endpoint_time (endpoint_path, recorded_at DESC),
    INDEX idx_status_code (status_code),
    INDEX idx_error_flag (error_flag)
);

-- Endpoint Performance Aggregates (pre-calculated for dashboard)
CREATE TABLE endpoint_performance_aggregates (
    id VARCHAR(255) PRIMARY KEY,
    endpoint_path VARCHAR(255) NOT NULL UNIQUE,
    method VARCHAR(10) NOT NULL,
    total_requests BIGINT DEFAULT 0,
    successful_requests BIGINT DEFAULT 0,
    failed_requests BIGINT DEFAULT 0,
    error_rate FLOAT DEFAULT 0,
    avg_response_time_ms FLOAT DEFAULT 0,
    p50_response_time_ms BIGINT DEFAULT 0,
    p95_response_time_ms BIGINT DEFAULT 0,    -- 95th percentile
    p99_response_time_ms BIGINT DEFAULT 0,    -- 99th percentile
    max_response_time_ms BIGINT DEFAULT 0,
    min_response_time_ms BIGINT DEFAULT 0,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_window VARCHAR(20),               -- 1min, 5min, 1hour, 1day
    INDEX idx_endpoint_path (endpoint_path),
    INDEX idx_last_updated (last_updated DESC)
);

-- Performance Alerts Table
CREATE TABLE performance_alerts (
    id VARCHAR(255) PRIMARY KEY,
    alert_type VARCHAR(50) NOT NULL,        -- high_response_time, high_error_rate, memory_spike, etc
    endpoint_path VARCHAR(255),
    method VARCHAR(10),
    severity VARCHAR(20) NOT NULL,           -- info, warning, critical
    threshold_value FLOAT NOT NULL,
    current_value FLOAT NOT NULL,
    message TEXT NOT NULL,
    recommendation TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP NULL,
    resolved_by_user_id VARCHAR(255),
    FOREIGN KEY (resolved_by_user_id) REFERENCES users(id),
    INDEX idx_severity (severity),
    INDEX idx_endpoint_path (endpoint_path),
    INDEX idx_created_at (created_at DESC),
    INDEX idx_resolved_at (resolved_at),
    INDEX idx_alert_type (alert_type)
);

-- Performance Alert Rules (define thresholds)
CREATE TABLE alert_rules (
    id VARCHAR(255) PRIMARY KEY,
    alert_type VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    condition_formula VARCHAR(255),
    threshold_warning FLOAT,
    threshold_critical FLOAT,
    evaluation_window_seconds INT DEFAULT 60,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_active (is_active)
);

-- ==================== PERFORMANCE TREND ANALYSIS ====================

-- Hourly Performance Summary (for trend analysis and dashboards)
CREATE TABLE hourly_performance_stats (
    id VARCHAR(255) PRIMARY KEY,
    endpoint_path VARCHAR(255) NOT NULL,
    method VARCHAR(10) NOT NULL,
    hour_bucket TIMESTAMP NOT NULL,          -- Start of hour
    request_count INT,
    error_count INT,
    avg_response_time_ms FLOAT,
    p95_response_time_ms BIGINT,
    p99_response_time_ms BIGINT,
    max_response_time_ms BIGINT,
    memory_usage_mb INT,
    cpu_usage_percent FLOAT,
    database_connections INT,
    active_goroutines INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY idx_endpoint_hour (endpoint_path, method, hour_bucket),
    INDEX idx_hour_bucket (hour_bucket DESC),
    INDEX idx_endpoint (endpoint_path)
);

-- ==================== MONITORING CONFIGURATION ====================

-- Monitoring Settings
CREATE TABLE monitoring_settings (
    setting_key VARCHAR(100) PRIMARY KEY,
    setting_value VARCHAR(1000),
    data_type VARCHAR(20),                  -- string, integer, float, boolean, json
    description TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(255)
);

-- Default monitoring settings to insert:
-- monitoring_settings:
-- - p95_threshold_ms: 500 (warning if p95 > 500ms)
-- - p99_threshold_ms: 1000 (critical if p99 > 1000ms)
-- - error_rate_threshold: 0.01 (1% - warning)
-- - memory_threshold_mb: 200 (warning if > 200MB)
-- - alert_cooldown_seconds: 300 (don't spam same alert more than every 5 min)
-- - retention_days: 90 (keep metrics for 90 days)

-- ==================== INDEXED VIEWS FOR REPORTING ====================

-- View: Current Performance Status
CREATE VIEW vw_current_performance_status AS
SELECT 
    epa.endpoint_path,
    epa.method,
    epa.total_requests,
    epa.error_rate,
    epa.avg_response_time_ms,
    epa.p95_response_time_ms,
    epa.p99_response_time_ms,
    CASE 
        WHEN epa.p99_response_time_ms > 1000 THEN 'CRITICAL'
        WHEN epa.p95_response_time_ms > 500 THEN 'WARNING'
        ELSE 'HEALTHY'
    END as health_status,
    epa.last_updated
FROM endpoint_performance_aggregates epa
ORDER BY epa.last_updated DESC;

-- View: Active Alerts
CREATE VIEW vw_active_alerts AS
SELECT 
    pa.id,
    pa.alert_type,
    pa.endpoint_path,
    pa.severity,
    pa.threshold_value,
    pa.current_value,
    pa.message,
    pa.recommendation,
    pa.created_at,
    TIMESTAMPDIFF(SECOND, pa.created_at, NOW()) as alert_age_seconds
FROM performance_alerts pa
WHERE pa.resolved_at IS NULL
ORDER BY 
    CASE pa.severity
        WHEN 'critical' THEN 1
        WHEN 'warning' THEN 2
        ELSE 3
    END,
    pa.created_at DESC;

-- View: Performance Trends (24-hour)
CREATE VIEW vw_performance_trends_24h AS
SELECT 
    hps.endpoint_path,
    hps.method,
    hps.hour_bucket,
    hps.request_count,
    hps.error_count,
    hps.avg_response_time_ms,
    hps.p95_response_time_ms,
    hps.max_response_time_ms,
    hps.memory_usage_mb,
    hps.cpu_usage_percent
FROM hourly_performance_stats hps
WHERE hps.hour_bucket >= DATE_SUB(NOW(), INTERVAL 24 HOUR)
ORDER BY hps.endpoint_path, hps.hour_bucket DESC;

-- ==================== SAMPLE DATA ====================

-- Insert sample supported currencies
INSERT IGNORE INTO supported_currencies (code, symbol, name, decimal_places, exchange_rate_to_usd, is_active, supported_countries) VALUES
('USD', '$', 'US Dollar', 2, 1.0, TRUE, 'US'),
('EUR', '€', 'Euro', 2, 1.10, TRUE, 'DE,FR,ES,IT,NL,BE,AT,IE,PT,GR,CY,LU,MT,SK,SI'),
('INR', '₹', 'Indian Rupee', 2, 0.012, TRUE, 'IN'),
('GBP', '£', 'British Pound', 2, 1.27, TRUE, 'GB'),
('JPY', '¥', 'Japanese Yen', 0, 0.0067, TRUE, 'JP'),
('SGD', 'S$', 'Singapore Dollar', 2, 0.75, TRUE, 'SG'),
('CAD', 'C$', 'Canadian Dollar', 2, 0.74, TRUE, 'CA'),
('MXN', '$', 'Mexican Peso', 2, 0.058, TRUE, 'MX'),
('AUD', 'A$', 'Australian Dollar', 2, 0.66, TRUE, 'AU'),
('CHF', 'CHF', 'Swiss Franc', 2, 1.13, TRUE, 'CH');

-- Insert sample supported languages
INSERT IGNORE INTO supported_languages (code, name, native_name, is_active, region_codes, rtl) VALUES
('en', 'English', 'English', TRUE, 'US,GB,CA,AU,NZ,IE', FALSE),
('es', 'Spanish', 'Español', TRUE, 'ES,MX,AR,CO,PE,CL', FALSE),
('fr', 'French', 'Français', TRUE, 'FR,CA,BE,CH,LU,SN', FALSE),
('de', 'German', 'Deutsch', TRUE, 'DE,AT,CH', FALSE),
('it', 'Italian', 'Italiano', TRUE, 'IT,CH', FALSE),
('pt', 'Portuguese', 'Português', TRUE, 'PT,BR', FALSE),
('nl', 'Dutch', 'Nederlands', TRUE, 'NL,BE', FALSE),
('hi', 'Hindi', 'हिन्दी', TRUE, 'IN', FALSE),
('ja', 'Japanese', '日本語', TRUE, 'JP', FALSE),
('zh', 'Chinese', '中文', TRUE, 'CN,TW,HK,SG', FALSE);

-- Insert default alert rules
INSERT IGNORE INTO alert_rules (id, alert_type, description, threshold_warning, threshold_critical, evaluation_window_seconds, is_active) VALUES
('rule-p95-response', 'high_response_time_p95', 'Alert if P95 response time exceeds threshold', 500, 1000, 60, TRUE),
('rule-p99-response', 'high_response_time_p99', 'Alert if P99 response time exceeds threshold', 1000, 2000, 60, TRUE),
('rule-error-rate', 'high_error_rate', 'Alert if error rate exceeds threshold', 0.01, 0.05, 300, TRUE),
('rule-memory', 'memory_spike', 'Alert if memory usage exceeds threshold', 200, 500, 60, TRUE),
('rule-db-conn', 'db_connection_pool', 'Alert if database connection pool usage high', 80, 95, 60, TRUE),
('rule-slow-query', 'slow_database_query', 'Alert if database query exceeds threshold', 2000, 5000, 60, TRUE);

-- Insert default monitoring settings
INSERT IGNORE INTO monitoring_settings (setting_key, setting_value, data_type, description) VALUES
('p95_threshold_ms', '500', 'integer', 'P95 response time warning threshold in milliseconds'),
('p99_threshold_ms', '1000', 'integer', 'P99 response time critical threshold in milliseconds'),
('error_rate_threshold', '0.01', 'float', 'Error rate warning threshold (percentage as decimal)'),
('memory_threshold_mb', '200', 'integer', 'Memory usage warning threshold in MB'),
('alert_cooldown_seconds', '300', 'integer', 'Cooldown period between same alert types (seconds)'),
('metrics_retention_days', '90', 'integer', 'How long to retain metrics data (days)'),
('dashboard_refresh_interval_seconds', '30', 'integer', 'Dashboard auto-refresh interval (seconds)'),
('enable_real_time_alerts', 'true', 'boolean', 'Enable real-time performance alerts'),
('enable_slack_notifications', 'false', 'boolean', 'Send alerts to Slack'),
('enable_pagerduty_integration', 'false', 'boolean', 'Integrate with PagerDuty for critical alerts');

-- ==================== INDEXES FOR PERFORMANCE ====================

-- Ensure optimal query performance
CREATE INDEX idx_pm_endpoint_method_time ON performance_metrics (endpoint_path, method, recorded_at DESC);
CREATE INDEX idx_pm_status_error ON performance_metrics (status_code, error_flag);
CREATE INDEX idx_alert_severity_created ON performance_alerts (severity, created_at DESC);
CREATE INDEX idx_hps_endpoint_hour ON hourly_performance_stats (endpoint_path, hour_bucket DESC);

-- ==================== MIGRATION NOTES ====================

/*
To apply this schema:

1. For Oracle Database:
   - Replace DEFAULT CURRENT_TIMESTAMP with SYSDATE
   - Replace AUTO_INCREMENT with GENERATED ALWAYS AS IDENTITY
   - Replace VARCHAR with appropriate Oracle types
   - Replace TIMESTAMP with DATE or TIMESTAMP datatype

2. For PostgreSQL:
   - Use DEFAULT CURRENT_TIMESTAMP or NOW()
   - Use UUID type for IDs or UUID() function
   - Use SERIAL or BIGSERIAL for auto-increment

3. Apply incrementally:
   a. Create user_preferences table first
   b. Add indexes
   c. Insert supported currencies/languages
   d. Create performance tables
   e. Create views
   f. Test with sample queries

4. Verify data integrity:
   - All users have preferences
   - All currencies have exchange rates
   - All languages are active

5. Migrate existing data:
   - Set default currency based on user's detected location
   - Set default language based on browser settings
   - Set default timezone based on IP geolocation
*/
