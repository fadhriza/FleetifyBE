-- Migration: Add id column to users table
-- Generated at: 2025-12-26T07:00:00+07:00

ALTER TABLE users ADD COLUMN IF NOT EXISTS id INTEGER;

-- Rollback
-- ALTER TABLE users DROP COLUMN IF EXISTS id;

