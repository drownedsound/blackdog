.mode column
.headers on

-- ============================================================================
-- 1. CONFIGURATION
-- ============================================================================
PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;

-- ============================================================================
-- 2. REFERENCES
-- ============================================================================

CREATE TABLE REF_APPLICATION_STATUS (
    id integer PRIMARY KEY,
    name text NOT NULL UNIQUE,
    description text NOT NULL,
    is_terminal integer NOT NULL DEFAULT 0,
    CONSTRAINT ck_app_status_terminal CHECK (is_terminal IN (0, 1)))
STRICT;

CREATE TABLE REF_CURRENCY (
    id integer PRIMARY KEY,
    code text NOT NULL UNIQUE,
    CONSTRAINT ck_currency_code_len CHECK (length(code) = 3))
STRICT;

CREATE TABLE REF_CREDIT_CARD (
    id integer PRIMARY KEY,
    name text NOT NULL UNIQUE,
    description text NOT NULL,
    product_ceiling integer NOT NULL,
    interest_rate integer NOT NULL,
    CONSTRAINT ck_ref_cc_ceiling CHECK (product_ceiling > 0),
    CONSTRAINT ck_ref_cc_rate CHECK (interest_rate > 0))
STRICT;

CREATE TABLE REF_PERSONAL_LOAN (
    id integer PRIMARY KEY,
    name text NOT NULL UNIQUE,
    description text NOT NULL,
    product_ceiling integer NOT NULL,
    interest_rate integer NOT NULL,
    CONSTRAINT ck_ref_pl_ceiling CHECK (product_ceiling > 0),
    CONSTRAINT ck_ref_pl_rate CHECK (interest_rate > 0))
STRICT;

CREATE TABLE REF_CONTACT_TYPE (
    id integer PRIMARY KEY,
    name text NOT NULL UNIQUE,
    description text NOT NULL)
STRICT;

-- ============================================================================
-- 3. PRODUCT
-- ============================================================================

CREATE TABLE CREDIT_CARD (
    id integer PRIMARY KEY,
    profile_id integer NOT NULL,
    currency_id integer NOT NULL,
    credit_limit integer NOT NULL DEFAULT 0,
    FOREIGN KEY (profile_id) REFERENCES REF_CREDIT_CARD (id),
    FOREIGN KEY (currency_id) REFERENCES REF_CURRENCY (id),
    CONSTRAINT ck_cc_limit CHECK (credit_limit > 0))
STRICT;

CREATE TABLE PERSONAL_LOAN (
    id integer PRIMARY KEY,
    profile_id integer NOT NULL,
    currency_id integer NOT NULL,
    loan_amount integer NOT NULL DEFAULT 0,
    FOREIGN KEY (profile_id) REFERENCES REF_PERSONAL_LOAN (id),
    FOREIGN KEY (currency_id) REFERENCES REF_CURRENCY (id),
    CONSTRAINT ck_pl_amount CHECK (loan_amount > 0))
STRICT;

-- ============================================================================
-- 4. APPLICATION
-- ============================================================================

CREATE TABLE APPLICATION (
    id integer PRIMARY KEY,
    member_reference_no text NOT NULL UNIQUE,
    status_id integer NOT NULL,
    credit_card_id integer,
    personal_loan_id integer,
    requested_amount integer NOT NULL,
    created_at integer NOT NULL DEFAULT (unixepoch ()),
    updated_at integer NOT NULL DEFAULT (unixepoch ()),
    FOREIGN KEY (status_id) REFERENCES REF_APPLICATION_STATUS (id),
    FOREIGN KEY (credit_card_id) REFERENCES CREDIT_CARD (id),
    FOREIGN KEY (personal_loan_id) REFERENCES PERSONAL_LOAN (id),
    CONSTRAINT ck_app_req_amount CHECK (requested_amount > 0),
    CONSTRAINT ck_app_product_exists CHECK (credit_card_id IS NOT NULL OR
	personal_loan_id IS NOT NULL),
    -- Ensure valid epoch time (> Jan 1 2020)
    CONSTRAINT ck_app_created_valid CHECK (created_at > 1577836800),
    CONSTRAINT ck_app_updated_valid CHECK (updated_at > 1577836800))
STRICT;

-- ============================================================================
-- 5. APPLICANT
-- ============================================================================

CREATE TABLE APPLICANT (
    id integer PRIMARY KEY,
    application_id integer NOT NULL,
    is_principal integer NOT NULL DEFAULT 0,
    last_name text NOT NULL,
    first_name text NOT NULL,
    middle_name text,
    birthday text NOT NULL,
    FOREIGN KEY (application_id) REFERENCES APPLICATION (id),
    CONSTRAINT ck_applicant_is_principal CHECK (is_principal IN (0, 1)),
    CONSTRAINT ck_applicant_birthday_valid CHECK (length(birthday) = 10 AND
	birthday IS date(birthday)))
STRICT;

CREATE UNIQUE INDEX idx_applicant_principal ON APPLICANT (application_id)
WHERE
    is_principal = 1;

CREATE INDEX idx_applicant_lookup ON APPLICANT (application_id, is_principal);

-- ============================================================================
-- 6. CONTACT NUMBER
-- ============================================================================

CREATE TABLE CONTACT_NUMBER (
    id integer PRIMARY KEY,
    applicant_id integer NOT NULL,
    type_id integer NOT NULL,
    value text NOT NULL,
    FOREIGN KEY (applicant_id) REFERENCES APPLICANT (id),
    FOREIGN KEY (type_id) REFERENCES REF_CONTACT_TYPE (id),
    CONSTRAINT ck_contact_value_len CHECK (length(value) >= 7))
STRICT;

CREATE INDEX idx_contact_lookup ON CONTACT_NUMBER (applicant_id);
