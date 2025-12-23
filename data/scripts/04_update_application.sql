.mode column
.headers on

-- ============================================================================
-- 1. CONFIGURATION
-- ============================================================================
PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;

BEGIN TRANSACTION;

-- ============================================================================
-- 2. UPDATE APP-CC-001 (Pasada King)
-- ============================================================================
-- Goal: Add contact to principal.
-- Limit: Orig 5M, Req 6M. Target: 5.5M.

-- 2.1 Add Contact Number to Principal (Home Phone)
INSERT INTO CONTACT_NUMBER (applicant_id, type_id, value)
SELECT a.id, 2, '0281112233'
FROM APPLICANT a
JOIN APPLICATION app ON a.application_id = app.id
WHERE app.member_reference_no = 'APP-CC-001' AND a.is_principal = 1;

-- 2.2 Update Credit Limit
UPDATE CREDIT_CARD
SET credit_limit = 5500000
WHERE id = (
    SELECT credit_card_id 
    FROM APPLICATION 
    WHERE member_reference_no = 'APP-CC-001'
);

-- 2.3 Touch Updated At
UPDATE APPLICATION
SET updated_at = unixepoch()
WHERE member_reference_no = 'APP-CC-001';


-- ============================================================================
-- 3. UPDATE APP-CC-002 (Amigas Platinum)
-- ============================================================================
-- Goal: Add new applicant. Add contact to principal.
-- Limit: Orig 30M, Req 30M. Target: 28M.

-- 3.1 Add New Supplementary Applicant
INSERT INTO APPLICANT (
    application_id, is_principal, last_name, first_name, birthday
)
SELECT id, 0, 'SANTOS', 'ANA', '1995-02-14'
FROM APPLICATION 
WHERE member_reference_no = 'APP-CC-002';

-- 3.2 Add Contact for the Principal (Office Phone)
INSERT INTO CONTACT_NUMBER (applicant_id, type_id, value)
SELECT a.id, 3, '0289990000'
FROM APPLICANT a
JOIN APPLICATION app ON a.application_id = app.id
WHERE app.member_reference_no = 'APP-CC-002' AND a.is_principal = 1;

-- 3.3 Update Credit Limit
UPDATE CREDIT_CARD
SET credit_limit = 28000000
WHERE id = (
    SELECT credit_card_id 
    FROM APPLICATION 
    WHERE member_reference_no = 'APP-CC-002'
);

-- 3.4 Touch Updated At
UPDATE APPLICATION
SET updated_at = unixepoch()
WHERE member_reference_no = 'APP-CC-002';

-- ============================================================================
-- 4. UPDATE APP-PL-001 (Personal Loan)
-- ============================================================================
-- Goal: Add contact number to applicant.

-- 4.1 Add Contact Number (Home Phone)
INSERT INTO CONTACT_NUMBER (applicant_id, type_id, value)
SELECT a.id, 2, '0287654321'
FROM APPLICANT a
JOIN APPLICATION app ON a.application_id = app.id
WHERE app.member_reference_no = 'APP-PL-001';

-- 4.2 Touch Updated At
UPDATE APPLICATION
SET updated_at = unixepoch()
WHERE member_reference_no = 'APP-PL-001';

COMMIT;

