-- ============================================================================
-- 1. GET ALL APPLICATIONS (WITH PAGING, FILTERING AND SORTING)
-- ============================================================================
-- Returns flat rows: App -> Product -> Applicant -> Contact

SELECT
    -- 1. APPLICATION CORE
    app.member_reference_no,
    app.created_at,
    app.updated_at,

    -- 2. STATUS DETAILS
    ras.name                    AS status_name,

    -- 3. PRODUCT: CREDIT CARD (NULL if PL)
    rcc.name                    AS cc_name,
    cc_cur.code                 AS cc_currency_code,
    cc.credit_limit             AS cc_limit,
    rcc.interest_rate           AS cc_rate,

    -- 4. PRODUCT: PERSONAL LOAN (NULL if CC)
    rpl.name                    AS pl_name,
    pl_cur.code                 AS pl_currency_code,
    pl.loan_amount              AS pl_amount,
    rpl.interest_rate           AS pl_rate,

    -- 5. APPLICANT DETAILS
    a.last_name,
    a.first_name,
    a.middle_name,
    a.birthday,
    a.is_principal,

    -- 6. CONTACT DETAILS
    rct.name                    AS contact_type_name,
    c.value                     AS contact_value

FROM APPLICATION app

-- Join Status
JOIN REF_APPLICATION_STATUS ras ON app.status_id = ras.id

-- Join Credit Card chain
LEFT JOIN CREDIT_CARD cc        ON app.credit_card_id = cc.id
LEFT JOIN REF_CREDIT_CARD rcc   ON cc.profile_id = rcc.id
LEFT JOIN REF_CURRENCY cc_cur   ON rcc.currency_id = cc_cur.id

-- Join Personal Loan chain
LEFT JOIN PERSONAL_LOAN pl      ON app.personal_loan_id = pl.id
LEFT JOIN REF_PERSONAL_LOAN rpl ON pl.profile_id = rpl.id
LEFT JOIN REF_CURRENCY pl_cur   ON rpl.currency_id = pl_cur.id

-- Join Applicants
JOIN APPLICANT a                ON app.id = a.application_id

-- Join Contacts
LEFT JOIN CONTACT_NUMBER c      ON a.id = c.applicant_id
LEFT JOIN REF_CONTACT_TYPE rct  ON c.type_id = rct.id

-- ============================================================================
-- DYNAMIC PREDICATES
-- ============================================================================

-- HOW TO FILTER:
-- * CC Only: WHERE app.credit_card_id IS NOT NULL
-- * PL Only: WHERE app.personal_loan_id IS NOT NULL

-- Sorting: Critical for Row Folding. Must group by App -> Applicant.
ORDER BY 
    app.created_at DESC, 
    app.id, 
    a.is_principal DESC, -- Principals first
    a.id

-- HOW TO PAGE:
-- * Page 1 (Size 50): LIMIT 50 OFFSET 0
-- * Page 2 (Size 50): LIMIT 50 OFFSET 50
-- * Page 3 (Size 50): LIMIT 50 OFFSET 100
LIMIT 50 OFFSET 0;

-- ============================================================================
-- 2. GET FULL APPLICATION DETAILS
-- ============================================================================

SELECT
    -- 1. APPLICATION CORE
    app.id                      AS app_id,
    app.member_reference_no,
    app.created_at,
    app.updated_at,
    app.requested_amount,

    -- 2. STATUS DETAILS
    app.status_id,
    ras.name                    AS status_name,
    ras.is_terminal             AS status_is_terminal,

    -- 3. PRODUCT: CREDIT CARD (NULL if PL)
    cc.id                       AS cc_id,
    cc.credit_limit             AS cc_limit,
    rcc.name                    AS cc_name,
    rcc.product_ceiling         AS cc_ceiling,
    rcc.interest_rate           AS cc_rate,
    rcc.currency_id             AS cc_currency_id,
    cc_cur.code                 AS cc_currency_code,

    -- 4. PRODUCT: PERSONAL LOAN (NULL if CC)
    pl.id                       AS pl_id,
    pl.loan_amount              AS pl_amount,
    rpl.name                    AS pl_name,
    rpl.product_ceiling         AS pl_ceiling,
    rpl.interest_rate           AS pl_rate,
    rpl.currency_id             AS pl_currency_id,
    pl_cur.code                 AS pl_currency_code,

    -- 5. APPLICANT DETAILS
    a.id                        AS applicant_id,
    a.is_principal,
    a.first_name,
    a.middle_name,
    a.last_name,
    a.birthday,

    -- 6. CONTACT DETAILS
    c.id                        AS contact_id,
    c.value                     AS contact_value,
    c.type_id                   AS contact_type_id,
    rct.name                    AS contact_type_name

FROM APPLICATION app

-- Join Status
JOIN REF_APPLICATION_STATUS ras ON app.status_id = ras.id

-- Join Credit Card chain
LEFT JOIN CREDIT_CARD cc        ON app.credit_card_id = cc.id
LEFT JOIN REF_CREDIT_CARD rcc   ON cc.profile_id = rcc.id
LEFT JOIN REF_CURRENCY cc_cur   ON rcc.currency_id = cc_cur.id

-- Join Personal Loan chain
LEFT JOIN PERSONAL_LOAN pl      ON app.personal_loan_id = pl.id
LEFT JOIN REF_PERSONAL_LOAN rpl ON pl.profile_id = rpl.id
LEFT JOIN REF_CURRENCY pl_cur   ON rpl.currency_id = pl_cur.id

-- Join Applicants
JOIN APPLICANT a                ON app.id = a.application_id

-- Join Contacts
LEFT JOIN CONTACT_NUMBER c      ON a.id = c.applicant_id
LEFT JOIN REF_CONTACT_TYPE rct  ON c.type_id = rct.id

WHERE app.member_reference_no = 'APP-CC-002';

