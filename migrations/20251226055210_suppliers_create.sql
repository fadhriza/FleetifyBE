-- Migration: Create table suppliers
-- Generated at: 2025-12-26T05:52:10+07:00
-- Generated from model: internal/models/suppliers.go

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS suppliers (
	suppliers_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	id INTEGER,
	name TEXT NOT NULL,
	email TEXT,
	address TEXT,
	phone TEXT,
	supplier_type TEXT,
	is_active BOOLEAN,
	created_at TIMESTAMPTZ,
	updated_at TIMESTAMPTZ,
	created_timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add table and column comments
COMMENT ON TABLE suppliers IS 'Table for suppliers';
COMMENT ON COLUMN suppliers.suppliers_id IS 'Primary key UUID';
COMMENT ON COLUMN suppliers.created_timestamp IS 'Record creation timestamp';
COMMENT ON COLUMN suppliers.updated_timestamp IS 'Record update timestamp';

-- Rollback
-- DROP TABLE IF EXISTS suppliers;
