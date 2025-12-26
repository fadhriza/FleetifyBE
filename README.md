# Fleetify Backend

## Getting Started

### Prerequisites
Before running the server, create a `.env` or `.env.development` file with the required environment variables (see Environment Variables section below).

### Commands
```bash
# Run in development mode (reloads on file changes)
go run cmd/server/main.go

# Run in production mode (recommended: build binary first)
go build -o fleetify cmd/server/main.go
./fleetify
```

**Note:** The application loads environment variables from `.env.{ENV}` (where ENV defaults to "development") or falls back to `.env` file. Make sure to configure your environment variables before running the server.


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

## Environment Variables

Create a `.env` or `.env.development` file in the project root with the following variables:

### Required Variables

**Server Configuration:**
```bash
PORT=3000
HOST=0.0.0.0
ENV=development
```

**Database Configuration:**
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=fleetify
DB_SSLMODE=disable
```

**JWT Configuration:**
```bash
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRES_IN=24h
```

**CORS Configuration:**
```bash
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173,http://127.0.0.1:5500
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,PATCH,OPTIONS
CORS_ALLOWED_HEADERS=Origin,Content-Type,Accept,Authorization
```

**Webhook Configuration (Optional):**
```bash
WEBHOOK_URL=https://webhook.site/79dadcf7-3cd0-4601-9efb-fdcdbd6a7568
```

**Note:** 
- The application loads from `.env.{ENV}` file first (e.g., `.env.development`), then falls back to `.env`, then system environment variables.
- Default values are used if variables are not set (see `internal/config/config.go` for defaults).
- When a purchasing is created, the API will send a webhook notification to the configured URL asynchronously (if `WEBHOOK_URL` is set).