package models

import (
	"fleetify/internal/migration"
)

type PurchasingDetails struct {
	PurchasingDetailsId string  `db:"purchasing_details_id" json:"purchasing_details_id"`
	PurchasingId        string  `db:"purchasing_id,notnull" json:"purchasing_id"`
	ItemId              string  `db:"item_id,notnull" json:"item_id"`
	Qty                 int     `db:"qty,notnull" json:"qty"`
	Subtotal            float64 `db:"subtotal,notnull" json:"subtotal"`
}

func (PurchasingDetails) TableName() string {
	return "purchasing_details"
}

func (PurchasingDetails) GetID() string {
	return "purchasing_details_id"
}

func init() {
	migration.RegisterSeeder("PurchasingDetails", func() interface{} {
		return SeedPurchasingDetails()
	})
}

func SeedPurchasingDetails() []PurchasingDetails {
	return []PurchasingDetails{
		// Seeder will be updated after UUID migration
		// {PurchasingId: "uuid", ItemId: "uuid", Qty: 30, Subtotal: 4500000},
	}
}
