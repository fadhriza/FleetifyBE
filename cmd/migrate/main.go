package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"fleetify/internal/config"
	"fleetify/internal/migration"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "createtable":
		if len(os.Args) < 3 {
			log.Fatal("ERROR: Table name is required. Usage: go run cmd/migrate/main.go createtable <table_name>")
		}
		tableName := os.Args[2]
		generateTable(tableName)
	case "altertable":
		if len(os.Args) < 3 {
			log.Fatal("ERROR: Table name is required. Usage: go run cmd/migrate/main.go altertable <table_name>")
		}
		tableName := os.Args[2]
		generateAlterTable(tableName)
	case "migrate":
		runMigrations()
	case "rollback":
		if len(os.Args) < 3 {
			log.Fatal("ERROR: Migration name is required. Usage: go run cmd/migrate/main.go rollback <migration_name>")
		}
		migrationName := os.Args[2]
		rollbackMigration(migrationName)
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Fleetify Migration Tool")
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/migrate/main.go createtable <table_name>  - Generate a new table migration")
	fmt.Println("  go run cmd/migrate/main.go altertable <table_name>   - Generate an alter table migration")
	fmt.Println("  go run cmd/migrate/main.go migrate                   - Run pending migrations")
	fmt.Println("  go run cmd/migrate/main.go rollback <migration_name> - Rollback a migration")
}

func generateTable(tableName string) {
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	generator := migration.NewGenerator()
	if err := generator.GenerateTable(tableName); err != nil {
		log.Fatalf("Failed to generate table: %v", err)
	}

	fmt.Printf("SUCCESS: Successfully generated table migration for: %s\n", tableName)
	fmt.Printf("NOTE: Model file: internal/models/%s.go\n", strings.ToLower(tableName))
	fmt.Printf("NOTE: Migration file: migrations/%s_*.sql\n", strings.ToLower(tableName))
}

func generateAlterTable(tableName string) {
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	generator := migration.NewGenerator()
	if err := generator.GenerateAlterTable(tableName); err != nil {
		log.Fatalf("Failed to generate alter table migration: %v", err)
	}

	fmt.Printf("SUCCESS: Successfully generated alter table migration for: %s\n", tableName)
	fmt.Printf("NOTE: Migration file: migrations/*_%s_alter.sql\n", strings.ToLower(tableName))
}

func runMigrations() {
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	runner := migration.NewRunner()
	if err := runner.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	fmt.Println("SUCCESS: All migrations completed successfully")
}

func rollbackMigration(migrationName string) {
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	runner := migration.NewRunner()
	if err := runner.RollbackMigration(migrationName); err != nil {
		log.Fatalf("Failed to rollback migration: %v", err)
	}

	fmt.Printf("SUCCESS: Successfully rolled back migration: %s\n", migrationName)
}
