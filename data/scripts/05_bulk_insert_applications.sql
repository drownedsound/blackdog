.mode column
.headers on

-- ============================================================================
-- 1. CONFIGURATION
-- ============================================================================
PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;

BEGIN TRANSACTION;

-- ============================================================================
-- 1. CTE GENERATOR
-- ============================================================================
WITH RECURSIVE 
  params(max_rows, start_cc_id, start_pl_id) AS (
    VALUES(
        1000, 
        (SELECT COALESCE(MAX(id), 0) FROM CREDIT_CARD), 
        (SELECT COALESCE(MAX(id), 0) FROM PERSONAL_LOAN)
    )
  ),
  
  -- Expanded Name Repositories (50 First * 50 Last = 2500 Combinations)
  first_names(name) AS (
    VALUES 
    ('James'),('Mary'),('Robert'),('Patricia'),('John'),('Jennifer'),
    ('Michael'),('Linda'),('David'),('Elizabeth'),('William'),('Barbara'),
    ('Richard'),('Susan'),('Joseph'),('Jessica'),('Thomas'),('Sarah'),
    ('Charles'),('Karen'),('Christopher'),('Lisa'),('Daniel'),('Nancy'),
    ('Matthew'),('Betty'),('Anthony'),('Margaret'),('Mark'),('Sandra'),
    ('Donald'),('Ashley'),('Steven'),('Kimberly'),('Paul'),('Emily'),
    ('Andrew'),('Donna'),('Joshua'),('Michelle'),('Kenneth'),('Carol'),
    ('Kevin'),('Amanda'),('Brian'),('Melissa'),('George'),('Deborah'),
    ('Edward'),('Stephanie')
  ),
  last_names(name) AS (
    VALUES 
    ('Smith'),('Johnson'),('Williams'),('Brown'),('Jones'),('Garcia'),
    ('Miller'),('Davis'),('Rodriguez'),('Martinez'),('Hernandez'),('Lopez'),
    ('Gonzalez'),('Wilson'),('Anderson'),('Thomas'),('Taylor'),('Moore'),
    ('Jackson'),('Martin'),('Lee'),('Perez'),('Thompson'),('White'),
    ('Harris'),('Sanchez'),('Clark'),('Ramirez'),('Lewis'),('Robinson'),
    ('Walker'),('Young'),('Allen'),('King'),('Wright'),('Scott'),
    ('Torres'),('Nguyen'),('Hill'),('Flores'),('Green'),('Adams'),
    ('Nelson'),('Baker'),('Hall'),('Rivera'),('Campbell'),('Mitchell'),
    ('Carter'),('Roberts')
  ),

  -- 1. Generate Sequence
  seq(n) AS (
    SELECT 1 
    UNION ALL 
    SELECT n + 1 FROM seq, params WHERE n < params.max_rows
  ),

  -- 2. Define Attributes & Calculate Linkage IDs
  dataset AS (
    SELECT 
      n,
      CASE WHEN n % 10 = 0 THEN 'PL' ELSE 'CC' END as app_type,
      -- Pre-calculate Foreign Key IDs using Window Functions
      params.start_cc_id + ROW_NUMBER() OVER (
          PARTITION BY CASE WHEN n % 10 = 0 THEN 'PL' ELSE 'CC' END 
          ORDER BY n
      ) as product_fk_id,
      CASE 
        WHEN n % 10 = 0 THEN 1 
        WHEN (abs(random()) % 100) < 30 THEN (abs(random()) % 4) + 2 
        ELSE 1 
      END as num_applicants,
      -- Generate pseudo-Nanoid (20 chars hex) for uniqueness and speed
      lower(hex(randomblob(10))) as ref_no
    FROM seq, params
  )

-- ============================================================================
-- BULK INSERT
-- ============================================================================

-- 1. Insert Products
INSERT INTO CREDIT_CARD (id, profile_id, credit_limit)
SELECT 
    product_fk_id,
    (abs(random()) % 4) + 1,
    ((abs(random()) % 50) + 10) * 100000
FROM dataset WHERE app_type = 'CC';

INSERT INTO PERSONAL_LOAN (id, profile_id, loan_amount)
SELECT 
    product_fk_id,
    (abs(random()) % 3) + 1,
    ((abs(random()) % 100) + 10) * 100000
FROM dataset WHERE app_type = 'PL';

-- 2. Insert Applications
INSERT INTO APPLICATION (
    member_reference_no, status_id, credit_card_id, personal_loan_id, 
    requested_amount, created_at, updated_at
)
SELECT 
    ref_no,
    (abs(random()) % 5) + 1,
    CASE WHEN app_type = 'CC' THEN product_fk_id ELSE NULL END,
    CASE WHEN app_type = 'PL' THEN product_fk_id ELSE NULL END,
    5000000,
    unixepoch(),
    unixepoch()
FROM dataset;

-- 3. Insert Applicants (Fan-out CTE with Random Names & Birthdays)
INSERT INTO APPLICANT (
    application_id, is_principal, last_name, first_name, birthday
)
WITH RECURSIVE 
  app_fanout(n, applicant_idx, total_applicants, app_id_lookup) AS (
    SELECT 
      d.n, 1, d.num_applicants,
      (SELECT id FROM APPLICATION WHERE member_reference_no = d.ref_no)
    FROM dataset d
    UNION ALL
    SELECT n, applicant_idx + 1, total_applicants, app_id_lookup
    FROM app_fanout 
    WHERE applicant_idx < total_applicants
  )
SELECT 
    app_id_lookup,
    CASE WHEN applicant_idx = 1 THEN 1 ELSE 0 END,
    -- Pick Random Last Name
    (SELECT name FROM last_names ORDER BY random() LIMIT 1),
    -- Pick Random First Name
    (SELECT name FROM first_names ORDER BY random() LIMIT 1),
    -- Random Birthday: 1970-01-01 (0) to 2000-12-31 (~978307200)
    date(abs(random()) % 978307200, 'unixepoch')
FROM app_fanout;

-- 4. Insert Contact Numbers
INSERT INTO CONTACT_NUMBER (applicant_id, type_id, value)
WITH RECURSIVE
  contact_gen(applicant_id, contact_idx, total_contacts) AS (
      -- Filter to only newly inserted applicants by scanning recent IDs
      -- This assumes strict sequential ID generation for the session
      SELECT id, 1, (abs(random()) % 3) + 1 
      FROM APPLICANT 
      WHERE id > (SELECT MAX(id) - 2000 FROM APPLICANT) 
      UNION ALL
      SELECT applicant_id, contact_idx + 1, total_contacts
      FROM contact_gen WHERE contact_idx < total_contacts
  )
SELECT 
    applicant_id,
    (abs(random()) % 3) + 1,
    '09' || (abs(random()) % 1000000000)
FROM contact_gen;

COMMIT;
