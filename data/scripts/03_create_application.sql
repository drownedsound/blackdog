.mode column
.headers on

-- ============================================================================
-- 1. CONFIGURATION
-- ============================================================================
PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;

BEGIN TRANSACTION;

-- ============================================================================
-- 2. CREDIT CARD
-- ============================================================================
INSERT INTO CREDIT_CARD (profile_id, credit_limit)
VALUES (1, 7500000);

-- ============================================================================
-- 3. APPLICATION
-- ============================================================================
INSERT INTO APPLICATION (
    member_reference_no,
    status_id,
    credit_card_id,
    requested_amount,
    created_at,
    updated_at
) VALUES (
    'APP-2023-QC-FULL-01',
    1,                       -- Status: CREATED
    last_insert_rowid(),     -- Links to the card created above
    7500000,             
    unixepoch(),
    unixepoch()
);

COMMIT;

