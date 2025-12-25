# Fleetify Backend

## Database Migration

### Commands

```bash
# Create model template (prompts for seeder)
go run cmd/migrate/main.go createtable <table_name>

# Generate SQL migration from model
go run cmd/migrate/main.go generatesql <table_name>

# Generate ALTER migration from model
go run cmd/migrate/main.go altertable <table_name>

# Run pending migrations
go run cmd/migrate/main.go migrate

# Rollback migration
go run cmd/migrate/main.go rollback <migration_file_name>

# Run seeders
go run cmd/migrate/main.go seed <table_name>
```

### Workflow

**CREATE:**
1. Create model template: `createtable <table_name>` (prompts: "Do you need dbseeds? [y/n]")
2. Edit model in `internal/models/<table_name>.go` (add seed data if y)
3. Generate SQL: `generatesql <table_name>` (reads model, generates SQL)
4. Run migrations: `migrate`
5. Run seeders: `seed <table_name>` (if seeders enabled)
6. Rollback if needed: `rollback <filename>`

**ALTER:**
1. Edit model in `internal/models/<table_name>.go`
2. Generate ALTER migration: `altertable <table_name>` (reads model, generates ALTER SQL)
3. Run migrations: `migrate`
4. Rollback if needed: `rollback <filename>`

Migrations are tracked in `schema_migrations` table.