-- Migration: Remove unused INTEGER id columns
-- Generated at: 2025-12-26T09:10:00+07:00
-- Purpose: Remove INTEGER id columns that are no longer needed after UUID migration

-- Drop INTEGER id columns from all tables
ALTER TABLE items DROP COLUMN IF EXISTS id;
ALTER TABLE suppliers DROP COLUMN IF EXISTS id;
ALTER TABLE purchasings DROP COLUMN IF EXISTS id;
ALTER TABLE purchasing_details DROP COLUMN IF EXISTS id;
ALTER TABLE users DROP COLUMN IF EXISTS id;

-- Drop sequences that are no longer needed
DROP SEQUENCE IF EXISTS items_id_seq;
DROP SEQUENCE IF EXISTS suppliers_id_seq;
DROP SEQUENCE IF EXISTS purchasings_id_seq;
DROP SEQUENCE IF EXISTS purchasing_details_id_seq;
DROP SEQUENCE IF EXISTS users_id_seq;

-- Rollback (if needed)
-- CREATE SEQUENCE IF NOT EXISTS items_id_seq;
-- CREATE SEQUENCE IF NOT EXISTS suppliers_id_seq;
-- CREATE SEQUENCE IF NOT EXISTS purchasings_id_seq;
-- CREATE SEQUENCE IF NOT EXISTS purchasing_details_id_seq;
-- CREATE SEQUENCE IF NOT EXISTS users_id_seq;
-- ALTER TABLE items ADD COLUMN id INTEGER;
-- ALTER TABLE suppliers ADD COLUMN id INTEGER;
-- ALTER TABLE purchasings ADD COLUMN id INTEGER;
-- ALTER TABLE purchasing_details ADD COLUMN id INTEGER;
-- ALTER TABLE users ADD COLUMN id INTEGER;

