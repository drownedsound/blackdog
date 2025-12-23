.mode column
.headers on

-- ============================================================================
-- 1. CONFIGURATION
-- ============================================================================
PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;

-- ============================================================================
-- 2. AUDIT LOGGING
-- ============================================================================

CREATE TABLE AUDIT_LOG (
    id integer PRIMARY KEY,
    table_name text NOT NULL,
    record_id integer NOT NULL,
    action text NOT NULL, -- 'INSERT', 'UPDATE'
    timestamp integer NOT NULL DEFAULT (unixepoch()),
    changes text NOT NULL -- JSON representation of the record
) STRICT;

-- ============================================================================
-- 3. REFERENCES
-- ============================================================================

CREATE TABLE REF_APPLICATION_STATUS (
    id integer PRIMARY KEY,
    name text NOT NULL UNIQUE,
    description text NOT NULL,
    is_terminal integer NOT NULL DEFAULT 0,
    CONSTRAINT ck_app_status_terminal CHECK (is_terminal IN (0, 1)))
STRICT;

CREATE TRIGGER trg_block_del_ref_app_status BEFORE DELETE ON REF_APPLICATION_STATUS
BEGIN
    SELECT RAISE(ABORT, 'Deletion is not allowed on REF_APPLICATION_STATUS');
END;

CREATE TABLE REF_CURRENCY (
    id integer PRIMARY KEY,
    code text NOT NULL UNIQUE,
    CONSTRAINT ck_currency_code_len CHECK (length(code) = 3))
STRICT;

CREATE TRIGGER trg_block_del_ref_currency BEFORE DELETE ON REF_CURRENCY
BEGIN
    SELECT RAISE(ABORT, 'Deletion is not allowed on REF_CURRENCY');
END;

CREATE TABLE REF_CREDIT_CARD (
    id integer PRIMARY KEY,
    name text NOT NULL UNIQUE,
    description text NOT NULL,
    currency_id integer NOT NULL,
    product_ceiling integer NOT NULL,
    interest_rate integer NOT NULL,
    FOREIGN KEY (currency_id) REFERENCES REF_CURRENCY (id),
    CONSTRAINT ck_ref_cc_ceiling CHECK (product_ceiling > 0),
    CONSTRAINT ck_ref_cc_rate CHECK (interest_rate > 0))
STRICT;

CREATE TRIGGER trg_block_del_ref_cc BEFORE DELETE ON REF_CREDIT_CARD
BEGIN
    SELECT RAISE(ABORT, 'Deletion is not allowed on REF_CREDIT_CARD');
END;

CREATE TABLE REF_PERSONAL_LOAN (
    id integer PRIMARY KEY,
    name text NOT NULL UNIQUE,
    description text NOT NULL,
    currency_id integer NOT NULL,
    product_ceiling integer NOT NULL,
    interest_rate integer NOT NULL,
    FOREIGN KEY (currency_id) REFERENCES REF_CURRENCY (id),
    CONSTRAINT ck_ref_pl_ceiling CHECK (product_ceiling > 0),
    CONSTRAINT ck_ref_pl_rate CHECK (interest_rate > 0))
STRICT;

CREATE TRIGGER trg_block_del_ref_pl BEFORE DELETE ON REF_PERSONAL_LOAN
BEGIN
    SELECT RAISE(ABORT, 'Deletion is not allowed on REF_PERSONAL_LOAN');
END;

CREATE TABLE REF_CONTACT_TYPE (
    id integer PRIMARY KEY,
    name text NOT NULL UNIQUE,
    description text NOT NULL)
STRICT;

CREATE TRIGGER trg_block_del_ref_contact_type BEFORE DELETE ON REF_CONTACT_TYPE
BEGIN
    SELECT RAISE(ABORT, 'Deletion is not allowed on REF_CONTACT_TYPE');
END;

-- ============================================================================
-- 4. PRODUCT
-- ============================================================================

CREATE TABLE CREDIT_CARD (
    id integer PRIMARY KEY,
    profile_id integer NOT NULL,
    credit_limit integer NOT NULL DEFAULT 0,
    FOREIGN KEY (profile_id) REFERENCES REF_CREDIT_CARD (id),
    CONSTRAINT ck_cc_limit CHECK (credit_limit > 0))
STRICT;

CREATE TRIGGER trg_audit_ins_cc AFTER INSERT ON CREDIT_CARD
BEGIN
    INSERT INTO AUDIT_LOG (table_name, record_id, action, changes)
    VALUES ('CREDIT_CARD', new.id, 'INSERT', 
        json_object(
            'profile_id', new.profile_id, 
            'credit_limit', new.credit_limit
        )
    );
END;

CREATE TRIGGER trg_audit_upd_cc AFTER UPDATE ON CREDIT_CARD
BEGIN
    INSERT INTO AUDIT_LOG (table_name, record_id, action, changes)
    VALUES ('CREDIT_CARD', new.id, 'UPDATE', 
        json_object(
            'profile_id', new.profile_id, 
            'credit_limit', new.credit_limit
        )
    );
END;

CREATE TRIGGER trg_block_del_cc BEFORE DELETE ON CREDIT_CARD
BEGIN
    SELECT RAISE(ABORT, 'Deletion is not allowed on CREDIT_CARD');
END;

CREATE TABLE PERSONAL_LOAN (
    id integer PRIMARY KEY,
    profile_id integer NOT NULL,
    loan_amount integer NOT NULL DEFAULT 0,
    FOREIGN KEY (profile_id) REFERENCES REF_PERSONAL_LOAN (id),
    CONSTRAINT ck_pl_amount CHECK (loan_amount > 0))
STRICT;

CREATE TRIGGER trg_audit_ins_pl AFTER INSERT ON PERSONAL_LOAN
BEGIN
    INSERT INTO AUDIT_LOG (table_name, record_id, action, changes)
    VALUES ('PERSONAL_LOAN', new.id, 'INSERT', 
        json_object(
            'profile_id', new.profile_id, 
            'loan_amount', new.loan_amount
        )
    );
END;

CREATE TRIGGER trg_audit_upd_pl AFTER UPDATE ON PERSONAL_LOAN
BEGIN
    INSERT INTO AUDIT_LOG (table_name, record_id, action, changes)
    VALUES ('PERSONAL_LOAN', new.id, 'UPDATE', 
        json_object(
            'profile_id', new.profile_id, 
            'loan_amount', new.loan_amount
        )
    );
END;

CREATE TRIGGER trg_block_del_pl BEFORE DELETE ON PERSONAL_LOAN
BEGIN
    SELECT RAISE(ABORT, 'Deletion is not allowed on PERSONAL_LOAN');
END;

-- ============================================================================
-- 5. APPLICATION
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
    CONSTRAINT ck_app_created_valid CHECK (created_at > 1577836800),
    CONSTRAINT ck_app_updated_valid CHECK (updated_at > 1577836800))
STRICT;

CREATE TRIGGER trg_audit_ins_app AFTER INSERT ON APPLICATION
BEGIN
    INSERT INTO AUDIT_LOG (table_name, record_id, action, changes)
    VALUES ('APPLICATION', new.id, 'INSERT', 
        json_object(
            'ref_no', new.member_reference_no, 
            'status', new.status_id,
            'cc_id', new.credit_card_id,
            'pl_id', new.personal_loan_id,
            'req_amt', new.requested_amount
        )
    );
END;

CREATE TRIGGER trg_audit_upd_app AFTER UPDATE ON APPLICATION
BEGIN
    INSERT INTO AUDIT_LOG (table_name, record_id, action, changes)
    VALUES ('APPLICATION', new.id, 'UPDATE', 
        json_object(
            'ref_no', new.member_reference_no, 
            'status', new.status_id,
            'cc_id', new.credit_card_id,
            'pl_id', new.personal_loan_id,
            'req_amt', new.requested_amount
        )
    );
END;

CREATE TRIGGER trg_block_del_app BEFORE DELETE ON APPLICATION
BEGIN
    SELECT RAISE(ABORT, 'Deletion is not allowed on APPLICATION');
END;

-- ============================================================================
-- 6. APPLICANT
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

CREATE TRIGGER trg_audit_ins_applicant AFTER INSERT ON APPLICANT
BEGIN
    INSERT INTO AUDIT_LOG (table_name, record_id, action, changes)
    VALUES ('APPLICANT', new.id, 'INSERT', 
        json_object(
            'app_id', new.application_id, 
            'is_principal', new.is_principal,
            'last', new.last_name,
            'first', new.first_name
        )
    );
END;

CREATE TRIGGER trg_audit_upd_applicant AFTER UPDATE ON APPLICANT
BEGIN
    INSERT INTO AUDIT_LOG (table_name, record_id, action, changes)
    VALUES ('APPLICANT', new.id, 'UPDATE', 
        json_object(
            'app_id', new.application_id, 
            'is_principal', new.is_principal,
            'last', new.last_name,
            'first', new.first_name
        )
    );
END;

CREATE TRIGGER trg_block_del_applicant BEFORE DELETE ON APPLICANT
BEGIN
    SELECT RAISE(ABORT, 'Deletion is not allowed on APPLICANT');
END;

-- ============================================================================
-- 7. CONTACT NUMBER
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

CREATE TRIGGER trg_audit_ins_contact AFTER INSERT ON CONTACT_NUMBER
BEGIN
    INSERT INTO AUDIT_LOG (table_name, record_id, action, changes)
    VALUES ('CONTACT_NUMBER', new.id, 'INSERT', 
        json_object(
            'applicant_id', new.applicant_id, 
            'type', new.type_id,
            'val', new.value
        )
    );
END;

CREATE TRIGGER trg_audit_upd_contact AFTER UPDATE ON CONTACT_NUMBER
BEGIN
    INSERT INTO AUDIT_LOG (table_name, record_id, action, changes)
    VALUES ('CONTACT_NUMBER', new.id, 'UPDATE', 
        json_object(
            'applicant_id', new.applicant_id, 
            'type', new.type_id,
            'val', new.value
        )
    );
END;

CREATE TRIGGER trg_block_del_contact BEFORE DELETE ON CONTACT_NUMBER
BEGIN
    SELECT RAISE(ABORT, 'Deletion is not allowed on CONTACT_NUMBER');
END;
