-- Migration: Alter table suppliers
-- Generated at: 2025-12-26T06:03:00+07:00
-- Generated from model: internal/models/suppliers.go

	ALTER TABLE suppliers ADD COLUMN IF NOT EXISTS id INTEGER;
	ALTER TABLE suppliers ADD COLUMN IF NOT EXISTS name TEXT NOT NULL;
	ALTER TABLE suppliers ADD COLUMN IF NOT EXISTS email TEXT;
	ALTER TABLE suppliers ADD COLUMN IF NOT EXISTS address TEXT;
	ALTER TABLE suppliers ADD COLUMN IF NOT EXISTS phone TEXT;
	ALTER TABLE suppliers ADD COLUMN IF NOT EXISTS supplier_type TEXT;
	ALTER TABLE suppliers ADD COLUMN IF NOT EXISTS is_active BOOLEAN;
	ALTER TABLE suppliers ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ;
	ALTER TABLE suppliers ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ;

-- Rollback
-- Customize rollback statements below (reverse the changes above):
-- ALTER TABLE suppliers DROP COLUMN IF EXISTS column_name;
