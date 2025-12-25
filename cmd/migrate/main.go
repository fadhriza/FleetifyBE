package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"fleetify/internal/config"
	"fleetify/internal/migration"
	_ "fleetify/internal/models"
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
		createTableModel(tableName)
	case "generatesql":
		if len(os.Args) < 3 {
			log.Fatal("ERROR: Table name is required. Usage: go run cmd/migrate/main.go generatesql <table_name>")
		}
		tableName := os.Args[2]
		generateSQLFromModel(tableName)
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
	case "seed":
		if len(os.Args) < 3 {
			log.Fatal("ERROR: Table name is required. Usage: go run cmd/migrate/main.go seed <table_name>")
		}
		tableName := os.Args[2]
		runSeeder(tableName)
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Fleetify Migration Tool")
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/migrate/main.go createtable <table_name>  - Create model template")
	fmt.Println("  go run cmd/migrate/main.go generatesql <table_name>   - Generate SQL migration from model")
	fmt.Println("  go run cmd/migrate/main.go altertable <table_name>    - Generate ALTER migration from model")
	fmt.Println("  go run cmd/migrate/main.go migrate                     - Run pending migrations")
	fmt.Println("  go run cmd/migrate/main.go rollback <migration_name> - Rollback a migration")
	fmt.Println("  go run cmd/migrate/main.go seed <table_name>           - Run seeders")
}

func createTableModel(tableName string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you need dbseeds? [y/n]: ")
	answer, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}
	answer = strings.TrimSpace(strings.ToLower(answer))
	needsSeeder := answer == "y" || answer == "yes"

	generator := migration.NewGenerator()
	if err := generator.CreateTableModel(tableName, needsSeeder); err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	fmt.Printf("SUCCESS: Created model template: internal/models/%s.go\n", strings.ToLower(tableName))
	if needsSeeder {
		fmt.Printf("NOTE: Seeder template included in model\n")
	}
	fmt.Printf("NOTE: Edit the model, then run: go run cmd/migrate/main.go generatesql %s\n", tableName)
}

func runSeeder(tableName string) {
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	runner := migration.NewSeeder()
	if err := runner.RunSeeder(tableName); err != nil {
		log.Fatalf("Failed to run seeder: %v", err)
	}

	fmt.Printf("SUCCESS: Seeder executed for: %s\n", tableName)
}

func generateSQLFromModel(tableName string) {
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	generator := migration.NewGenerator()
	if err := generator.GenerateSQLMigration(tableName); err != nil {
		log.Fatalf("Failed to generate SQL migration: %v", err)
	}

	fmt.Printf("SUCCESS: Generated SQL migration from model: migrations/*_%s_create.sql\n", strings.ToLower(tableName))
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
