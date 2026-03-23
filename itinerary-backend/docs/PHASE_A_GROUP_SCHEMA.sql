-- ============================================================================
-- PHASE A: GROUP COLLABORATION DATABASE SCHEMA
-- ============================================================================
-- Created: March 24, 2026
-- Purpose: Add group trip, expense, and voting tables for multi-user collaboration
-- 
-- Supports both Oracle and PostgreSQL (commented syntax variations)
-- ============================================================================

-- ============================================================================
-- TABLE 1: GROUP_TRIPS
-- Description: Base table for group trips
-- ============================================================================
CREATE TABLE group_trips (
  id VARCHAR2(36) PRIMARY KEY,
  title VARCHAR2(255) NOT NULL,
  destination_id VARCHAR2(36) NOT NULL,
  owner_id VARCHAR2(36) NOT NULL,
  budget NUMBER(12, 2) NOT NULL CHECK (budget > 0),
  duration NUMBER(3) NOT NULL CHECK (duration > 0),
  start_date DATE,
  status VARCHAR2(20) DEFAULT 'draft' CHECK (status IN ('draft', 'planning', 'published', 'completed')),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (owner_id) REFERENCES users(id),
  FOREIGN KEY (destination_id) REFERENCES destinations(id)
);

-- INDEX for faster queries
CREATE INDEX idx_group_trips_owner_id ON group_trips(owner_id);
CREATE INDEX idx_group_trips_destination_id ON group_trips(destination_id);
CREATE INDEX idx_group_trips_status ON group_trips(status);

-- ============================================================================
-- TABLE 2: GROUP_MEMBERS
-- Description: Track who is in each group trip and their role
-- ============================================================================
CREATE TABLE group_members (
  id VARCHAR2(36) PRIMARY KEY,
  group_trip_id VARCHAR2(36) NOT NULL,
  user_id VARCHAR2(36) NOT NULL,
  role VARCHAR2(20) DEFAULT 'member' CHECK (role IN ('owner', 'editor', 'member', 'viewer')),
  joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  status VARCHAR2(20) DEFAULT 'pending' CHECK (status IN ('pending', 'active', 'declined', 'left')),
  FOREIGN KEY (group_trip_id) REFERENCES group_trips(id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  UNIQUE (group_trip_id, user_id)
);

CREATE INDEX idx_group_members_group_trip_id ON group_members(group_trip_id);
CREATE INDEX idx_group_members_user_id ON group_members(user_id);
CREATE INDEX idx_group_members_status ON group_members(status);

-- ============================================================================
-- TABLE 3: EXPENSES
-- Description: Track expenses in group trips
-- ============================================================================
CREATE TABLE expenses (
  id VARCHAR2(36) PRIMARY KEY,
  group_trip_id VARCHAR2(36) NOT NULL,
  description VARCHAR2(255) NOT NULL,
  amount NUMBER(12, 2) NOT NULL CHECK (amount > 0),
  paid_by VARCHAR2(36) NOT NULL,
  category VARCHAR2(50) DEFAULT 'other' CHECK (category IN ('accommodation', 'food', 'transport', 'activity', 'other')),
  paid_date DATE DEFAULT TRUNC(SYSDATE),
  status VARCHAR2(20) DEFAULT 'pending' CHECK (status IN ('pending', 'settled')),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (group_trip_id) REFERENCES group_trips(id) ON DELETE CASCADE,
  FOREIGN KEY (paid_by) REFERENCES users(id)
);

CREATE INDEX idx_expenses_group_trip_id ON expenses(group_trip_id);
CREATE INDEX idx_expenses_paid_by ON expenses(paid_by);
CREATE INDEX idx_expenses_status ON expenses(status);

-- ============================================================================
-- TABLE 4: EXPENSE_SPLITS
-- Description: How each expense is split among members
-- ============================================================================
CREATE TABLE expense_splits (
  id VARCHAR2(36) PRIMARY KEY,
  expense_id VARCHAR2(36) NOT NULL,
  user_id VARCHAR2(36) NOT NULL,
  amount_owed NUMBER(12, 2) NOT NULL CHECK (amount_owed >= 0),
  FOREIGN KEY (expense_id) REFERENCES expenses(id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  UNIQUE (expense_id, user_id)
);

CREATE INDEX idx_expense_splits_expense_id ON expense_splits(expense_id);
CREATE INDEX idx_expense_splits_user_id ON expense_splits(user_id);

-- ============================================================================
-- TABLE 5: POLLS
-- Description: Voting on group trip decisions
-- ============================================================================
CREATE TABLE polls (
  id VARCHAR2(36) PRIMARY KEY,
  group_trip_id VARCHAR2(36) NOT NULL,
  created_by VARCHAR2(36) NOT NULL,
  question VARCHAR2(500) NOT NULL,
  poll_type VARCHAR2(20) DEFAULT 'itinerary' CHECK (poll_type IN ('itinerary', 'budget', 'date', 'activity', 'destination')),
  status VARCHAR2(20) DEFAULT 'active' CHECK (status IN ('active', 'locked', 'resolved')),
  expires_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (group_trip_id) REFERENCES group_trips(id) ON DELETE CASCADE,
  FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE INDEX idx_polls_group_trip_id ON polls(group_trip_id);
CREATE INDEX idx_polls_status ON polls(status);

-- ============================================================================
-- TABLE 6: POLL_OPTIONS
-- Description: Answer choices for each poll
-- ============================================================================
CREATE TABLE poll_options (
  id VARCHAR2(36) PRIMARY KEY,
  poll_id VARCHAR2(36) NOT NULL,
  option_text VARCHAR2(500) NOT NULL,
  vote_count NUMBER(10) DEFAULT 0,
  sequence NUMBER(3),
  FOREIGN KEY (poll_id) REFERENCES polls(id) ON DELETE CASCADE
);

CREATE INDEX idx_poll_options_poll_id ON poll_options(poll_id);

-- ============================================================================
-- TABLE 7: POLL_VOTES
-- Description: Individual votes on poll options
-- ============================================================================
CREATE TABLE poll_votes (
  id VARCHAR2(36) PRIMARY KEY,
  poll_option_id VARCHAR2(36) NOT NULL,
  user_id VARCHAR2(36) NOT NULL,
  voted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (poll_option_id) REFERENCES poll_options(id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  UNIQUE (poll_option_id, user_id)
);

CREATE INDEX idx_poll_votes_poll_option_id ON poll_votes(poll_option_id);
CREATE INDEX idx_poll_votes_user_id ON poll_votes(user_id);

-- ============================================================================
-- TABLE 8: SETTLEMENTS
-- Description: Track who owes whom after all expenses
-- ============================================================================
CREATE TABLE settlements (
  id VARCHAR2(36) PRIMARY KEY,
  group_trip_id VARCHAR2(36) NOT NULL,
  debtor_id VARCHAR2(36) NOT NULL,
  creditor_id VARCHAR2(36) NOT NULL,
  amount NUMBER(12, 2) NOT NULL,
  status VARCHAR2(20) DEFAULT 'pending' CHECK (status IN ('pending', 'settled')),
  settled_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (group_trip_id) REFERENCES group_trips(id) ON DELETE CASCADE,
  FOREIGN KEY (debtor_id) REFERENCES users(id),
  FOREIGN KEY (creditor_id) REFERENCES users(id)
);

CREATE INDEX idx_settlements_group_trip_id ON settlements(group_trip_id);
CREATE INDEX idx_settlements_debtor_id ON settlements(debtor_id);
CREATE INDEX idx_settlements_creditor_id ON settlements(creditor_id);
CREATE INDEX idx_settlements_status ON settlements(status);

-- ============================================================================
-- TRIGGERS & VIEWS FOR COMMON QUERIES
-- ============================================================================

-- View: Group Trip Summary with member count and total expenses
CREATE VIEW vw_group_trips_summary AS
SELECT 
  gt.id,
  gt.title,
  gt.destination_id,
  gt.owner_id,
  gt.budget,
  gt.duration,
  gt.status,
  COUNT(DISTINCT gm.user_id) as member_count,
  COALESCE(SUM(e.amount), 0) as total_expenses,
  gt.created_at
FROM group_trips gt
LEFT JOIN group_members gm ON gt.id = gm.group_trip_id AND gm.status = 'active'
LEFT JOIN expenses e ON gt.id = e.group_trip_id
GROUP BY gt.id, gt.title, gt.destination_id, gt.owner_id, gt.budget, gt.duration, gt.status, gt.created_at;

-- View: Settlement summary for a group trip
CREATE VIEW vw_settlements_summary AS
SELECT 
  gt.id as group_trip_id,
  gt.title,
  s.debtor_id,
  s.creditor_id,
  s.amount,
  s.status,
  u_debtor.name as debtor_name,
  u_creditor.name as creditor_name
FROM settlements s
JOIN group_trips gt ON s.group_trip_id = gt.id
JOIN users u_debtor ON s.debtor_id = u_debtor.id
JOIN users u_creditor ON s.creditor_id = u_creditor.id;

-- ============================================================================
-- SAMPLE DATA FOR TESTING
-- ============================================================================
-- Insert sample group trip (for testing)
-- INSERT INTO group_trips VALUES (
--   'grp-001', 
--   'Goa Beach Trip', 
--   'dest-002', 
--   'user-001', 
--   100000, 
--   5, 
--   TO_DATE('2026-06-01', 'YYYY-MM-DD'),
--   'planning',
--   CURRENT_TIMESTAMP,
--   CURRENT_TIMESTAMP
-- );

-- Insert sample members
-- INSERT INTO group_members VALUES (
--   'grpmem-001', 'grp-001', 'user-001', 'owner', CURRENT_TIMESTAMP, 'active'
-- );
-- INSERT INTO group_members VALUES (
--   'grpmem-002', 'grp-001', 'user-002', 'editor', CURRENT_TIMESTAMP, 'active'
-- );

-- ============================================================================
-- COMMENTS & NOTES
-- ============================================================================
-- Status tracking:
--   group_trips: draft -> planning -> published -> completed
--   group_members: pending -> active (or declined/left)
--   expenses: pending -> settled
--   polls: active -> locked -> resolved
--
-- Role hierarchy:
--   owner: Can add/remove members, modify trip, delete trip
--   editor: Can add items, expenses, create polls
--   member: Can see everything, vote, contribute expenses
--   viewer: Read-only access
--
-- Expense splitting:
--   Equal split: amount / member_count
--   Custom split: specified by expense_splits table
--
-- Settlements:
--   Calculated by algorithm: Who paid what, who owes whom
--   Can be marked settled (to track payment history)

-- ============================================================================
-- PostgreSQL EQUIVALENT SYNTAX (Uncomment to use PostgreSQL)
-- ============================================================================
/*

-- PostgreSQL Version:
CREATE TABLE group_trips (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title VARCHAR(255) NOT NULL,
  destination_id UUID NOT NULL,
  owner_id UUID NOT NULL,
  budget DECIMAL(12, 2) NOT NULL CHECK (budget > 0),
  duration INTEGER NOT NULL CHECK (duration > 0),
  start_date DATE,
  status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'planning', 'published', 'completed')),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (owner_id) REFERENCES users(id),
  FOREIGN KEY (destination_id) REFERENCES destinations(id)
);

-- Similar for all other tables, just replace:
--   VARCHAR2(n) -> VARCHAR(n)
--   NUMBER(12,2) -> DECIMAL(12,2)
--   NUMBER(3) -> INTEGER
--   NUMBER(10) -> INTEGER
--   TRUNC(SYSDATE) -> CURRENT_DATE
--   UUID instead of VARCHAR(36)
--   ON DELETE CASCADE works same as Oracle

*/

-- ============================================================================
-- END OF PHASE A SCHEMA
-- ============================================================================
