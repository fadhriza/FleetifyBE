package models

import (
	"fleetify/internal/migration"
	"time"
)

type Items struct {
	Id        int64     `db:"id" json:"id"`
	Name      string    `db:"name,notnull" json:"name"`
	Stock     int       `db:"stock" json:"stock"`
	Price     float64   `db:"price,notnull" json:"price"`
	Category  string    `db:"category" json:"category"`
	Unit      string    `db:"unit" json:"unit"`
	MinStock  int       `db:"min_stock" json:"min_stock"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (Items) TableName() string {
	return "items"
}

func (Items) GetID() string {
	return "id"
}

func init() {
	migration.RegisterSeeder("Items", func() interface{} {
		return SeedItems()
	})
}

func SeedItems() []Items {
	now := time.Now()
	return []Items{
		{Name: "Engine Oil 5W-30", Stock: 50, Price: 150000, Category: "oil", Unit: "liter", MinStock: 10, CreatedAt: now, UpdatedAt: now},
		{Name: "Engine Oil 10W-40", Stock: 45, Price: 140000, Category: "oil", Unit: "liter", MinStock: 10, CreatedAt: now, UpdatedAt: now},
		{Name: "Brake Pad Front", Stock: 30, Price: 250000, Category: "parts", Unit: "set", MinStock: 5, CreatedAt: now, UpdatedAt: now},
		{Name: "Brake Pad Rear", Stock: 25, Price: 200000, Category: "parts", Unit: "set", MinStock: 5, CreatedAt: now, UpdatedAt: now},
		{Name: "Air Filter", Stock: 40, Price: 75000, Category: "parts", Unit: "pcs", MinStock: 10, CreatedAt: now, UpdatedAt: now},
		{Name: "Fuel Filter", Stock: 35, Price: 85000, Category: "parts", Unit: "pcs", MinStock: 10, CreatedAt: now, UpdatedAt: now},
		{Name: "Tire 205/55R16", Stock: 20, Price: 800000, Category: "tire", Unit: "pcs", MinStock: 4, CreatedAt: now, UpdatedAt: now},
		{Name: "Tire 215/60R16", Stock: 18, Price: 850000, Category: "tire", Unit: "pcs", MinStock: 4, CreatedAt: now, UpdatedAt: now},
		{Name: "Battery 12V 60Ah", Stock: 15, Price: 1200000, Category: "battery", Unit: "pcs", MinStock: 3, CreatedAt: now, UpdatedAt: now},
		{Name: "Battery 12V 70Ah", Stock: 12, Price: 1400000, Category: "battery", Unit: "pcs", MinStock: 3, CreatedAt: now, UpdatedAt: now},
		{Name: "Spark Plug", Stock: 60, Price: 45000, Category: "parts", Unit: "pcs", MinStock: 20, CreatedAt: now, UpdatedAt: now},
		{Name: "Radiator Coolant", Stock: 30, Price: 95000, Category: "oil", Unit: "liter", MinStock: 10, CreatedAt: now, UpdatedAt: now},
		{Name: "Windshield Wiper", Stock: 25, Price: 55000, Category: "parts", Unit: "set", MinStock: 5, CreatedAt: now, UpdatedAt: now},
		{Name: "Headlight Bulb H4", Stock: 20, Price: 125000, Category: "parts", Unit: "pcs", MinStock: 5, CreatedAt: now, UpdatedAt: now},
		{Name: "Brake Fluid", Stock: 35, Price: 65000, Category: "oil", Unit: "liter", MinStock: 10, CreatedAt: now, UpdatedAt: now},
	}
}
