-- ============================================================================
-- SCRIPT: 04_insert_app_full_report.sql
-- DESCRIPTION: Atomic Application Creation with Comprehensive Output
-- ============================================================================

-- Force Foreign Key enforcement for this session
PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;

BEGIN TRANSACTION;

-- 1. Create the Product Instance (Pasada King, 75k limit)
INSERT INTO CREDIT_CARD (profile_id, credit_limit)
VALUES (1, 7500000);

-- 2. Create the Application linked to the new Product
-- Uses last_insert_rowid() to capture the ID from step 1
INSERT INTO APPLICATION (
    member_reference_no,
    status_id,
    credit_card_id,
    requested_amount,
    created_at,
    updated_at
) VALUES (
    'APP-2023-QC-FULL-01',   -- Unique Reference
    1,                       -- Status: CREATED
    last_insert_rowid(),     -- Links to the card created above
    7500000,                 -- Requested Amount
    unixepoch(),
    unixepoch()
);

COMMIT;

-- -- ============================================================================
-- -- FULL VERIFICATION REPORT
-- -- ============================================================================
-- SELECT 
--     -- APPLICATION TABLE FIELDS
--     a.id AS App_ID,
--     a.member_reference_no,
--     a.status_id,
--     s.name AS Status_Name,
--     
--     -- LINKED PRODUCT FIELDS
--     a.credit_card_id AS Card_ID,
--     c.profile_id AS Profile_ID,
--     p.name AS Product_Name,
--     c.credit_limit AS Limit_Minor,
--     
--     -- FINANCIALS
--     a.requested_amount,
--     
--     -- AUDIT TRAIL
--     datetime(a.created_at, 'unixepoch') AS Created_UTC,
--     datetime(a.updated_at, 'unixepoch') AS Updated_UTC,
--     
--     -- METADATA
--     p.description AS Product_Desc
-- FROM APPLICATION a
-- JOIN REF_APPLICATION_STATUS s ON a.status_id = s.id
-- LEFT JOIN CREDIT_CARD c ON a.credit_card_id = c.id
-- LEFT JOIN REF_CREDIT_CARD p ON c.profile_id = p.id
-- WHERE a.member_reference_no = 'APP-2023-QC-FULL-01';
