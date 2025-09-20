-- +migrate Up
CREATE TABLE IF NOT EXISTS tb_payment_configs (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    provider_name VARCHAR(100) NOT NULL,
    provider_code VARCHAR(50) NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT true,
    min_amount NUMERIC(10,2) NOT NULL,
    max_amount NUMERIC(10,2) NOT NULL,
    transaction_fee NUMERIC(10,2) NOT NULL,
    processing_time VARCHAR(50),
    supported_modes TEXT NOT NULL,
    required_fields TEXT NOT NULL,
    provider_config TEXT NOT NULL,
    cache_expiry_mins INTEGER DEFAULT 60
);

CREATE INDEX idx_payment_configs_provider_code ON tb_payment_configs(provider_code);
CREATE INDEX idx_payment_configs_is_active ON tb_payment_configs(is_active);

-- +migrate Down
DROP TABLE IF EXISTS tb_payment_configs;
