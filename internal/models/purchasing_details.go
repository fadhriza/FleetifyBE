package models

import (
	"fleetify/internal/migration"
)

type PurchasingDetails struct {
	Id           int64   `db:"id" json:"id"`
	PurchasingId int64   `db:"purchasing_id,notnull" json:"purchasing_id"`
	ItemId       int64   `db:"item_id,notnull" json:"item_id"`
	Qty          int     `db:"qty,notnull" json:"qty"`
	Subtotal     float64 `db:"subtotal,notnull" json:"subtotal"`
}

func (PurchasingDetails) TableName() string {
	return "purchasing_details"
}

func (PurchasingDetails) GetID() string {
	return "id"
}

func init() {
	migration.RegisterSeeder("PurchasingDetails", func() interface{} {
		return SeedPurchasingDetails()
	})
}

func SeedPurchasingDetails() []PurchasingDetails {
	return []PurchasingDetails{
		{PurchasingId: 1, ItemId: 1, Qty: 30, Subtotal: 4500000},
		{PurchasingId: 1, ItemId: 2, Qty: 20, Subtotal: 2800000},
		{PurchasingId: 1, ItemId: 5, Qty: 10, Subtotal: 200000},
		{PurchasingId: 2, ItemId: 7, Qty: 4, Subtotal: 3200000},
		{PurchasingId: 3, ItemId: 1, Qty: 20, Subtotal: 3000000},
		{PurchasingId: 3, ItemId: 12, Qty: 10, Subtotal: 950000},
		{PurchasingId: 3, ItemId: 15, Qty: 5, Subtotal: 325000},
		{PurchasingId: 4, ItemId: 9, Qty: 2, Subtotal: 2400000},
		{PurchasingId: 5, ItemId: 3, Qty: 4, Subtotal: 1000000},
		{PurchasingId: 5, ItemId: 4, Qty: 4, Subtotal: 800000},
		{PurchasingId: 6, ItemId: 11, Qty: 10, Subtotal: 450000},
		{PurchasingId: 6, ItemId: 13, Qty: 5, Subtotal: 275000},
		{PurchasingId: 6, ItemId: 14, Qty: 2, Subtotal: 250000},
		{PurchasingId: 7, ItemId: 6, Qty: 10, Subtotal: 850000},
		{PurchasingId: 7, ItemId: 10, Qty: 1, Subtotal: 500000},
		{PurchasingId: 8, ItemId: 8, Qty: 2, Subtotal: 1700000},
	}
}
