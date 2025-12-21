.mode column
.headers on

-- ============================================================================
-- 1. CONFIGURATION
-- ============================================================================
PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;

-- ============================================================================
-- 2. REFERENCES 
-- ============================================================================

CREATE TABLE REF_APPLICATION_STATUS (
    -- StatusCreated (1)
    -- StatusInProgress (2)
    -- StatusApproved (3)
    -- StatusDeclined (4)
    -- StatusCancelled (5)
    id INTEGER PRIMARY KEY, 
    -- CREATED
    -- IN_PROGRESS
    -- APPROVED
    -- DECLINED
    -- CANCELLED
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    is_terminal INTEGER NOT NULL DEFAULT 0 

    CHECK (is_terminal IN (0, 1))
) STRICT;

CREATE TABLE REF_CREDIT_CARD (
    id  INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    credit_limit INTEGER NOT NULL CHECK (credit_limit > 0), 
    interest_rate INTEGER NOT NULL CHECK (interest_rate > 0) 
) STRICT;

-- ============================================================================
-- 3. APPLICATION DETAILS
-- ============================================================================

CREATE TABLE APPLICATION (
    id INTEGER PRIMARY KEY, 
    member_reference_no TEXT NOT NULL UNIQUE,
    status_id INTEGER NOT NULL,
    credit_card_id INTEGER,
    loan_id INTEGER,
    requested_amount INTEGER NOT NULL,
    created_at INTEGER NOT NULL DEFAULT (unixepoch()),
    updated_at INTEGER NOT NULL DEFAULT (unixepoch()),

    FOREIGN KEY (status_id) 
    REFERENCES REF_APPLICATION_STATUS(id)
    FOREIGN KEY (credit_card_id) 
    REFERENCES CREDIT_CARD(id)

    CHECK (requested_amount > 0)
    -- Ensures that an application will have at least one product
    CHECK (credit_card_id IS NOT NULL OR loan_id IS NOT NULL)
    -- Ensure valid epoch time (> Jan 1 2020) 
    CHECK (created_at > 1577836800)
    CHECK (updated_at > 1577836800)
) STRICT;

-- ============================================================================
-- 4. PRODUCT DETAILS
-- ============================================================================

CREATE TABLE CREDIT_CARD (
    id INTEGER PRIMARY KEY,
    profile_id INTEGER NOT NULL,
    credit_limit INTEGER NOT NULL CHECK (credit_limit > 0), 

    FOREIGN KEY (profile_id)
    REFERENCES REF_CREDIT_CARD(id)
) STRICT;

