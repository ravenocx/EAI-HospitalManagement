CREATE TABLE patients (
    identity_number BIGINT PRIMARY KEY NOT NULL,
    phone_number VARCHAR(15) NOT NULL,
    name VARCHAR(30) NOT NULL,
    birth_date DATE NOT NULL,
    gender VARCHAR(10) NOT NULL,
    identity_card_scan_img TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for suffix search on phone_number
CREATE INDEX idx_patients_phone_number_suffix ON patients(phone_number text_pattern_ops);

-- Index for wildcard search on name
CREATE INDEX idx_patients_name_trgm ON patients USING gin(name gin_trgm_ops);

-- Index for sorting by created_at
CREATE INDEX idx_patients_created_at ON patients(created_at);