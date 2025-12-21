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
    id, name, description, credit_limit, interest_rate
)
VALUES 
    -- 1. Pasada King
    -- Limit: 75,000.00 (7.5M minor units), Rate: 3.25% (325 bps)
    (1, 'Pasada King', 
     'Earn double points on fuel and auto parts', 
     7500000, 325),

    -- 2. Amigas Platinum
    -- Limit: 300,000.00 (30M minor units), Rate: 2.25% (225 bps)
    (2, 'Amigas Platinum', 
     'Massive rebates on dining, salons, and designer products.', 
     30000000, 225),

    -- 3. Gold Steam
    -- Limit: 25,000.00 (2.5M minor units), Rate: 3.50% (350 bps)
    (3, 'Gold Steam', 
     '5% cashback on Steam purchases', 
     2500000, 350),

    -- 4. Boracay Black
    -- Limit: 1,000,000.00 (100M minor units), Rate: 1.75% (175 bps)
    (4, 'Boracay Black', 
     'No foreign transaction fees. Points convert to free ' || 
     'cocktails at partner resorts.', 
     100000000, 175);

COMMIT;
