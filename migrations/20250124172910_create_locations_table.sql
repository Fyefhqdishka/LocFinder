-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS locations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ip_address VARCHAR(45) UNIQUE NOT NULL,
    country VARCHAR(100),
    city VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_ip_country_city ON locations(ip_address, country, city);
CREATE INDEX idx_created_at ON locations(created_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP EXTENSION IF EXISTS "uuid-ossp";
DROP TABLE IF EXISTS locations;
-- +goose StatementEnd
