-- MFA Configuration Table
CREATE TABLE IF NOT EXISTS mfa_configs (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    mfa_type TEXT NOT NULL CHECK(mfa_type IN ('totp', 'sms', 'email')),
    enabled BOOLEAN NOT NULL DEFAULT FALSE,
    secret_hash TEXT NOT NULL, -- Encrypted using AES-256
    backup_codes TEXT NOT NULL, -- Encrypted JSON array of hashes
    created_at TIMESTAMP NOT NULL,
    verified_at TIMESTAMP,
    last_used_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, mfa_type)
);

-- MFA Verification Attempts
CREATE TABLE IF NOT EXISTS mfa_attempts (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    mfa_type TEXT NOT NULL,
    code_type TEXT NOT NULL CHECK(code_type IN ('totp', 'backup')),
    success BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    ip_address TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Backup Code Usage
CREATE TABLE IF NOT EXISTS backup_code_usage (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    used_at TIMESTAMP NOT NULL,
    ip_address TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_mfa_configs_user_id ON mfa_configs(user_id);
CREATE INDEX IF NOT EXISTS idx_mfa_attempts_user_id ON mfa_attempts(user_id);
CREATE INDEX IF NOT EXISTS idx_mfa_attempts_created ON mfa_attempts(created_at);
CREATE INDEX IF NOT EXISTS idx_backup_usage_user_id ON backup_code_usage(user_id);

-- OAuth Linked Accounts Table
CREATE TABLE IF NOT EXISTS linked_accounts (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    provider TEXT NOT NULL CHECK(provider IN ('github', 'google', 'microsoft')),
    provider_id TEXT NOT NULL,
    email TEXT,
    name TEXT,
    avatar TEXT,
    linked_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(provider, provider_id)
);

-- OAuth State Table (for CSRF protection)
CREATE TABLE IF NOT EXISTS oauth_states (
    state TEXT PRIMARY KEY,
    provider TEXT NOT NULL,
    user_id TEXT,
    redirect_uri TEXT,
    created_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL
);

-- Indexes for OAuth
CREATE INDEX IF NOT EXISTS idx_linked_accounts_user_id ON linked_accounts(user_id);
CREATE INDEX IF NOT EXISTS idx_oauth_states_created ON oauth_states(created_at);
