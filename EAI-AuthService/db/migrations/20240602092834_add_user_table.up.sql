CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

CREATE TABLE users (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    nip VARCHAR(20) NOT NULL,
    name VARCHAR(50) NOT NULL,
    password VARCHAR(100),
    role VARCHAR(20) NOT NULL, -- Admin, Nurse.
    identity_card_scan_img TEXT,
    access BOOLEAN NOT NULL DEFAULT true,
    refresh_token VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL
);

-- Index for normal WHERE = searches on nip
CREATE INDEX idx_users_nip ON users(nip);

-- Index for wildcard searches on name
CREATE INDEX idx_users_name ON users USING GIN (name gin_trgm_ops);

-- Index for normal WHERE = searches on role
CREATE INDEX idx_users_role ON users(role);

-- Index for sorting by created_at
CREATE INDEX idx_users_created_at ON users(created_at);