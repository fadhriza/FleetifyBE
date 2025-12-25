package migration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"fleetify/internal/database"
	"fleetify/pkg/errors"

	"github.com/jackc/pgx/v5/pgconn"
)

type Runner struct {
	migrationsDir string
}

func NewRunner() *Runner {
	return &Runner{
		migrationsDir: "migrations",
	}
}

func (r *Runner) RunMigrations() error {
	if err := database.Connect(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer database.Close()

	if err := r.ensureMigrationTable(); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}

	migrations, err := r.getPendingMigrations()
	if err != nil {
		return fmt.Errorf("failed to get pending migrations: %w", err)
	}

	if len(migrations) == 0 {
		fmt.Println("SUCCESS: No pending migrations")
		return nil
	}

	fmt.Printf("NOTE: Found %d pending migration(s)\n", len(migrations))

	for _, migration := range migrations {
		if err := r.runMigration(migration); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", migration, err)
		}
	}

	return nil
}

func (r *Runner) RollbackMigration(migrationName string) error {
	if err := database.Connect(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer database.Close()

	migrationPath := filepath.Join(r.migrationsDir, migrationName)
	if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
		return fmt.Errorf("migration file not found: %s", migrationName)
	}

	content, err := os.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	rollbackSQL := r.extractRollbackSQL(string(content))
	if rollbackSQL == "" {
		return fmt.Errorf("no rollback SQL found in migration file")
	}

	ctx := context.Background()
	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, rollbackSQL); err != nil {
		errors.LogError("Migration Rollback Error", err)
		return fmt.Errorf("failed to execute rollback: %w", err)
	}

	if _, err := tx.Exec(ctx, "DELETE FROM schema_migrations WHERE name = $1", migrationName); err != nil {
		errors.LogError("Migration Record Delete Error", err)
		return fmt.Errorf("failed to delete migration record: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	fmt.Printf("SUCCESS: Rolled back migration: %s\n", migrationName)
	return nil
}

func (r *Runner) ensureMigrationTable() error {
	ctx := context.Background()
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			executed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`
	_, err := database.DB.Exec(ctx, query)
	return err
}

func (r *Runner) getPendingMigrations() ([]string, error) {
	ctx := context.Background()

	var executedMigrations []string
	rows, err := database.DB.Query(ctx, "SELECT name FROM schema_migrations ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		executedMigrations = append(executedMigrations, name)
	}

	executedMap := make(map[string]bool)
	for _, name := range executedMigrations {
		executedMap[name] = true
	}

	files, err := os.ReadDir(r.migrationsDir)
	if err != nil {
		return nil, err
	}

	var pending []string
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}
		if !executedMap[file.Name()] {
			pending = append(pending, file.Name())
		}
	}

	sort.Strings(pending)
	return pending, nil
}

func (r *Runner) runMigration(migrationName string) error {
	migrationPath := filepath.Join(r.migrationsDir, migrationName)
	content, err := os.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	ctx := context.Background()
	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	sql := string(content)
	statements := r.parseSQLStatements(sql)

	if len(statements) > 0 {
		fmt.Println("Executing SQL:")
		for i, stmt := range statements {
			r.displaySQLStatement(i+1, stmt)
			result, err := tx.Exec(ctx, stmt)
			if err != nil {
				errors.LogError("Migration Execution Error", err)
				return fmt.Errorf("failed to execute migration: %w", err)
			}
			r.displayPostgreSQLResponse(result)
		}
		fmt.Println()
	}

	if _, err := tx.Exec(ctx, "INSERT INTO schema_migrations (name) VALUES ($1)", migrationName); err != nil {
		errors.LogError("Migration Record Insert Error", err)
		return fmt.Errorf("failed to record migration: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	fmt.Printf("SUCCESS: Executed migration: %s\n", migrationName)
	return nil
}

func (r *Runner) parseSQLStatements(sql string) []string {
	lines := strings.Split(sql, "\n")
	var statements []string
	var currentStatement strings.Builder

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "" || strings.HasPrefix(trimmed, "--") {
			continue
		}

		currentStatement.WriteString(line)
		currentStatement.WriteString("\n")

		if strings.HasSuffix(trimmed, ";") {
			stmt := strings.TrimSpace(currentStatement.String())
			if stmt != "" {
				statements = append(statements, stmt)
			}
			currentStatement.Reset()
		}
	}

	if currentStatement.Len() > 0 {
		stmt := strings.TrimSpace(currentStatement.String())
		if stmt != "" {
			statements = append(statements, stmt)
		}
	}

	return statements
}

func (r *Runner) displaySQLStatement(num int, stmt string) {
	stmtLines := strings.Split(stmt, "\n")
	for j, stmtLine := range stmtLines {
		trimmed := strings.TrimSpace(stmtLine)
		if trimmed == "" {
			continue
		}
		if j == 0 {
			fmt.Printf("  [%d] %s\n", num, trimmed)
		} else {
			fmt.Printf("      %s\n", trimmed)
		}
	}
}

func (r *Runner) displayPostgreSQLResponse(result pgconn.CommandTag) {
	response := result.String()
	if response != "" {
		fmt.Printf("      → PostgreSQL: %s\n", response)
	}
}

func (r *Runner) extractRollbackSQL(content string) string {
	lines := strings.Split(content, "\n")
	var rollbackLines []string
	inRollback := false

	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "-- rollback") {
			inRollback = true
			continue
		}
		if inRollback {
			if strings.HasPrefix(strings.TrimSpace(line), "--") && strings.Contains(strings.ToLower(line), "end rollback") {
				break
			}
			rollbackLines = append(rollbackLines, line)
		}
	}

	return strings.Join(rollbackLines, "\n")
}
