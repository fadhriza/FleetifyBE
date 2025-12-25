# Fleetify Backend

## Database Migration

### Commands

```bash
# Create new table migration
go run cmd/migrate/main.go createtable <table_name>

# Create alter table migration
go run cmd/migrate/main.go altertable <table_name>

# Run pending migrations
go run cmd/migrate/main.go migrate

# Rollback migration
go run cmd/migrate/main.go rollback <migration_file_name>
```

### Workflow

1. Generate migration: `createtable` or `altertable`
2. Edit generated SQL file in `migrations/` directory
3. Run migrations: `migrate`
4. Rollback if needed: `rollback <filename>`

Migrations are tracked in `schema_migrations` table.