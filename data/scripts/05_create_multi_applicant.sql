BEGIN TRANSACTION;

-- 1. Create the Parent Application
INSERT INTO APPLICATION (
    member_reference_no, category_code, status_code, requested_amount
    ) VALUES (
    'APP-2023-MULTI', 'CARD', 1, 10000000 -- 100k Shared Limit
    );

INSERT INTO CARD_DETAIL (
    application_id, card_design_code, credit_limit
    ) VALUES (
    (SELECT id FROM APPLICATION WHERE member_reference_no = 'APP-2023-MULTI'),
    'PMT_BLACK',
    10000000
    );

-- ============================================================================
-- 2. APPLICANT 1: Maria Clara (PRIMARY)
-- ============================================================================
INSERT INTO APPLICANT (
    application_id, role_code, product_type_code, 
    first_name, last_name, date_of_birth
    ) VALUES (
    (SELECT id FROM APPLICATION WHERE member_reference_no = 'APP-2023-MULTI'),
    'PRIMARY_CARDHOLDER', 'VISA_PLATINUM', 'Maria', 'Clara', '1980-01-01'
    );

-- Maria's Details
INSERT INTO APPLICANT_IDENTIFICATION (
    applicant_id, id_type_code, id_value, issuing_authority, expiry_date
    ) VALUES (
    (SELECT id FROM APPLICANT WHERE last_name='Clara' AND first_name='Maria'),
    'SSS', '01-2345678-9', 'SSS Main', 'N/A'
    );

INSERT INTO APPLICANT_ADDRESS (
    applicant_id, address_type_code, line_1, line_2, city_municipality, 
    province, postal_code
    ) VALUES (
    (SELECT id FROM APPLICANT WHERE last_name='Clara' AND first_name='Maria'),
    'RESIDENTIAL', '123 Dasma Village', 'Makati Ave', 'Makati', 'NCR', '1200'
    );

INSERT INTO APPLICANT_CONTACT (applicant_id, contact_type_code, contact_number)
VALUES ((SELECT id FROM APPLICANT WHERE first_name='Maria'), 
    'MOBILE', '+639180000001');

INSERT INTO APPLICANT_EMPLOYMENT (
    applicant_id, employment_status_code, employer_name, gross_monthly_income
    ) VALUES (
    (SELECT id FROM APPLICANT WHERE first_name='Maria'), 
    'SELF_EMPLOYED', 'Clara Trading Inc', 15000000
    );

INSERT INTO APPLICANT_EDUCATION (
    applicant_id, education_level_code, institution_name
    ) VALUES (
    (SELECT id FROM APPLICANT WHERE first_name='Maria'), 
    'POST_GRAD', 'Asian Institute of Management'
    );

-- ============================================================================
-- 3. APPLICANT 2: Crisostomo Ibarra (SUPPLEMENTARY - Husband)
-- ============================================================================
INSERT INTO APPLICANT (
    application_id, role_code, product_type_code, 
    first_name, last_name, date_of_birth
    ) VALUES (
    (SELECT id FROM APPLICATION WHERE member_reference_no = 'APP-2023-MULTI'),
    'SUPPLEMENTARY_CARDHOLDER', 'VISA_PLATINUM', 'Crisostomo', 'Ibarra', 
    '1978-02-14'
    );

-- Crisostomo's Details (Simplified for brevity, but all tables populated)
INSERT INTO APPLICANT_IDENTIFICATION (
    applicant_id, id_type_code, id_value, issuing_authority, expiry_date
    ) VALUES (
    (SELECT id FROM APPLICANT WHERE last_name='Ibarra'),
    'PASSPORT', 'P9876543B', 'DFA', '2029-01-01'
    );

INSERT INTO APPLICANT_ADDRESS (
    applicant_id, address_type_code, line_1, line_2, city_municipality, 
    province, postal_code
    ) VALUES (
    (SELECT id FROM APPLICANT WHERE last_name='Ibarra'),
    'RESIDENTIAL', '123 Dasma Village', 'Makati Ave', 'Makati', 'NCR', '1200'
    );

INSERT INTO APPLICANT_CONTACT (applicant_id, contact_type_code, contact_number)
VALUES ((SELECT id FROM APPLICANT WHERE last_name='Ibarra'), 
    'MOBILE', '+639180000002');

INSERT INTO APPLICANT_EMPLOYMENT (
    applicant_id, employment_status_code, employer_name, gross_monthly_income
    ) VALUES (
    (SELECT id FROM APPLICANT WHERE last_name='Ibarra'), 
    'OFW', 'Madrid Architects', 20000000
    );

INSERT INTO APPLICANT_EDUCATION (
    applicant_id, education_level_code, institution_name
    ) VALUES (
    (SELECT id FROM APPLICANT WHERE last_name='Ibarra'), 
    'COLLEGE', 'Ateneo de Manila'
    );

-- ============================================================================
-- 4. APPLICANT 3: Salome Clara (SUPPLEMENTARY - Daughter)
-- ============================================================================
INSERT INTO APPLICANT (
    application_id, role_code, product_type_code, 
    first_name, last_name, date_of_birth
    ) VALUES (
    (SELECT id FROM APPLICATION WHERE member_reference_no = 'APP-2023-MULTI'),
    'SUPPLEMENTARY_CARDHOLDER', 'VISA_PLATINUM', 'Salome', 'Clara', 
    '2002-12-25'
    );

INSERT INTO APPLICANT_IDENTIFICATION (
    applicant_id, id_type_code, id_value, issuing_authority, expiry_date
    ) VALUES (
    (SELECT id FROM APPLICANT WHERE first_name='Salome'),
    'PHILSYS', '1111-2222-3333', 'PSA', 'N/A'
    );

INSERT INTO APPLICANT_ADDRESS (
    applicant_id, address_type_code, line_1, line_2, city_municipality, 
    province, postal_code
    ) VALUES (
    (SELECT id FROM APPLICANT WHERE first_name='Salome'),
    'RESIDENTIAL', 'Condo Unit 101', 'Katipunan Ave', 'Quezon City', 
    'NCR', '1108'
    );

INSERT INTO APPLICANT_CONTACT (applicant_id, contact_type_code, contact_number)
VALUES ((SELECT id FROM APPLICANT WHERE first_name='Salome'), 
    'MOBILE', '+639180000003');

INSERT INTO APPLICANT_EMPLOYMENT (
    applicant_id, employment_status_code, employer_name, gross_monthly_income
    ) VALUES (
    (SELECT id FROM APPLICANT WHERE first_name='Salome'), 
    'FREELANCE', 'Graphic Design', 2000000
    );

INSERT INTO APPLICANT_EDUCATION (
    applicant_id, education_level_code, institution_name
    ) VALUES (
    (SELECT id FROM APPLICANT WHERE first_name='Salome'), 
    'COLLEGE', 'UP Fine Arts'
    );

COMMIT;

SELECT 'Multi-Applicant Created: APP-2023-MULTI (3 Applicants)' as status;
