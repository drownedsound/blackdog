.mode column
.headers on

-- ============================================================================
-- 1. CONFIGURATION
-- ============================================================================

PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;

BEGIN TRANSACTION;

-- ============================================================================
-- 2. SINGLE APPLICANT CREDIT CARD APPLICATION
-- ============================================================================

INSERT INTO CREDIT_CARD (profile_id, credit_limit) VALUES (1, 5000000);

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
-- 3. SINGLE APPLICANT PERSONAL LOAN APPLICATION
-- ============================================================================

INSERT INTO PERSONAL_LOAN (profile_id, loan_amount) VALUES (1,200000000);

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
    200000000,
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
    'DOE',
    'JANE',
    '1985-05-20'
);

INSERT INTO CONTACT_NUMBER (
    applicant_id, type_id, value
) VALUES (
    last_insert_rowid(), 
    1,                  
    '09171112222'
);


-- ============================================================================
-- 4. MULTI APPLICANT CREDIT CARD APPLICATION
-- ============================================================================

INSERT INTO CREDIT_CARD (profile_id, credit_limit) VALUES (2, 30000000);

INSERT INTO APPLICATION (
    member_reference_no,
    status_id,
    credit_card_id,
    requested_amount,
    created_at,
    updated_at
) VALUES (
    'APP-CC-002',
    1,
    last_insert_rowid(),
    30000000,
    unixepoch(),
    unixepoch()
);

INSERT INTO APPLICANT (
    application_id, is_principal, last_name, first_name, birthday
) VALUES (
    last_insert_rowid(), -- Links to APPLICATION
    1,
    'SANTOS',
    'MARIA',
    '1990-08-15'
);

-- Principal Contact 1: Mobile
INSERT INTO CONTACT_NUMBER (
    applicant_id, type_id, value
) VALUES (
    last_insert_rowid(), -- Links to PRINCIPAL
    1,
    '09181234567'
);

-- Principal Contact 2: Landline
-- Note: We must lookup the ID again because last_insert_rowid changed above
INSERT INTO CONTACT_NUMBER (
    applicant_id, type_id, value
) VALUES (
    (SELECT id FROM APPLICANT WHERE last_name = 'SANTOS' AND is_principal = 1),
    2,
    '0288881234'
);

INSERT INTO APPLICANT (
    application_id, is_principal, last_name, first_name, birthday
) VALUES (
    -- We cannot use last_insert_rowid() here; it points to a contact number.
    -- We lookup the Application ID by its unique reference number.
    (SELECT id FROM APPLICATION WHERE member_reference_no = 'APP-CC-002'),
    0, -- Supplementary
    'SANTOS',
    'JUAN',
    '1988-12-01'
);

-- Supplementary Contact: Mobile
INSERT INTO CONTACT_NUMBER (
    applicant_id, type_id, value
) VALUES (
    last_insert_rowid(), -- Links to SUPPLEMENTARY
    1,
    '09199876543'
);

COMMIT;
