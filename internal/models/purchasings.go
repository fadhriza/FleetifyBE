package models

import (
	"fleetify/internal/migration"
	"time"
)

type Purchasings struct {
	PurchasingsId string    `db:"purchasings_id" json:"purchasings_id"`
	Date          time.Time `db:"date,notnull" json:"date"`
	SupplierId    string    `db:"supplier_id,notnull" json:"supplier_id"`
	UserId        string    `db:"user_id,notnull" json:"user_id"`
	GrandTotal    float64   `db:"grand_total,notnull" json:"grand_total"`
	Status        string    `db:"status" json:"status"`
	Notes         string    `db:"notes" json:"notes"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

func (Purchasings) TableName() string {
	return "purchasings"
}

func (Purchasings) GetID() string {
	return "purchasings_id"
}

func init() {
	migration.RegisterSeeder("Purchasings", func() interface{} {
		return SeedPurchasings()
	})
}

func SeedPurchasings() []Purchasings {
	return []Purchasings{
		// Seeder will be updated after UUID migration
	}
}
