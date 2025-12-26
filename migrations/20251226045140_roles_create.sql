-- Migration: Create table roles
-- Generated at: 2025-12-26T04:51:40+07:00
-- Generated from model: internal/models/roles.go

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS roles (
	roles_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	role_oid TEXT NOT NULL,
	role_name TEXT NOT NULL,
	role_description TEXT,
	created_timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add table and column comments
COMMENT ON TABLE roles IS 'Table for roles';
COMMENT ON COLUMN roles.roles_id IS 'Primary key UUID';
COMMENT ON COLUMN roles.created_timestamp IS 'Record creation timestamp';
COMMENT ON COLUMN roles.updated_timestamp IS 'Record update timestamp';

-- Rollback
-- DROP TABLE IF EXISTS roles;
