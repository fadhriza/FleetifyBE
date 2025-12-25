package migration

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"fleetify/internal/database"
	"fleetify/pkg/errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Seeder struct {
	modelsDir string
}

func NewSeeder() *Seeder {
	return &Seeder{
		modelsDir: "internal/models",
	}
}

func (s *Seeder) RunSeeder(tableName string) error {
	if err := database.Connect(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer database.Close()

	modelName := toPascalCase(tableName)
	seedFuncName := fmt.Sprintf("Seed%s", modelName)
	
	seedData, err := s.getSeedData(modelName, seedFuncName)
	if err != nil {
		return fmt.Errorf("failed to get seed data: %w", err)
	}

	if len(seedData) == 0 {
		fmt.Printf("NOTE: No seed data found in %s\n", seedFuncName)
		return nil
	}

	ctx := context.Background()
	tableNameLower := strings.ToLower(tableName)
	idFieldName := fmt.Sprintf("%s_id", tableNameLower)

	fmt.Printf("Seeding %s table...\n", tableNameLower)
	fmt.Printf("Found %d record(s) to seed\n\n", len(seedData))

	for i, record := range seedData {
		recordMap := s.structToMap(record)
		
		if idField := recordMap[idFieldName]; idField == "" || idField == nil {
			recordMap[idFieldName] = uuid.New().String()
		}

		if _, ok := recordMap["created_timestamp"]; !ok || recordMap["created_timestamp"] == nil {
			recordMap["created_timestamp"] = time.Now()
		}
		if _, ok := recordMap["updated_timestamp"]; !ok || recordMap["updated_timestamp"] == nil {
			recordMap["updated_timestamp"] = time.Now()
		}

		s.hashPasswordFields(recordMap)

		columns := []string{}
		values := []interface{}{}
		placeholders := []string{}

		idx := 1
		for col, val := range recordMap {
			columns = append(columns, col)
			values = append(values, val)
			placeholders = append(placeholders, fmt.Sprintf("$%d", idx))
			idx++
		}

		query := fmt.Sprintf(
			"INSERT INTO %s (%s) VALUES (%s) ON CONFLICT DO NOTHING",
			tableNameLower,
			strings.Join(columns, ", "),
			strings.Join(placeholders, ", "),
		)

		fmt.Printf("  [%d] Inserting record...\n", i+1)
		result, err := database.DB.Exec(ctx, query, values...)
		if err != nil {
			errors.LogError("Seeder Insert Error", err)
			return fmt.Errorf("failed to insert seed data: %w", err)
		}

		response := result.String()
		if response != "" {
			fmt.Printf("      → PostgreSQL: %s\n", response)
		}
	}

	fmt.Println()
	fmt.Printf("SUCCESS: Seeded %d record(s) into %s\n", len(seedData), tableNameLower)
	return nil
}

func (s *Seeder) getSeedData(modelName, seedFuncName string) ([]interface{}, error) {
	seedResult := getSeedDataFromRegistry(modelName)
	if seedResult == nil {
		return nil, fmt.Errorf("seed function for %s not found. Add case in seeder_registry.go", modelName)
	}

	seedSlice := reflect.ValueOf(seedResult)
	if seedSlice.Kind() != reflect.Slice {
		return nil, fmt.Errorf("seed function must return a slice")
	}

	var seedData []interface{}
	for i := 0; i < seedSlice.Len(); i++ {
		seedData = append(seedData, seedSlice.Index(i).Interface())
	}

	return seedData, nil
}

func (s *Seeder) structToMap(v interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		dbTag := field.Tag.Get("db")
		if dbTag == "" || dbTag == "-" {
			continue
		}

		colName := strings.Split(dbTag, ",")[0]
		if colName == "" {
			continue
		}

		if fieldValue.CanInterface() {
			result[colName] = fieldValue.Interface()
		}
	}

	return result
}

func (s *Seeder) hashPasswordFields(recordMap map[string]interface{}) {
	for key, value := range recordMap {
		if strings.Contains(strings.ToLower(key), "password") {
			if strValue, ok := value.(string); ok && strValue != "" {
				hashed, err := bcrypt.GenerateFromPassword([]byte(strValue), bcrypt.DefaultCost)
				if err == nil {
					recordMap[key] = string(hashed)
				}
			}
		}
	}
}

func toPascalCase(s string) string {
	parts := strings.Split(s, "_")
	var result strings.Builder
	for _, part := range parts {
		if len(part) > 0 {
			result.WriteString(strings.ToUpper(part[:1]) + strings.ToLower(part[1:]))
		}
	}
	return result.String()
}

