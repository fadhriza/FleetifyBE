-- Migration: Migrate all foreign keys to use UUID
-- Generated at: 2025-12-26T09:00:00+07:00
-- Purpose: Change all foreign keys from INTEGER id to UUID primary keys

-- Step 1: Add new UUID columns for foreign keys
ALTER TABLE purchasings ADD COLUMN IF NOT EXISTS supplier_id_uuid UUID;
ALTER TABLE purchasings ADD COLUMN IF NOT EXISTS user_id_uuid UUID;
ALTER TABLE purchasing_details ADD COLUMN IF NOT EXISTS purchasing_id_uuid UUID;
ALTER TABLE purchasing_details ADD COLUMN IF NOT EXISTS item_id_uuid UUID;

-- Step 2: Migrate data from INTEGER id to UUID
-- Map supplier_id INTEGER to suppliers_id UUID
UPDATE purchasings p
SET supplier_id_uuid = s.suppliers_id
FROM suppliers s
WHERE p.supplier_id = s.id;

-- Map user_id INTEGER to users_id UUID
UPDATE purchasings p
SET user_id_uuid = u.users_id
FROM users u
WHERE p.user_id = u.id;

-- Map purchasing_id INTEGER to purchasings_id UUID
UPDATE purchasing_details pd
SET purchasing_id_uuid = p.purchasings_id
FROM purchasings p
WHERE pd.purchasing_id = p.id;

-- Map item_id INTEGER to items_id UUID
UPDATE purchasing_details pd
SET item_id_uuid = i.items_id
FROM items i
WHERE pd.item_id = i.id;

-- Step 3: Make new UUID columns NOT NULL
ALTER TABLE purchasings ALTER COLUMN supplier_id_uuid SET NOT NULL;
ALTER TABLE purchasings ALTER COLUMN user_id_uuid SET NOT NULL;
ALTER TABLE purchasing_details ALTER COLUMN purchasing_id_uuid SET NOT NULL;
ALTER TABLE purchasing_details ALTER COLUMN item_id_uuid SET NOT NULL;

-- Step 4: Drop old INTEGER foreign key columns
ALTER TABLE purchasings DROP COLUMN IF EXISTS supplier_id;
ALTER TABLE purchasings DROP COLUMN IF EXISTS user_id;
ALTER TABLE purchasing_details DROP COLUMN IF EXISTS purchasing_id;
ALTER TABLE purchasing_details DROP COLUMN IF EXISTS item_id;

-- Step 5: Rename UUID columns to original names
ALTER TABLE purchasings RENAME COLUMN supplier_id_uuid TO supplier_id;
ALTER TABLE purchasings RENAME COLUMN user_id_uuid TO user_id;
ALTER TABLE purchasing_details RENAME COLUMN purchasing_id_uuid TO purchasing_id;
ALTER TABLE purchasing_details RENAME COLUMN item_id_uuid TO item_id;

-- Step 6: Add foreign key constraints with UUID
ALTER TABLE purchasings
ADD CONSTRAINT fk_purchasings_supplier
FOREIGN KEY (supplier_id) REFERENCES suppliers(suppliers_id);

ALTER TABLE purchasings
ADD CONSTRAINT fk_purchasings_user
FOREIGN KEY (user_id) REFERENCES users(users_id);

ALTER TABLE purchasing_details
ADD CONSTRAINT fk_purchasing_details_purchasing
FOREIGN KEY (purchasing_id) REFERENCES purchasings(purchasings_id) ON DELETE CASCADE;

ALTER TABLE purchasing_details
ADD CONSTRAINT fk_purchasing_details_item
FOREIGN KEY (item_id) REFERENCES items(items_id);

-- Rollback (if needed)
-- ALTER TABLE purchasings DROP CONSTRAINT IF EXISTS fk_purchasings_supplier;
-- ALTER TABLE purchasings DROP CONSTRAINT IF EXISTS fk_purchasings_user;
-- ALTER TABLE purchasing_details DROP CONSTRAINT IF EXISTS fk_purchasing_details_purchasing;
-- ALTER TABLE purchasing_details DROP CONSTRAINT IF EXISTS fk_purchasing_details_item;
-- ALTER TABLE purchasings RENAME COLUMN supplier_id TO supplier_id_uuid;
-- ALTER TABLE purchasings RENAME COLUMN user_id TO user_id_uuid;
-- ALTER TABLE purchasing_details RENAME COLUMN purchasing_id TO purchasing_id_uuid;
-- ALTER TABLE purchasing_details RENAME COLUMN item_id TO item_id_uuid;
-- ALTER TABLE purchasings ADD COLUMN supplier_id INTEGER;
-- ALTER TABLE purchasings ADD COLUMN user_id INTEGER;
-- ALTER TABLE purchasing_details ADD COLUMN purchasing_id INTEGER;
-- ALTER TABLE purchasing_details ADD COLUMN item_id INTEGER;

