-- ============================================================================
-- 1. CONFIGURATION
-- ============================================================================

PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;

BEGIN TRANSACTION;

-- ============================================================================
-- 2. SINGLE APPLICANT CREDIT CARD APPLICATION
-- ============================================================================

INSERT INTO CREDIT_CARD (
    id, 
    profile_id, 
    credit_limit
) VALUES (
    7410754270103867392, 
    1, 
    5000000
);

INSERT INTO APPLICATION (
    id,
    member_reference_no,
    status_id,
    credit_card_id,
    requested_amount,
    created_at,
    updated_at
) VALUES (
    7410754167536357376,
    'APP-CC-001',
    1,                   
    7410754270103867392,
    6000000,
    unixepoch(),
    unixepoch()
);

INSERT INTO APPLICANT (
    id,
    application_id,
    is_principal,
    last_name,
    first_name,
    middle_name,
    birthday
) VALUES (
    7410754725357817856,
    7410754167536357376,
    1,
    'SMITH',
    'JOHN',
    'DOE',
    '1980-01-01'
);

INSERT INTO CONTACT_NUMBER (
    id,
    applicant_id, 
    type_id, 
    value
) VALUES (
    7410758405897326592,
    7410754725357817856,
    1,
    '1234567'
);

-- ============================================================================
-- 3. SINGLE APPLICANT PERSONAL LOAN APPLICATION
-- ============================================================================

INSERT INTO PERSONAL_LOAN (
    id,
    profile_id, 
    loan_amount
) VALUES (
    7410755240611287040,
    1,
    200000000
);

INSERT INTO APPLICATION (
    id,
    member_reference_no,
    status_id,
    personal_loan_id,
    requested_amount,
    created_at,
    updated_at
) VALUES (
    7410755338086912000,
    'APP-PL-001',
    1,                   
    7410755240611287040,
    200000000,
    unixepoch(),
    unixepoch()
);

INSERT INTO APPLICANT (
    id,
    application_id,
    is_principal,
    last_name,
    first_name,
    middle_name,
    birthday
) VALUES (
    7410755498057666560,
    7410755338086912000,
    1,                  
    'SMITH',
    'DOE',
    'JANE',
    '1985-05-20'
);

INSERT INTO CONTACT_NUMBER (
    id,
    applicant_id, 
    type_id, 
    value
) VALUES (
    7410755610020417536,
    7410755498057666560,
    1,                  
    '09171112222'
);


-- ============================================================================
-- 4. MULTI APPLICANT CREDIT CARD APPLICATION
-- ============================================================================

INSERT INTO CREDIT_CARD (
    id,
    profile_id, 
    credit_limit
) VALUES (
    7410755959405940736,
    2, 
    30000000
);

INSERT INTO APPLICATION (
    id,
    member_reference_no,
    status_id,
    credit_card_id,
    requested_amount,
    created_at,
    updated_at
) VALUES (
    7410756033057918976,
    'APP-CC-002',
    1,
    7410755959405940736,
    30000000,
    unixepoch(),
    unixepoch()
);

INSERT INTO APPLICANT (
    id,
    application_id, 
    is_principal, 
    last_name, 
    first_name, 
    birthday
) VALUES (
    7410756138540470272,
    7410756033057918976,
    1,
    'SANTOS',
    'MARIA',
    '1990-08-15'
);

-- Principal Contact 1: Mobile
INSERT INTO CONTACT_NUMBER (
    id,
    applicant_id, 
    type_id, 
    value
) VALUES (
    7410756396527915008,
    7410756138540470272,
    1,
    '09181234567'
);

-- Principal Contact 2: Landline
-- Note: We must lookup the ID again because last_insert_rowid changed above
INSERT INTO CONTACT_NUMBER (
    id,
    applicant_id, 
    type_id, 
    value
) VALUES (
    7410756463674527744,
    7410756138540470272,
    2,
    '0288881234'
);

INSERT INTO APPLICANT (
    id,
    application_id, 
    is_principal, 
    last_name, 
    first_name, 
    birthday
) VALUES (
    7410756670961225728,
    7410756033057918976,
    0, -- Supplementary
    'SANTOS',
    'JUAN',
    '1988-12-01'
);

-- Supplementary Contact: Mobile
INSERT INTO CONTACT_NUMBER (
    id,
    applicant_id, 
    type_id, 
    value
) VALUES (
    7410756919549235200,
    7410756670961225728,
    1,
    '09199876543'
);

COMMIT;
