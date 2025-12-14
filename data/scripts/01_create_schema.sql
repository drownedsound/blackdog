.mode column
.headers on

-- 1. CONFIGURATION
-- ============================================================================
PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;

-- 2. REFERENCE DATA (CONTROLLED VOCABULARY)
-- ============================================================================

CREATE TABLE REF_PRODUCT_CATEGORY (
    -- CategoryCard (1)
    -- CategoryLoan (2)
    id INTEGER PRIMARY KEY, 
    -- Credit Card
    -- Personal Loan
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL
) STRICT;

-- CREATE TABLE REF_PRODUCT_TYPE (
--     id INTEGER PRIMARY KEY, 
--     name TEXT NOT NULL UNIQUE,
--     description TEXT NOT NULL,
--     category_id INTEGER NOT NULL,
--
--     FOREIGN KEY (category_id) 
--     REFERENCES REF_PRODUCT_CATEGORY(id)
-- ) STRICT;

CREATE TABLE REF_APP_STATUS (
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

-- CREATE TABLE REF_APPLICANT_ROLE (
--     code TEXT PRIMARY KEY,
--     description TEXT NOT NULL,
--     category_code TEXT NOT NULL,
--     CONSTRAINT fk_role_category 
--     FOREIGN KEY (category_code) 
--     REFERENCES REF_PRODUCT_CATEGORY(code)
-- );
--
-- CREATE TABLE REF_ADDRESS_TYPE (
--     code TEXT PRIMARY KEY,
--     description TEXT NOT NULL
-- );
--
-- CREATE TABLE REF_CONTACT_TYPE (
--     code TEXT PRIMARY KEY,
--     description TEXT NOT NULL
-- );
--
-- CREATE TABLE REF_ID_TYPE (
--     code TEXT PRIMARY KEY,
--     description TEXT NOT NULL
-- );
--
-- CREATE TABLE REF_EMPLOYMENT_STATUS (
--     code TEXT PRIMARY KEY,
--     description TEXT NOT NULL
-- );
--
-- CREATE TABLE REF_EDUCATION_LEVEL (
--     code TEXT PRIMARY KEY,
--     description TEXT NOT NULL
-- );

-- 3. TRANSACTIONAL CORE
-- ============================================================================

CREATE TABLE APPLICATION (
    id INTEGER PRIMARY KEY, 
    member_reference_no TEXT NOT NULL UNIQUE,
    category_id INTEGER NOT NULL,
    status_id INTEGER NOT NULL,
    requested_amount INTEGER NOT NULL,
    created_at INTEGER NOT NULL DEFAULT (unixepoch()),
    updated_at INTEGER NOT NULL DEFAULT (unixepoch()),

    FOREIGN KEY (category_id) 
    REFERENCES REF_PRODUCT_CATEGORY(id),
    FOREIGN KEY (status_id) 
    REFERENCES REF_APP_STATUS(id)

    CHECK (requested_amount > 0)
-- Ensure valid epoch time (> Jan 1 2020) 
CHECK (created_at > 1577836800)
CHECK (updated_at > 1577836800)
) STRICT;

-- CREATE TABLE APPLICANT (
--     id INTEGER PRIMARY KEY,
--     application_id INTEGER NOT NULL,
--     role_code TEXT NOT NULL,
--     product_type_code TEXT NOT NULL,
--     first_name TEXT NOT NULL CHECK (length(trim(first_name)) > 0),
--     last_name TEXT NOT NULL CHECK (length(trim(last_name)) > 0),
--     date_of_birth TEXT NOT NULL,
--
--     CONSTRAINT fk_applicant_app 
--     FOREIGN KEY (application_id) 
--     REFERENCES APPLICATION(id),
--     CONSTRAINT fk_applicant_role 
--     FOREIGN KEY (role_code) 
--     REFERENCES REF_APPLICANT_ROLE(code),
--     CONSTRAINT fk_applicant_prod_type 
--     FOREIGN KEY (product_type_code) 
--     REFERENCES REF_PRODUCT_TYPE(code),
--     -- Date Format Check (YYYY-MM-DD)
--     CONSTRAINT chk_dob_fmt 
--     CHECK (date_of_birth GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]')
-- );
--
-- -- 4. APPLICANT IDENTITY (NORMALIZED)
-- -- ============================================================================
-- CREATE TABLE APPLICANT_IDENTIFICATION (
--     id INTEGER PRIMARY KEY,
--     applicant_id INTEGER NOT NULL,
--     id_type_code TEXT NOT NULL,
--     id_value TEXT NOT NULL CHECK (length(trim(id_value)) > 0),
--     issuing_authority TEXT NOT NULL,
--     expiry_date TEXT, -- Nullable for non-expiring IDs
--
--     CONSTRAINT fk_ident_applicant 
--     FOREIGN KEY (applicant_id) 
--     REFERENCES APPLICANT(id),
--     CONSTRAINT fk_ident_type 
--     FOREIGN KEY (id_type_code) 
--     REFERENCES REF_ID_TYPE(code),
--     -- Prevent duplicate IDs in system
--     CONSTRAINT uq_ident_type_val 
--     UNIQUE (id_type_code, id_value) 
-- );
--
-- -- 5. PRODUCT SPECIFICS (POLYMORPHIC EXTENSIONS)
-- -- ============================================================================
--
-- CREATE TABLE LOAN_DETAIL (
--     application_id INTEGER PRIMARY KEY,
--     term_months INTEGER NOT NULL 
--     CHECK (term_months > 0 AND term_months <= 120),
--     interest_rate INTEGER NOT NULL 
--     CHECK (interest_rate >= 0), -- Basis points
--     loan_amount INTEGER NOT NULL CHECK (loan_amount > 0), -- Centavos
--
--     CONSTRAINT fk_loan_app 
--     FOREIGN KEY (application_id) 
--     REFERENCES APPLICATION(id)
-- );
--
-- CREATE TABLE CARD_DETAIL (
--     application_id INTEGER PRIMARY KEY,
--     card_design_code TEXT NOT NULL,
--     credit_limit INTEGER NOT NULL CHECK (credit_limit > 0), -- Centavos
--
--     CONSTRAINT fk_card_app 
--     FOREIGN KEY (application_id) 
--     REFERENCES APPLICATION(id)
-- );
--
-- -- 6. APPLICANT DETAILS (DEMOGRAPHICS)
-- -- ============================================================================
--
-- CREATE TABLE APPLICANT_ADDRESS (
--     id INTEGER PRIMARY KEY,
--     applicant_id INTEGER NOT NULL,
--     address_type_code TEXT NOT NULL,
--     line_1 TEXT NOT NULL CHECK (length(trim(line_1)) > 0),
--     line_2 TEXT,
--     city_municipality TEXT NOT NULL CHECK (length(trim(city_municipality)) > 0),
--     province TEXT NOT NULL CHECK (length(trim(province)) > 0),
--     postal_code TEXT NOT NULL CHECK (length(trim(postal_code)) > 0),
--
--     CONSTRAINT fk_addr_applicant 
--     FOREIGN KEY (applicant_id) 
--     REFERENCES APPLICANT(id),
--     CONSTRAINT fk_addr_type 
--     FOREIGN KEY (address_type_code) 
--     REFERENCES REF_ADDRESS_TYPE(code)
-- );
--
-- CREATE TABLE APPLICANT_CONTACT (
--     id INTEGER PRIMARY KEY,
--     applicant_id INTEGER NOT NULL,
--     contact_type_code TEXT NOT NULL,
--     contact_number TEXT NOT NULL CHECK (length(trim(contact_number)) > 0),
--
--     CONSTRAINT fk_contact_applicant 
--     FOREIGN KEY (applicant_id) 
--     REFERENCES APPLICANT(id),
--     CONSTRAINT fk_contact_type 
--     FOREIGN KEY (contact_type_code) 
--     REFERENCES REF_CONTACT_TYPE(code)
-- );
--
-- CREATE TABLE APPLICANT_EMPLOYMENT (
--     id INTEGER PRIMARY KEY,
--     applicant_id INTEGER NOT NULL,
--     employment_status_code TEXT NOT NULL,
--     employer_name TEXT NOT NULL,
--     gross_monthly_income INTEGER NOT NULL 
--     CHECK (gross_monthly_income >= 0), -- Centavos
--
--     CONSTRAINT fk_emp_applicant 
--     FOREIGN KEY (applicant_id) 
--     REFERENCES APPLICANT(id),
--     CONSTRAINT fk_emp_status 
--     FOREIGN KEY (employment_status_code) 
--     REFERENCES REF_EMPLOYMENT_STATUS(code)
-- );
--
-- CREATE TABLE APPLICANT_EDUCATION (
--     id INTEGER PRIMARY KEY,
--     applicant_id INTEGER NOT NULL,
--     education_level_code TEXT NOT NULL,
--     institution_name TEXT NOT NULL,
--
--     CONSTRAINT fk_edu_applicant 
--     FOREIGN KEY (applicant_id) 
--     REFERENCES APPLICANT(id),
--     CONSTRAINT fk_edu_lvl 
--     FOREIGN KEY (education_level_code) 
--     REFERENCES REF_EDUCATION_LEVEL(code)
-- );
--
-- 7. AUDIT TRAIL
-- ============================================================================

CREATE TABLE AUDIT_LOG (
    id INTEGER PRIMARY KEY,
    table_name TEXT NOT NULL,
    record_id INTEGER NOT NULL,
    action TEXT NOT NULL CHECK (action IN ('INSERT', 'UPDATE')),
    old_value TEXT, -- JSON Payload
    new_value TEXT NOT NULL, -- JSON Payload
    changed_at INTEGER NOT NULL DEFAULT (unixepoch())
) STRICT;

-- ============================================================================
-- 8. BUSINESS LOGIC TRIGGERS
-- ============================================================================

-- 8.1 IMMUTABILITY 
-- ----------------------------------------------------------------------------
-- TODO: Verify deletion is not allowed
-- Implementing the requirement that records cannot be deleted.
CREATE TRIGGER no_delete_application BEFORE DELETE ON APPLICATION 
BEGIN SELECT RAISE(ABORT, 'Access Denied: Records are immutable.'); END;

-- CREATE TRIGGER no_delete_applicant BEFORE DELETE ON APPLICANT 
-- BEGIN SELECT RAISE(ABORT, 'Access Denied: Records are immutable.'); END;
--
-- CREATE TRIGGER no_delete_ident BEFORE DELETE ON APPLICANT_IDENTIFICATION 
-- BEGIN SELECT RAISE(ABORT, 'Access Denied: Records are immutable.'); END;
--
-- CREATE TRIGGER no_delete_loan BEFORE DELETE ON LOAN_DETAIL 
-- BEGIN SELECT RAISE(ABORT, 'Access Denied: Records are immutable.'); END;
--
-- CREATE TRIGGER no_delete_card BEFORE DELETE ON CARD_DETAIL 
-- BEGIN SELECT RAISE(ABORT, 'Access Denied: Records are immutable.'); END;
--
-- CREATE TRIGGER no_delete_addr BEFORE DELETE ON APPLICANT_ADDRESS 
-- BEGIN SELECT RAISE(ABORT, 'Access Denied: Records are immutable.'); END;
--
-- CREATE TRIGGER no_delete_contact BEFORE DELETE ON APPLICANT_CONTACT 
-- BEGIN SELECT RAISE(ABORT, 'Access Denied: Records are immutable.'); END;
--
-- CREATE TRIGGER no_delete_emp BEFORE DELETE ON APPLICANT_EMPLOYMENT 
-- BEGIN SELECT RAISE(ABORT, 'Access Denied: Records are immutable.'); END;
--
-- -- 8.2 POLYMORPHIC EXCLUSIVITY 
-- -- ----------------------------------------------------------------------------
-- -- Ensure an Application cannot be both a Loan and a Card.
--
-- CREATE TRIGGER check_exclusive_loan_insert
-- BEFORE INSERT ON LOAN_DETAIL
-- BEGIN
-- SELECT RAISE(ABORT, 'Integrity Error: Application already has CARD details.')
-- WHERE EXISTS (
-- SELECT 1 FROM CARD_DETAIL WHERE application_id = NEW.application_id
-- );
-- END;
--
-- CREATE TRIGGER check_exclusive_card_insert
-- BEFORE INSERT ON CARD_DETAIL
-- BEGIN
-- SELECT RAISE(ABORT, 'Integrity Error: Application already has LOAN details.')
-- WHERE EXISTS (
-- SELECT 1 FROM LOAN_DETAIL WHERE application_id = NEW.application_id
-- );
-- END;
--
-- -- 8.3 CROSS-TABLE CONSISTENCY (ROLE vs CATEGORY)
-- -- ----------------------------------------------------------------------------
-- -- Ensure Applicant Role (e.g. Guarantor) matches Product Category (e.g. Loan).
--
-- CREATE TRIGGER validate_applicant_role_category
-- BEFORE INSERT ON APPLICANT
-- BEGIN
-- SELECT RAISE(ABORT, 
-- 'Business Error: Applicant Role invalid for Product Category.')
-- WHERE (
-- SELECT r.category_code 
-- FROM REF_APPLICANT_ROLE r 
-- WHERE r.code = NEW.role_code
-- ) != (
-- SELECT a.category_code 
-- FROM APPLICATION a 
-- WHERE a.id = NEW.application_id
-- );
-- END;

-- 8.4 AUDIT LOGGING (JSON via SQLite)
-- ----------------------------------------------------------------------------
-- Capture state changes for key tables.

-- APPLICATION Audit
CREATE TRIGGER audit_app_insert AFTER INSERT ON APPLICATION
BEGIN
INSERT INTO AUDIT_LOG (table_name, record_id, action, new_value)
    VALUES ('APPLICATION', NEW.id, 'INSERT', 
json_object(
'member_reference_no', NEW.member_reference_no,
'category', NEW.category_id, 
'status', NEW.status_id, 
'amount', NEW.requested_amount
));
END;

-- CREATE TRIGGER audit_app_update AFTER UPDATE ON APPLICATION
-- BEGIN
--     INSERT INTO AUDIT_LOG (table_name, record_id, action, old_value, new_value)
--     VALUES ('APPLICATION', NEW.id, 'UPDATE', 
--         json_object(
--         'member_reference_no', OLD.member_reference_no,
--         'category', OLD.category_id, 
--         'status', OLD.status_id, 
--         'amount', OLD.requested_amount
--         ),
--         json_object(
--         'member_reference_no', NEW.member_reference_no,
--         'category', NEW.category_id, 
--         'status', NEW.status_id, 
--         'amount', NEW.requested_amount
--         ));
-- END;

CREATE TRIGGER audit_app_update AFTER UPDATE ON APPLICATION
-- Add this guard clause:
WHEN OLD.status_id IS NOT NEW.status_id 
OR OLD.requested_amount IS NOT NEW.requested_amount
BEGIN
INSERT INTO AUDIT_LOG (table_name, record_id, action, old_value, new_value)
VALUES ('APPLICATION', NEW.id, 'UPDATE', 
json_object(
'status', OLD.status_id, 
'amount', OLD.requested_amount
),
json_object(
'status', NEW.status_id, 
'amount', NEW.requested_amount
)
);
END;

-- -- APPLICANT Audit
-- CREATE TRIGGER audit_applicant_update AFTER UPDATE ON APPLICANT
-- BEGIN
--     INSERT INTO AUDIT_LOG (table_name, record_id, action, old_value, new_value)
--     VALUES ('APPLICANT', NEW.id, 'UPDATE', 
--         json_object(
--         'first_name', OLD.first_name, 
--         'last_name', OLD.last_name
--         ),
--         json_object(
--         'first_name', NEW.first_name, 
--         'last_name', NEW.last_name
--         ));
-- END;
--
-- Timestamp Maintenance
CREATE TRIGGER update_timestamp_app AFTER UPDATE ON APPLICATION
BEGIN
UPDATE APPLICATION SET updated_at = unixepoch() WHERE id = NEW.id;
END;
