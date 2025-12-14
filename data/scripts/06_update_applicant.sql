BEGIN TRANSACTION;

-- 1. Perform the Update
UPDATE APPLICATION 
SET status_id = 2
WHERE member_reference_no = 'APP-2023-SINGLE';

COMMIT;

-- -- 2. Verify Output (Optional Display)
-- SELECT 'Update Complete. Checking Audit Log...' as info;
--
-- SELECT 
--     id as audit_id, 
--     record_id, 
--     action, 
--     old_value, 
--     new_value, 
--     changed_at 
-- FROM AUDIT_LOG 
-- WHERE table_name = 'APPLICATION' 
-- AND record_id = (
--     SELECT id FROM APPLICATION WHERE member_reference_no = 'APP-2023-SINGLE'
-- )
-- ORDER BY id DESC LIMIT 1;
