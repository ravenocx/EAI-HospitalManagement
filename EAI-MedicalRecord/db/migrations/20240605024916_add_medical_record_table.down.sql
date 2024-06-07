DROP TABLE IF EXISTS medical_records CASCADE;

DROP INDEX IF EXISTS idx_medical_records_identity_number CASCADE;
DROP INDEX IF EXISTS idx_medical_records_created_by_nip CASCADE;
DROP INDEX IF EXISTS idx_medical_records_created_by_user_id CASCADE;
DROP INDEX IF EXISTS idx_medical_records_created_at CASCADE;