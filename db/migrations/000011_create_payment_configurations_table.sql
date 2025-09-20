-- +migrate Up
CREATE TABLE IF NOT EXISTS tb_payment_configurations (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    country_code VARCHAR(10) NOT NULL UNIQUE,
    secret_key TEXT NOT NULL,
    public_key TEXT NOT NULL
);

CREATE INDEX idx_payment_configurations_country_code ON tb_payment_configurations(country_code);

-- +migrate Down
DROP TABLE IF EXISTS tb_payment_configurations;
