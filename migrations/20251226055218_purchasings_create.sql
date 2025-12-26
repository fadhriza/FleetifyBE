-- Migration: Create table purchasings
-- Generated at: 2025-12-26T05:52:18+07:00
-- Generated from model: internal/models/purchasings.go

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS purchasings (
	purchasings_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	id INTEGER,
	date TIMESTAMPTZ NOT NULL,
	supplier_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	grand_total NUMERIC(10, 2) NOT NULL,
	status TEXT,
	notes TEXT,
	created_at TIMESTAMPTZ,
	created_timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add table and column comments
COMMENT ON TABLE purchasings IS 'Table for purchasings';
COMMENT ON COLUMN purchasings.purchasings_id IS 'Primary key UUID';
COMMENT ON COLUMN purchasings.created_timestamp IS 'Record creation timestamp';
COMMENT ON COLUMN purchasings.updated_timestamp IS 'Record update timestamp';

-- Rollback
-- DROP TABLE IF EXISTS purchasings;
