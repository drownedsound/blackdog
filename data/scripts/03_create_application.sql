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

INSERT INTO CREDIT_CARD (
    profile_id, currency_id, credit_limit
) VALUES (
    1, 1, 5000000
);

INSERT INTO APPLICATION (
    member_reference_no,
    status_id,
    credit_card_id,
    requested_amount,
    created_at,
    updated_at
) VALUES (
    'APP-CC-001',
    1,                   -- CREATED
    last_insert_rowid(), -- Links to CREDIT_CARD
    6000000,
    unixepoch(),
    unixepoch()
);

INSERT INTO APPLICANT (
    application_id,
    is_principal,
    last_name,
    first_name,
    middle_name,
    birthday
) VALUES (
    last_insert_rowid(),
    1,
    'SMITH',
    'JOHN',
    'DOE',
    '1980-01-01'
);

INSERT INTO CONTACT_NUMBER (
    applicant_id, type_id, value
) VALUES (
    last_insert_rowid(),
    1,
    '1234567'
);

-- ============================================================================
-- 3. PERSONAL LOAN APPLICATION
-- ============================================================================

-- INSERT INTO PERSONAL_LOAN (profile_id, loan_amount)
-- VALUES (1, 5000000);
--
-- INSERT INTO APPLICATION (
--      member_reference_no,
--      status_id,
--      personal_loan_id,
--      requested_amount,
--      created_at,
--      updated_at
-- ) VALUES (
--      'APP-PL-001',
--      1,
--      last_insert_rowid(),
--      6000000,
--      unixepoch(),
--      unixepoch()
-- );

COMMIT;
