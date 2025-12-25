package migration

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Generator struct {
	migrationsDir string
	modelsDir     string
}

func NewGenerator() *Generator {
	return &Generator{
		migrationsDir: "migrations",
		modelsDir:     "internal/models",
	}
}

func (g *Generator) CreateTableModel(tableName string, needsSeeder bool) error {
	if err := g.ensureDirectories(); err != nil {
		return fmt.Errorf("failed to ensure directories: %w", err)
	}

	tableNameLower := strings.ToLower(tableName)
	idFieldName := fmt.Sprintf("%s_id", tableNameLower)
	modelPath := filepath.Join(g.modelsDir, fmt.Sprintf("%s.go", tableNameLower))

	if _, err := os.Stat(modelPath); err == nil {
		return fmt.Errorf("model already exists: %s", modelPath)
	}

	modelName := g.toPascalCase(tableName)
	modelContent := g.generateModelContent(modelName, tableNameLower, idFieldName, needsSeeder)
	if err := os.WriteFile(modelPath, []byte(modelContent), 0644); err != nil {
		return fmt.Errorf("failed to write model file: %w", err)
	}

	return nil
}

func (g *Generator) GenerateSQLMigration(tableName string) error {
	if err := g.ensureDirectories(); err != nil {
		return fmt.Errorf("failed to ensure directories: %w", err)
	}

	tableNameLower := strings.ToLower(tableName)
	idFieldName := fmt.Sprintf("%s_id", tableNameLower)
	modelPath := filepath.Join(g.modelsDir, fmt.Sprintf("%s.go", tableNameLower))

	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return fmt.Errorf("model file not found: %s. Create model first with: createtable %s", modelPath, tableName)
	}

	sqlColumns, err := g.GenerateSQLFromModel(tableName)
	if err != nil {
		return fmt.Errorf("failed to generate SQL from model: %w", err)
	}

	timestamp := time.Now().Format("20060102150405")
	migrationFileName := fmt.Sprintf("%s_%s_%s.sql", timestamp, tableNameLower, "create")
	migrationPath := filepath.Join(g.migrationsDir, migrationFileName)

	migrationContent := g.generateCreateMigrationFromModel(tableNameLower, idFieldName, sqlColumns)

	if err := os.WriteFile(migrationPath, []byte(migrationContent), 0644); err != nil {
		return fmt.Errorf("failed to write migration file: %w", err)
	}

	return nil
}

func (g *Generator) GenerateAlterTable(tableName string) error {
	if err := g.ensureDirectories(); err != nil {
		return fmt.Errorf("failed to ensure directories: %w", err)
	}

	tableNameLower := strings.ToLower(tableName)
	modelPath := filepath.Join(g.modelsDir, fmt.Sprintf("%s.go", tableNameLower))

	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return fmt.Errorf("model file not found: %s. Create model first", modelPath)
	}

	sqlColumns, err := g.GenerateSQLFromModel(tableName)
	if err != nil {
		return fmt.Errorf("failed to generate SQL from model: %w", err)
	}

	timestamp := time.Now().Format("20060102150405")
	migrationFileName := fmt.Sprintf("%s_%s_%s.sql", timestamp, tableNameLower, "alter")
	migrationPath := filepath.Join(g.migrationsDir, migrationFileName)

	migrationContent := g.generateAlterMigrationFromModel(tableNameLower, sqlColumns)

	if err := os.WriteFile(migrationPath, []byte(migrationContent), 0644); err != nil {
		return fmt.Errorf("failed to write migration file: %w", err)
	}

	return nil
}

func (g *Generator) ensureDirectories() error {
	if err := os.MkdirAll(g.migrationsDir, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(g.modelsDir, 0755); err != nil {
		return err
	}
	return nil
}

func (g *Generator) toPascalCase(s string) string {
	parts := strings.Split(s, "_")
	var result strings.Builder
	for _, part := range parts {
		if len(part) > 0 {
			result.WriteString(strings.ToUpper(part[:1]) + strings.ToLower(part[1:]))
		}
	}
	return result.String()
}

func (g *Generator) generateModelContent(modelName, tableName, idFieldName string, needsSeeder bool) string {
	idFieldNamePascal := g.toPascalCase(idFieldName)

	imports := `import (
	"time"
)`
	if needsSeeder {
		imports = `import (
	"fleetify/internal/migration"
	"time"
)`
	}

	baseContent := fmt.Sprintf(`package models

%s

// %s represents the %s table
type %s struct {
	// PK UUID v4 
	%s string `+"`db:\"%s\" json:\"%s\"`"+`


	// 
	// Example: Text field (nullable)
	// FieldName string `+"`db:\"field_name\" json:\"field_name\"`"+`
	//
	// Example: Text field (NOT NULL)
	// FieldName string `+"`db:\"field_name,notnull\" json:\"field_name\"`"+`
	//
	// Example: Text field (UNIQUE)
	// FieldName string `+"`db:\"field_name,unique\" json:\"field_name\"`"+`
	//
	// Example: JSONB field
	// Data map[string]interface{} `+"`db:\"data\" json:\"data\"`"+`
	//
	// Example: Numeric field
	// Amount float64 `+"`db:\"amount\" json:\"amount\"`"+`
	//
	// Example: Boolean field with default
	// IsActive bool `+"`db:\"is_active\" json:\"is_active\"`"+`
	//
	// Example: Timestamp field
	// EventDate time.Time `+"`db:\"event_date\" json:\"event_date\"`"+`

	// Timestamps
	CreatedTimestamp time.Time `+"`db:\"created_timestamp\" json:\"created_timestamp\"`"+`
	UpdatedTimestamp time.Time `+"`db:\"updated_timestamp\" json:\"updated_timestamp\"`"+`
}

// TableName returns the table name
func (%s) TableName() string {
	return "%s"
}

// GetID returns the primary key field name
func (%s) GetID() string {
	return "%s"
}
`, imports, modelName, tableName, modelName, idFieldNamePascal, idFieldName, idFieldName, modelName, tableName, modelName, idFieldName)

	if needsSeeder {
		seederContent := fmt.Sprintf(`

func init() {
	migration.RegisterSeeder("%s", func() interface{} {
		return Seed%s()
	})
}

// Seed%s seeds the %s table
func Seed%s() []%s {
	return []%s{
		//set here
		// Example seeder data:
		// {
		// 	%s: "uuid-here",
		// 	FieldName: "value",
		// 	CreatedTimestamp: time.Now(),
		// 	UpdatedTimestamp: time.Now(),
		// },
	}
}
`, modelName, modelName, modelName, tableName, modelName, modelName, modelName, idFieldNamePascal)
		return baseContent + seederContent
	}

	return baseContent
}

func (g *Generator) generateMigrationContent(tableName, idFieldName string) string {
	return fmt.Sprintf(`-- Migration: Create table %s
-- Generated at: %s

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS %s (
	%s UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	

	--
	-- Text columns:
	-- name TEXT NOT NULL,
	-- email TEXT UNIQUE,
	-- description TEXT,
	-- status TEXT DEFAULT 'active' NOT NULL,
	--
	-- Numeric columns:
	-- amount NUMERIC(10, 2) NOT NULL,
	-- quantity INTEGER NOT NULL DEFAULT 0,
	-- price DECIMAL(12, 2),
	--
	-- JSONB columns:
	-- data JSONB DEFAULT '{}'::jsonb,
	-- metadata JSONB,
	--
	-- Boolean columns:
	-- is_active BOOLEAN NOT NULL DEFAULT false,
	-- is_deleted BOOLEAN NOT NULL DEFAULT false,
	--
	-- Timestamp columns:
	-- event_date TIMESTAMPTZ,
	-- deleted_at TIMESTAMPTZ,
	--
	-- Foreign keys:
	-- user_id UUID REFERENCES users(user_id),
	-- category_id UUID REFERENCES categories(category_id) ON DELETE CASCADE,
	
	created_timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes (customize as needed)
-- CREATE INDEX IF NOT EXISTS idx_%s_created_at ON %s(created_timestamp);
-- CREATE INDEX IF NOT EXISTS idx_%s_status ON %s(status);
-- CREATE INDEX IF NOT EXISTS idx_%s_user_id ON %s(user_id);

-- Add table and column comments
COMMENT ON TABLE %s IS 'Table for %s';
COMMENT ON COLUMN %s.%s IS 'Primary key UUID';
COMMENT ON COLUMN %s.created_timestamp IS 'Record creation timestamp';
COMMENT ON COLUMN %s.updated_timestamp IS 'Record update timestamp';

-- Rollback
-- DROP TABLE IF EXISTS %s;
`, tableName, time.Now().Format(time.RFC3339), tableName, idFieldName, tableName, tableName, tableName, tableName, tableName, tableName, tableName, tableName, tableName, idFieldName, tableName, tableName, tableName)
}

func (g *Generator) generateAlterMigrationContent(tableName string) string {
	template := `-- Migration: Alter table {{TABLE_NAME}}
-- Generated at: {{TIMESTAMP}}


-- Add columns:
-- ALTER TABLE {{TABLE_NAME}} ADD COLUMN IF NOT EXISTS column_name TEXT;
-- ALTER TABLE {{TABLE_NAME}} ADD COLUMN IF NOT EXISTS column_name TEXT NOT NULL DEFAULT '';
-- ALTER TABLE {{TABLE_NAME}} ADD COLUMN IF NOT EXISTS column_name INTEGER NOT NULL DEFAULT 0;
-- ALTER TABLE {{TABLE_NAME}} ADD COLUMN IF NOT EXISTS column_name BOOLEAN NOT NULL DEFAULT false;
-- ALTER TABLE {{TABLE_NAME}} ADD COLUMN IF NOT EXISTS column_name JSONB DEFAULT '{}'::jsonb;
-- ALTER TABLE {{TABLE_NAME}} ADD COLUMN IF NOT EXISTS column_name TIMESTAMPTZ;
-- ALTER TABLE {{TABLE_NAME}} ADD COLUMN IF NOT EXISTS column_name UUID REFERENCES other_table(other_id);
--
-- Drop columns:
-- ALTER TABLE {{TABLE_NAME}} DROP COLUMN IF EXISTS column_name;
--
-- Rename columns:
-- ALTER TABLE {{TABLE_NAME}} RENAME COLUMN old_column_name TO new_column_name;
--
-- Modify column types:
-- ALTER TABLE {{TABLE_NAME}} ALTER COLUMN column_name TYPE TEXT;
-- ALTER TABLE {{TABLE_NAME}} ALTER COLUMN column_name TYPE INTEGER USING column_name::integer;
--
-- Set column defaults:
-- ALTER TABLE {{TABLE_NAME}} ALTER COLUMN column_name SET DEFAULT 'default_value';
-- ALTER TABLE {{TABLE_NAME}} ALTER COLUMN column_name SET DEFAULT NOW();
--
-- Drop column defaults:
-- ALTER TABLE {{TABLE_NAME}} ALTER COLUMN column_name DROP DEFAULT;
--
-- Set NOT NULL constraint:
-- ALTER TABLE {{TABLE_NAME}} ALTER COLUMN column_name SET NOT NULL;
--
-- Drop NOT NULL constraint:
-- ALTER TABLE {{TABLE_NAME}} ALTER COLUMN column_name DROP NOT NULL;
--
-- Add unique constraint:
-- ALTER TABLE {{TABLE_NAME}} ADD CONSTRAINT unique_column_name UNIQUE (column_name);
--
-- Drop unique constraint:
-- ALTER TABLE {{TABLE_NAME}} DROP CONSTRAINT IF EXISTS unique_column_name;
--
-- Add check constraint:
-- ALTER TABLE {{TABLE_NAME}} ADD CONSTRAINT check_column_name CHECK (column_name > 0);
--
-- Drop check constraint:
-- ALTER TABLE {{TABLE_NAME}} DROP CONSTRAINT IF EXISTS check_column_name;
--
-- Add foreign key:
-- ALTER TABLE {{TABLE_NAME}} ADD CONSTRAINT fk_column_name FOREIGN KEY (column_name) REFERENCES other_table(other_id) ON DELETE CASCADE;
--
-- Drop foreign key:
-- ALTER TABLE {{TABLE_NAME}} DROP CONSTRAINT IF EXISTS fk_column_name;
--
-- Create indexes:
-- CREATE INDEX IF NOT EXISTS idx_{{TABLE_NAME}}_column_name ON {{TABLE_NAME}}(column_name);
-- CREATE INDEX IF NOT EXISTS idx_{{TABLE_NAME}}_column1_column2 ON {{TABLE_NAME}}(column1, column2);
-- CREATE UNIQUE INDEX IF NOT EXISTS idx_{{TABLE_NAME}}_unique_column ON {{TABLE_NAME}}(column_name);
--
-- Drop indexes:
-- DROP INDEX IF EXISTS idx_{{TABLE_NAME}}_column_name;
--
-- Add column comments:
-- COMMENT ON COLUMN {{TABLE_NAME}}.column_name IS 'Description of the column';
--
-- Update table comment:
-- COMMENT ON TABLE {{TABLE_NAME}} IS 'Updated table description';

-- Rollback
-- Customize rollback statements below (reverse the changes above):
-- ALTER TABLE {{TABLE_NAME}} DROP COLUMN IF EXISTS column_name;
-- ALTER TABLE {{TABLE_NAME}} RENAME COLUMN new_column_name TO old_column_name;
-- ALTER TABLE {{TABLE_NAME}} DROP CONSTRAINT IF EXISTS constraint_name;
-- DROP INDEX IF EXISTS idx_{{TABLE_NAME}}_column_name;
`

	result := strings.ReplaceAll(template, "{{TABLE_NAME}}", tableName)
	result = strings.ReplaceAll(result, "{{TIMESTAMP}}", time.Now().Format(time.RFC3339))
	return result
}

func (g *Generator) GenerateSQLFromModel(tableName string) (string, error) {
	modelPath := filepath.Join(g.modelsDir, fmt.Sprintf("%s.go", strings.ToLower(tableName)))

	src, err := os.ReadFile(modelPath)
	if err != nil {
		return "", fmt.Errorf("failed to read model file: %w", err)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, modelPath, src, parser.ParseComments)
	if err != nil {
		return "", fmt.Errorf("failed to parse model file: %w", err)
	}

	var columns []string
	tableNameLower := strings.ToLower(tableName)
	idFieldName := fmt.Sprintf("%s_id", tableNameLower)

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.StructType:
			if x.Fields != nil {
				for _, field := range x.Fields.List {
					if field.Tag != nil {
						tag := field.Tag.Value
						colName, constraintsStr := extractDBTag(tag)
						if colName == "" {
							continue
						}

						fieldName := field.Names[0].Name
						if fieldName == g.toPascalCase(idFieldName) {
							continue
						}
						if fieldName == "CreatedTimestamp" || fieldName == "UpdatedTimestamp" {
							continue
						}

						goType := getGoType(field.Type)
						sqlType := goTypeToSQL(goType)
						constraints := parseConstraints(constraintsStr)

						colDef := fmt.Sprintf("\t%s %s%s", colName, sqlType, constraints)
						columns = append(columns, colDef)
					}
				}
			}
		}
		return true
	})

	return strings.Join(columns, ",\n"), nil
}

func extractDBTag(tag string) (string, string) {
	tag = strings.Trim(tag, "`")
	parts := strings.Split(tag, " ")
	for _, part := range parts {
		if strings.HasPrefix(part, "db:") {
			value := strings.TrimPrefix(part, "db:")
			value = strings.Trim(value, "\"")
			if idx := strings.Index(value, ","); idx != -1 {
				colName := value[:idx]
				constraints := value[idx+1:]
				return colName, constraints
			}
			return value, ""
		}
	}
	return "", ""
}

func parseConstraints(constraintsStr string) string {
	var constraints []string

	if strings.Contains(constraintsStr, "notnull") {
		constraints = append(constraints, " NOT NULL")
	}
	if strings.Contains(constraintsStr, "unique") {
		constraints = append(constraints, " UNIQUE")
	}

	return strings.Join(constraints, "")
}

func getGoType(expr ast.Expr) string {
	switch x := expr.(type) {
	case *ast.Ident:
		return x.Name
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", x.X, x.Sel)
	case *ast.MapType:
		return "map"
	case *ast.ArrayType:
		return "array"
	}
	return "unknown"
}

func goTypeToSQL(goType string) string {
	switch goType {
	case "string":
		return "TEXT"
	case "int", "int8", "int16", "int32", "int64":
		return "INTEGER"
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return "INTEGER"
	case "float32", "float64":
		return "NUMERIC(10, 2)"
	case "bool":
		return "BOOLEAN"
	case "time.Time":
		return "TIMESTAMPTZ"
	case "map":
		return "JSONB"
	default:
		return "TEXT"
	}
}

func (g *Generator) generateCreateMigrationFromModel(tableName, idFieldName, sqlColumns string) string {
	return fmt.Sprintf(`-- Migration: Create table %s
-- Generated at: %s
-- Generated from model: internal/models/%s.go

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS %s (
	%s UUID PRIMARY KEY DEFAULT gen_random_uuid(),
%s,
	created_timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add table and column comments
COMMENT ON TABLE %s IS 'Table for %s';
COMMENT ON COLUMN %s.%s IS 'Primary key UUID';
COMMENT ON COLUMN %s.created_timestamp IS 'Record creation timestamp';
COMMENT ON COLUMN %s.updated_timestamp IS 'Record update timestamp';

-- Rollback
-- DROP TABLE IF EXISTS %s;
`, tableName, time.Now().Format(time.RFC3339), tableName, tableName, idFieldName, sqlColumns, tableName, tableName, tableName, idFieldName, tableName, tableName, tableName)
}

func (g *Generator) generateAlterMigrationFromModel(tableName, sqlColumns string) string {
	lines := strings.Split(sqlColumns, "\n")
	var alterStatements []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		line = strings.TrimSuffix(line, ",")
		alterStatements = append(alterStatements, fmt.Sprintf("\tALTER TABLE %s ADD COLUMN IF NOT EXISTS %s;", tableName, line))
	}

	return fmt.Sprintf(`-- Migration: Alter table %s
-- Generated at: %s
-- Generated from model: internal/models/%s.go

%s

-- Rollback
-- Customize rollback statements below (reverse the changes above):
-- ALTER TABLE %s DROP COLUMN IF EXISTS column_name;
`, tableName, time.Now().Format(time.RFC3339), tableName, strings.Join(alterStatements, "\n"), tableName)
}
