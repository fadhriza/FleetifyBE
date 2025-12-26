-- Migration: Create table users
-- Generated at: 2025-12-26T03:46:23+07:00
-- Generated from model: internal/models/users.go

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
	users_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	username TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	role TEXT NOT NULL,
	full_name TEXT NOT NULL,
	email TEXT,
	phone TEXT,
	is_active BOOLEAN,
	created_at TIMESTAMPTZ,
	updated_at TIMESTAMPTZ,
	created_timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add table and column comments
COMMENT ON TABLE users IS 'Table for users';
COMMENT ON COLUMN users.users_id IS 'Primary key UUID';
COMMENT ON COLUMN users.created_timestamp IS 'Record creation timestamp';
COMMENT ON COLUMN users.updated_timestamp IS 'Record update timestamp';

-- Rollback
-- DROP TABLE IF EXISTS users;
