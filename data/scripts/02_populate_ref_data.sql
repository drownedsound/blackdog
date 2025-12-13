.mode column
.headers on

BEGIN TRANSACTION;

-- 1. PRODUCT CATEGORIES
-- ============================================================================
INSERT INTO REF_PRODUCT_CATEGORY (id, name, description) VALUES 
    (1, 'Credit Card', 'A revolving credit facility allowing the holder to make purchases up to a specific limit. Balances can be paid in full monthly or over time with interest.'),
    (2, 'Personal Loan', 'A fixed-term, lump-sum loan repaid in regular monthly installments (principal + interest) over a set period. Typically unsecured.');

-- 2. PRODUCT TYPES
-- ============================================================================
-- Credit Cards
INSERT INTO REF_PRODUCT_TYPE (id, name, description, category_id) VALUES 
    (1, 'VISA_CLASSIC', 'Visa Classic', 1),
    (2, 'VISA_GOLD', 'Visa Gold', 1),
    (3, 'VISA_PLATINUM', 'Visa Platinum', 1),
    (4, 'MC_CLASSIC', 'Mastercard Classic', 1),
    (5, 'MC_GOLD', 'Mastercard Gold', 1),
    (6, 'MC_PLATINUM', 'Mastercard Platinum', 1);

    -- -- Personal Loans
    -- INSERT INTO REF_PRODUCT_TYPE (code, description, category_code) VALUES 
    --     ('PL_STANDARD', 'Standard Personal Loan', 'LOAN'),
    --     ('PL_SALARY', 'Salary Deduction Loan', 'LOAN'), 
    --     ('PL_OFW', 'OFW Reintegration Loan', 'LOAN'), 
    --     ('PL_MICRO', 'Micro-Finance Personal Loan', 'LOAN');

    -- 3. APPLICATION STATUS (Strict Workflow)
    -- ============================================================================
INSERT INTO REF_APP_STATUS (id, name, description, is_terminal) VALUES 
    (1, 'CREATED', 'Draft created but not submitted', 0),
    (2, 'IN_PROGRESS', 'Undergoing manual review', 0),
    (3, 'APPROVED', 'Approved but pending booking', 1), 
    (4, 'DECLINED', 'Rejected based on credit policy', 1),
    (5, 'CANCELLED', 'Withdrawn or terminated', 1);

-- -- 4. APPLICANT ROLES 
-- -- ============================================================================
-- -- Card Roles
-- INSERT INTO REF_APPLICANT_ROLE (code, description, category_code) VALUES 
--     ('PRIMARY_CARDHOLDER', 'Primary account holder', 'CARD'),
--     ('SUPPLEMENTARY_CARDHOLDER', 'Authorized user', 'CARD');
--
-- -- Loan Roles
-- INSERT INTO REF_APPLICANT_ROLE (code, description, category_code) VALUES 
--     ('PRINCIPAL_BORROWER', 'Primary obligor', 'LOAN'),
--     ('CO_BORROWER', 'Jointly liable obligor', 'LOAN'),
--     ('CO_SIGNER', 'Secondary obligor in case of default', 'LOAN'),
--     ('GUARANTOR', 'Third-party securitizer', 'LOAN');
--
-- -- 5. ADDRESS TYPES
-- -- ============================================================================
-- INSERT INTO REF_ADDRESS_TYPE (code, description) VALUES 
--     ('RESIDENTIAL', 'Current place of residence'),
--     ('PERMANENT', 'Permanent provincial address'), 
--     ('OFFICE', 'Place of employment');
--
-- -- 6. CONTACT TYPES
-- -- ============================================================================
-- INSERT INTO REF_CONTACT_TYPE (code, description) VALUES 
--     ('MOBILE', 'Mobile Phone'),
--     ('LANDLINE', 'Fixed Line Telephone'),
--     ('EMAIL', 'Personal Email Address'),
--     ('WORK_EMAIL', 'Corporate Email Address');
--
-- -- 7. IDENTIFICATION TYPES (Philippine KYC Compliance)
-- -- ============================================================================
-- INSERT INTO REF_ID_TYPE (code, description) VALUES 
--     ('PASSPORT', 'Philippine Passport (DFA)'),
--     ('PHILSYS', 'Philippine Identification System ID (National ID)'),
--     ('DRIVERS_LICENSE', 'LTO Driver License'),
--     ('SSS', 'Social Security System ID'),
--     ('GSIS', 'Government Service Insurance System ID'),
--     ('UMID', 'Unified Multi-Purpose ID'),
--     ('TIN', 'Tax Identification Number Card (BIR)'),
--     ('PRC', 'Professional Regulation Commission ID'),
--     ('VOTERS', 'COMELEC Voter ID'),
--     ('POSTAL', 'PhilPost Postal ID');
--
-- -- 8. EMPLOYMENT STATUS
-- -- ============================================================================
-- INSERT INTO REF_EMPLOYMENT_STATUS (code, description) VALUES 
--     ('REGULAR', 'Permanent/Regular Employee'),
--     ('PROBATIONARY', 'Probationary Employee (< 6 months)'),
--     ('CONTRACTUAL', 'Fixed-term Contract Employee'),
--     ('SELF_EMPLOYED', 'Business Owner / Professional'),
--     ('FREELANCE', 'Freelancer / Gig Economy Worker'),
--     ('OFW', 'Overseas Filipino Worker'), 
--     ('RETIRED', 'Retired with Pension'),
--     ('UNEMPLOYED', 'Currently Unemployed');
--
-- -- 9. EDUCATION LEVEL
-- -- ============================================================================
-- INSERT INTO REF_EDUCATION_LEVEL (code, description) VALUES 
--     ('HIGH_SCHOOL', 'High School Diploma'),
--     ('VOCATIONAL', 'Vocational / Trade Certificate'),
--     ('COLLEGE', 'Bachelor''s Degree'),
--     ('POST_GRAD', 'Master''s or Doctorate Degree'),
--     ('NONE', 'No Formal Education');

COMMIT;
