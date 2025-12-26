-- Migration: Create table purchasing_details
-- Generated at: 2025-12-26T05:52:55+07:00
-- Generated from model: internal/models/purchasing_details.go

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS purchasing_details (
	purchasing_details_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	id INTEGER,
	purchasing_id INTEGER NOT NULL,
	item_id INTEGER NOT NULL,
	qty INTEGER NOT NULL,
	subtotal NUMERIC(10, 2) NOT NULL,
	created_timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add table and column comments
COMMENT ON TABLE purchasing_details IS 'Table for purchasing_details';
COMMENT ON COLUMN purchasing_details.purchasing_details_id IS 'Primary key UUID';
COMMENT ON COLUMN purchasing_details.created_timestamp IS 'Record creation timestamp';
COMMENT ON COLUMN purchasing_details.updated_timestamp IS 'Record update timestamp';

-- Rollback
-- DROP TABLE IF EXISTS purchasing_details;
