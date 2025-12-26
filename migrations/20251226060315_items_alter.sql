-- Migration: Alter table items
-- Generated at: 2025-12-26T06:03:15+07:00
-- Generated from model: internal/models/items.go

	ALTER TABLE items ADD COLUMN IF NOT EXISTS id INTEGER;
	ALTER TABLE items ADD COLUMN IF NOT EXISTS name TEXT NOT NULL;
	ALTER TABLE items ADD COLUMN IF NOT EXISTS stock INTEGER;
	ALTER TABLE items ADD COLUMN IF NOT EXISTS price NUMERIC(10, 2) NOT NULL;
	ALTER TABLE items ADD COLUMN IF NOT EXISTS category TEXT;
	ALTER TABLE items ADD COLUMN IF NOT EXISTS unit TEXT;
	ALTER TABLE items ADD COLUMN IF NOT EXISTS min_stock INTEGER;
	ALTER TABLE items ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ;
	ALTER TABLE items ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ;

-- Rollback
-- Customize rollback statements below (reverse the changes above):
-- ALTER TABLE items DROP COLUMN IF EXISTS column_name;
