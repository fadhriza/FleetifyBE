-- Migration: Create table items
-- Generated at: 2025-12-26T05:52:42+07:00
-- Generated from model: internal/models/items.go

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS items (
	items_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	id INTEGER,
	name TEXT NOT NULL,
	stock INTEGER,
	price NUMERIC(10, 2) NOT NULL,
	category TEXT,
	unit TEXT,
	min_stock INTEGER,
	created_at TIMESTAMPTZ,
	updated_at TIMESTAMPTZ,
	created_timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add table and column comments
COMMENT ON TABLE items IS 'Table for items';
COMMENT ON COLUMN items.items_id IS 'Primary key UUID';
COMMENT ON COLUMN items.created_timestamp IS 'Record creation timestamp';
COMMENT ON COLUMN items.updated_timestamp IS 'Record update timestamp';

-- Rollback
-- DROP TABLE IF EXISTS items;
