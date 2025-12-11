-- TODO: Rename application_number to member_reference_no

BEGIN TRANSACTION;

-- 1. Create the Parent Application
-- ----------------------------------------------------------------------------
INSERT INTO APPLICATION (
    application_number, category_code, status_code, requested_amount
    ) VALUES (
    'APP-2023-SINGLE', 'CARD', 'CREATED', 5000000 -- 50,000.00 PHP
    );

-- 2. Create Product Specifics (Card Details)
-- ----------------------------------------------------------------------------
INSERT INTO CARD_DETAIL (
    application_id, card_design_code, credit_limit
    ) VALUES (
    (SELECT id FROM APPLICATION WHERE application_number = 'APP-2023-SINGLE'),
    'STD_BLUE',
    5000000
    );

-- 3. Create the Applicant (Primary)
-- ----------------------------------------------------------------------------
INSERT INTO APPLICANT (
    application_id, role_code, product_type_code, 
    first_name, last_name, date_of_birth
    ) VALUES (
    (SELECT id FROM APPLICATION WHERE application_number = 'APP-2023-SINGLE'),
    'PRIMARY_CARDHOLDER',
    'VISA_GOLD',
    'Juan',
    'Dela Cruz',
    '1985-05-15'
    );

-- 4. Create Applicant Details
--    (Using a CTE-like subquery strategy for the Applicant ID)
-- ----------------------------------------------------------------------------

-- 4.1 Identification
INSERT INTO APPLICANT_IDENTIFICATION (
    applicant_id, id_type_code, id_value, issuing_authority, expiry_date
    ) VALUES (
    (SELECT id FROM APPLICANT WHERE last_name = 'Dela Cruz' 
        AND first_name = 'Juan'),
    'PASSPORT',
    'P1234567A',
    'DFA Manila',
    '2028-05-15'
    );

-- 4.2 Address
INSERT INTO APPLICANT_ADDRESS (
    applicant_id, address_type_code, line_1, line_2, 
    city_municipality, province, postal_code
    ) VALUES (
    (SELECT id FROM APPLICANT WHERE last_name = 'Dela Cruz' 
        AND first_name = 'Juan'),
    'RESIDENTIAL',
    'Unit 404, Building A',
    'Rizal Avenue, Brgy 678',
    'Manila',
    'Metro Manila',
    '1000'
    );

-- 4.3 Contact
INSERT INTO APPLICANT_CONTACT (
    applicant_id, contact_type_code, contact_number
    ) VALUES (
    (SELECT id FROM APPLICANT WHERE last_name = 'Dela Cruz' 
        AND first_name = 'Juan'),
    'MOBILE',
    '+639171234567'
    );

-- 4.4 Employment
INSERT INTO APPLICANT_EMPLOYMENT (
    applicant_id, employment_status_code, employer_name, gross_monthly_income
    ) VALUES (
    (SELECT id FROM APPLICANT WHERE last_name = 'Dela Cruz' 
        AND first_name = 'Juan'),
    'REGULAR',
    'Acme Corp Philippines',
    6000000 -- 60,000.00 PHP
    );

-- 4.5 Education
INSERT INTO APPLICANT_EDUCATION (
    applicant_id, education_level_code, institution_name
    ) VALUES (
    (SELECT id FROM APPLICANT WHERE last_name = 'Dela Cruz' 
        AND first_name = 'Juan'),
    'COLLEGE',
    'University of the Philippines Diliman'
    );

COMMIT;

SELECT 'Single Applicant Created: APP-2023-SINGLE' as status;
