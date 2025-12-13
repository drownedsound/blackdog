.mode column
.headers on

-- Enforce Foreign Keys (Default is OFF in SQLite)
PRAGMA foreign_keys = ON;
-- Enable Write-Ahead Logging for concurrency and reliability
PRAGMA journal_mode = WAL;

SELECT '----------------------------------------' AS '======================';
SELECT 'PHASE 1: CONFIGURATION CHECK' AS 'Step';
SELECT '----------------------------------------' AS '----------------------';

-- 1. Check Pragma settings
SELECT 
    'Foreign Key Enforcement' AS setting, 
    CASE 
        WHEN (SELECT * FROM pragma_foreign_keys()) = 1 THEN 'PASS' 
        ELSE 'FAIL' 
    END AS status
UNION ALL
SELECT 
    'Journal Mode (WAL)', 
    CASE 
        WHEN (SELECT * FROM pragma_journal_mode()) = 'wal' THEN 'PASS' 
        ELSE 'FAIL' 
    END;

SELECT '----------------------------------------' AS '======================';
SELECT 'PHASE 2: SCHEMA INVENTORY' AS 'Step';
SELECT '----------------------------------------' AS '----------------------';

-- 2. Check Core Tables
SELECT name AS table_name, 'Exists' AS status 
FROM sqlite_master 
WHERE type='table' AND name NOT LIKE 'sqlite_%' 
ORDER BY name;

-- 3. Check Critical Triggers (Immutability & Audit)
SELECT '----------------------------------------' AS '----------------------';
SELECT 'Verifying Triggers...' AS info;

SELECT name AS trigger_name, 'Active' AS status 
FROM sqlite_master 
WHERE type='trigger' 
ORDER BY name;

SELECT '----------------------------------------' AS '======================';
SELECT 'PHASE 3: REFERENCE DATA HYDRATION' AS 'Step';
SELECT '----------------------------------------' AS '----------------------';

-- 4. Count Reference Data to ensure 02_populate_ref_data.sql ran correctly
SELECT 'REF_PRODUCT_CATEGORY' AS table_name, count(*) AS row_count 
FROM REF_PRODUCT_CATEGORY
UNION ALL
SELECT 'REF_PRODUCT_TYPE', count(*) FROM REF_PRODUCT_TYPE
UNION ALL
SELECT 'REF_APP_STATUS', count(*) FROM REF_APP_STATUS
UNION ALL
SELECT 'REF_APPLICANT_ROLE', count(*) FROM REF_APPLICANT_ROLE
UNION ALL
SELECT 'REF_ADDRESS_TYPE', count(*) FROM REF_ADDRESS_TYPE
UNION ALL
SELECT 'REF_CONTACT_TYPE', count(*) FROM REF_CONTACT_TYPE
UNION ALL
SELECT 'REF_ID_TYPE', count(*) FROM REF_ID_TYPE
UNION ALL
SELECT 'REF_EMPLOYMENT_STATUS', count(*) FROM REF_EMPLOYMENT_STATUS
UNION ALL
SELECT 'REF_EDUCATION_LEVEL', count(*) FROM REF_EDUCATION_LEVEL;

SELECT '----------------------------------------' AS '======================';
SELECT 'PHASE 4: FUNCTIONAL SMOKE TEST' AS 'Step';
SELECT '----------------------------------------' AS '----------------------';

-- 5. Functional Test: Simulate a valid application to test FKs and Audit
--    We wrap this in a TRANSACTION and ROLLBACK to avoid polluting the DB.

BEGIN TRANSACTION;

SELECT 'Attempting Insert: Application (Standard Visa)...' AS action;

-- A. Insert Parent Application
INSERT INTO APPLICATION (
    member_reference_no, category_code, status_code, requested_amount
    )
VALUES ('TEST-APP-001', 'CARD', 1, 5000000); -- 50k PHP

-- B. Insert Applicant (Primary)
INSERT INTO APPLICANT (
    application_id, role_code, product_type_code, 
    first_name, last_name, date_of_birth
    )
VALUES (
    last_insert_rowid(), 'PRIMARY_CARDHOLDER', 'VISA_GOLD', 
    'Maria', 'Clara', '1990-01-01'
    );

-- C. Insert Identity
INSERT INTO APPLICANT_IDENTIFICATION (
    applicant_id, id_type_code, id_value, issuing_authority, expiry_date
    )
VALUES (
    (SELECT id FROM APPLICANT WHERE last_name='Clara'), 
    'PASSPORT', 'P9999999A', 'DFA', '2030-01-01'
    );

-- D. Insert Polymorphic Detail (Card)
INSERT INTO CARD_DETAIL (application_id, card_design_code, credit_limit)
VALUES (
    (SELECT id FROM APPLICATION WHERE member_reference_no='TEST-APP-001'), 
    'STD_BLUE', 5000000
    );

SELECT '----------------------------------------' AS '----------------------';
SELECT 'Validation 1: Data Integrity Check' AS info;

-- Verify the Join works
SELECT 
    m.member_reference_no, 
    p.first_name || ' ' || p.last_name AS applicant_name,
    c.description AS card_type,
    cd.credit_limit
FROM APPLICATION a
JOIN APPLICANT p ON a.id = p.application_id
JOIN REF_PRODUCT_TYPE c ON p.product_type_code = c.code
JOIN CARD_DETAIL cd ON a.id = cd.application_id
WHERE m.member_reference_no = 'TEST-APP-001';

SELECT '----------------------------------------' AS '----------------------';
SELECT 'Validation 2: Audit Log Trigger Check' AS info;

-- Verify Audit Log captured the INSERT
SELECT 
    table_name, 
    action, 
    new_value 
FROM AUDIT_LOG 
WHERE record_id = (
    SELECT id FROM APPLICATION WHERE member_reference_no='TEST-APP-001'
);

SELECT '----------------------------------------' AS '----------------------';
SELECT 'Test Complete. Rolling back changes...' AS info;

ROLLBACK;

SELECT '----------------------------------------' AS '======================';
SELECT 'SCHEMA VERIFICATION COMPLETE' AS 'Status';
SELECT '----------------------------------------' AS '----------------------';
