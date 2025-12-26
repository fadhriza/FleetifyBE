-- Migration: Fix ID sequences for all tables
-- Generated at: 2025-12-26T08:00:00+07:00
-- Purpose: Add sequences and auto-increment for id columns

-- Create sequences for each table
CREATE SEQUENCE IF NOT EXISTS items_id_seq;
CREATE SEQUENCE IF NOT EXISTS suppliers_id_seq;
CREATE SEQUENCE IF NOT EXISTS purchasings_id_seq;
CREATE SEQUENCE IF NOT EXISTS purchasing_details_id_seq;
CREATE SEQUENCE IF NOT EXISTS users_id_seq;

-- Set sequence ownership
ALTER SEQUENCE items_id_seq OWNED BY items.id;
ALTER SEQUENCE suppliers_id_seq OWNED BY suppliers.id;
ALTER SEQUENCE purchasings_id_seq OWNED BY purchasings.id;
ALTER SEQUENCE purchasing_details_id_seq OWNED BY purchasing_details.id;
ALTER SEQUENCE users_id_seq OWNED BY users.id;

-- Update existing id values to use sequence (only for rows with id = 0 or NULL)
UPDATE items SET id = nextval('items_id_seq') WHERE id IS NULL OR id = 0;
UPDATE suppliers SET id = nextval('suppliers_id_seq') WHERE id IS NULL OR id = 0;
UPDATE purchasings SET id = nextval('purchasings_id_seq') WHERE id IS NULL OR id = 0;
UPDATE purchasing_details SET id = nextval('purchasing_details_id_seq') WHERE id IS NULL OR id = 0;
UPDATE users SET id = nextval('users_id_seq') WHERE id IS NULL OR id = 0;

-- Set sequence to current max value + 1 to avoid conflicts
SELECT setval('items_id_seq', COALESCE((SELECT MAX(id) FROM items), 0) + 1, false);
SELECT setval('suppliers_id_seq', COALESCE((SELECT MAX(id) FROM suppliers), 0) + 1, false);
SELECT setval('purchasings_id_seq', COALESCE((SELECT MAX(id) FROM purchasings), 0) + 1, false);
SELECT setval('purchasing_details_id_seq', COALESCE((SELECT MAX(id) FROM purchasing_details), 0) + 1, false);
SELECT setval('users_id_seq', COALESCE((SELECT MAX(id) FROM users), 0) + 1, false);

-- Alter columns to use sequence as default
ALTER TABLE items ALTER COLUMN id SET DEFAULT nextval('items_id_seq');
ALTER TABLE suppliers ALTER COLUMN id SET DEFAULT nextval('suppliers_id_seq');
ALTER TABLE purchasings ALTER COLUMN id SET DEFAULT nextval('purchasings_id_seq');
ALTER TABLE purchasing_details ALTER COLUMN id SET DEFAULT nextval('purchasing_details_id_seq');
ALTER TABLE users ALTER COLUMN id SET DEFAULT nextval('users_id_seq');

-- Make id columns NOT NULL
ALTER TABLE items ALTER COLUMN id SET NOT NULL;
ALTER TABLE suppliers ALTER COLUMN id SET NOT NULL;
ALTER TABLE purchasings ALTER COLUMN id SET NOT NULL;
ALTER TABLE purchasing_details ALTER COLUMN id SET NOT NULL;
ALTER TABLE users ALTER COLUMN id SET NOT NULL;

-- Rollback
-- ALTER TABLE items ALTER COLUMN id DROP DEFAULT;
-- ALTER TABLE suppliers ALTER COLUMN id DROP DEFAULT;
-- ALTER TABLE purchasings ALTER COLUMN id DROP DEFAULT;
-- ALTER TABLE purchasing_details ALTER COLUMN id DROP DEFAULT;
-- ALTER TABLE users ALTER COLUMN id DROP DEFAULT;
-- DROP SEQUENCE IF EXISTS items_id_seq;
-- DROP SEQUENCE IF EXISTS suppliers_id_seq;
-- DROP SEQUENCE IF EXISTS purchasings_id_seq;
-- DROP SEQUENCE IF EXISTS purchasing_details_id_seq;
-- DROP SEQUENCE IF EXISTS users_id_seq;

