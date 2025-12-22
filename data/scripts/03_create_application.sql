.mode column
.headers on

-- ============================================================================
-- 1. CONFIGURATION
-- ============================================================================
PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;

BEGIN TRANSACTION;

-- ============================================================================
-- 2. CREDIT CARD APPLICATION
-- ============================================================================
INSERT INTO CREDIT_CARD (profile_id, credit_limit)
VALUES (1, 5000000);

INSERT INTO APPLICATION (
    member_reference_no,
    status_id,
    credit_card_id,
    requested_amount,
    created_at,
    updated_at
) VALUES (
    'APP-CC-001',
    1,                       
    last_insert_rowid(),    
    6000000,             
    unixepoch(),
    unixepoch()
);

-- ============================================================================
-- 3. PERSONAL LOAN APPLICATION
-- ============================================================================

INSERT INTO PERSONAL_LOAN (profile_id, loan_amount)
VALUES (1, 5000000);

INSERT INTO APPLICATION (
    member_reference_no,
    status_id,
    personal_loan_id,
    requested_amount,
    created_at,
    updated_at
) VALUES (
    'APP-PL-001',
    1,                       
    last_insert_rowid(),    
    6000000,             
    unixepoch(),
    unixepoch()
);
COMMIT;

 
