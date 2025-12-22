.mode column
.headers on

-- ============================================================================
-- 1. CONFIGURATION
-- ============================================================================
PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;

BEGIN TRANSACTION;

-- ============================================================================
-- 2. APPLICATION STATUS
-- ============================================================================

INSERT INTO REF_APPLICATION_STATUS (id, name, description, is_terminal)
VALUES 
    (1, 'CREATED', 'Draft created but not submitted', 0),
    (2, 'IN_PROGRESS', 'Undergoing manual review', 0),
    (3, 'APPROVED', 'Approved but pending booking', 1), 
    (4, 'DECLINED', 'Rejected based on credit policy', 1),
    (5, 'CANCELLED', 'Withdrawn or terminated', 1);

-- ============================================================================
-- 3. PRODUCT DETAILS (REF_CREDIT_CARD)
-- ============================================================================

INSERT INTO REF_CREDIT_CARD (
    id, name, description, product_ceiling, interest_rate
    )
VALUES 
    -- 1. Pasada King
    -- Limit: PHP 75,000.00 
    -- Rate: 3.25% 
    (1, 'Pasada King', 
    'Earn double points on fuel and auto parts', 
    7500000, 325),

    -- 2. Amigas Platinum
    -- Limit: PHP 300,000.00
    -- Rate: 2.25% 
    (2, 'Amigas Platinum', 
    'Massive rebates on dining, salons, and designer products.', 
    30000000, 225),

    -- 3. Gold Steam
    -- Limit: PHP 25,000.00
    -- Rate: 3.50% 
    (3, 'Gold Steam', 
    '5% cashback on Steam purchases', 
    2500000, 350),

    -- 4. Boracay Black
    -- Limit: PHP 1,000,000.00
    -- Rate: 1.75% 
    (4, 'Boracay Black', 
    'No foreign transaction fees. Points convert to free ' || 
    'cocktails at partner resorts.', 
    100000000, 175);

-- ============================================================================
-- 4. PRODUCT DETAILS (REF_PERSONAL_LOAN)
-- ============================================================================

INSERT INTO REF_PERSONAL_LOAN (
    id, name, description, product_ceiling, interest_rate
    )
VALUES
    -- 1. Home Improvement
    -- Limit: PHP 2,000,000.00 
    -- Rate: 14.50% per annum 
    (1, 'Home Improvement Loan',
    'Financing for home renovations, repairs, and construction upgrades.',
    200000000,
    1450
    ),

    -- 2. Medical Emergency Loan
    -- Limit: PHP 500,000
    -- Rate: 12.00% per annum 
    (2, 'Medical Emergency Loan',
    'Urgent cash assistance for hospitalization, surgery, and medical ' ||
    'bills.',
    50000000,
    1200
    ),

    -- 3. Travel Loan
    -- Limit: PHP 300,000
    -- Rate: 18.00% per annum 
    (3, 'Travel Loan',
    'Personal financing specifically for domestic and international ' ||
    'travel expenses.',
    30000000,
    1800
    );

COMMIT;
