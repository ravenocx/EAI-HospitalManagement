CREATE TABLE medical_records (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    identity_number BIGINT NOT NULL REFERENCES patients(identity_number),
    symptoms TEXT NOT NULL,
    medications TEXT NOT NULL,
    created_by_nip VARCHAR(20) NOT NULL,
    created_by_name VARCHAR(50) NOT NULL,
    created_by_user_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_medical_records_identity_number ON medical_records(identity_number);
CREATE INDEX idx_medical_records_created_by_nip ON medical_records(created_by_nip);
CREATE INDEX idx_medical_records_created_by_user_id ON medical_records(created_by_user_id);
CREATE INDEX idx_medical_records_created_at ON medical_records(created_at);