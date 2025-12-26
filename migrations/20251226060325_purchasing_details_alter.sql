-- Migration: Alter table purchasing_details
-- Generated at: 2025-12-26T06:03:25+07:00
-- Generated from model: internal/models/purchasing_details.go

	ALTER TABLE purchasing_details ADD COLUMN IF NOT EXISTS id INTEGER;
	ALTER TABLE purchasing_details ADD COLUMN IF NOT EXISTS purchasing_id INTEGER NOT NULL;
	ALTER TABLE purchasing_details ADD COLUMN IF NOT EXISTS item_id INTEGER NOT NULL;
	ALTER TABLE purchasing_details ADD COLUMN IF NOT EXISTS qty INTEGER NOT NULL;
	ALTER TABLE purchasing_details ADD COLUMN IF NOT EXISTS subtotal NUMERIC(10, 2) NOT NULL;

-- Rollback
-- Customize rollback statements below (reverse the changes above):
-- ALTER TABLE purchasing_details DROP COLUMN IF EXISTS column_name;
