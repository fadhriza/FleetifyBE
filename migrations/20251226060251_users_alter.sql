-- Migration: Alter table users
-- Generated at: 2025-12-26T06:02:51+07:00
-- Generated from model: internal/models/users.go

	ALTER TABLE users ADD COLUMN IF NOT EXISTS username TEXT NOT NULL UNIQUE;
	ALTER TABLE users ADD COLUMN IF NOT EXISTS password TEXT NOT NULL;
	ALTER TABLE users ADD COLUMN IF NOT EXISTS role TEXT NOT NULL;
	ALTER TABLE users ADD COLUMN IF NOT EXISTS full_name TEXT NOT NULL;
	ALTER TABLE users ADD COLUMN IF NOT EXISTS email TEXT;
	ALTER TABLE users ADD COLUMN IF NOT EXISTS phone TEXT;
	ALTER TABLE users ADD COLUMN IF NOT EXISTS is_active BOOLEAN;
	ALTER TABLE users ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ;
	ALTER TABLE users ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ;

-- Rollback
-- Customize rollback statements below (reverse the changes above):
-- ALTER TABLE users DROP COLUMN IF EXISTS column_name;
