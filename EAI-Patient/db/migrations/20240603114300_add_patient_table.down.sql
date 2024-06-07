DROP TABLE IF EXISTS patients CASCADE;

DROP INDEX IF EXISTS idx_patients_phone_number_suffix CASCADE;
DROP INDEX IF EXISTS idx_patients_name_trgm CASCADE;
DROP INDEX IF EXISTS idx_patients_created_at CASCADE;