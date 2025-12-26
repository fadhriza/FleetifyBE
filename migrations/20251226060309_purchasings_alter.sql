-- Migration: Alter table purchasings
-- Generated at: 2025-12-26T06:03:09+07:00
-- Generated from model: internal/models/purchasings.go

	ALTER TABLE purchasings ADD COLUMN IF NOT EXISTS id INTEGER;
	ALTER TABLE purchasings ADD COLUMN IF NOT EXISTS date TIMESTAMPTZ NOT NULL;
	ALTER TABLE purchasings ADD COLUMN IF NOT EXISTS supplier_id INTEGER NOT NULL;
	ALTER TABLE purchasings ADD COLUMN IF NOT EXISTS user_id INTEGER NOT NULL;
	ALTER TABLE purchasings ADD COLUMN IF NOT EXISTS grand_total NUMERIC(10, 2) NOT NULL;
	ALTER TABLE purchasings ADD COLUMN IF NOT EXISTS status TEXT;
	ALTER TABLE purchasings ADD COLUMN IF NOT EXISTS notes TEXT;
	ALTER TABLE purchasings ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ;

-- Rollback
-- Customize rollback statements below (reverse the changes above):
-- ALTER TABLE purchasings DROP COLUMN IF EXISTS column_name;
