package models

import (
	"fleetify/internal/migration"
	"time"
)

type Purchasings struct {
	Id         int64     `db:"id" json:"id"`
	Date       time.Time `db:"date,notnull" json:"date"`
	SupplierId int64     `db:"supplier_id,notnull" json:"supplier_id"`
	UserId     int64     `db:"user_id,notnull" json:"user_id"`      
	GrandTotal float64   `db:"grand_total,notnull" json:"grand_total"`
	Status     string    `db:"status" json:"status"`
	Notes      string    `db:"notes" json:"notes"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

func (Purchasings) TableName() string {
	return "purchasings"
}

func (Purchasings) GetID() string {
	return "id"
}

func init() {
	migration.RegisterSeeder("Purchasings", func() interface{} {
		return SeedPurchasings()
	})
}

func SeedPurchasings() []Purchasings {
	now := time.Now()
	return []Purchasings{
		{Date: now.AddDate(0, 0, -30), SupplierId: 1, UserId: 2, GrandTotal: 7500000, Status: "completed", Notes: "Monthly stock replenishment", CreatedAt: now},
		{Date: now.AddDate(0, 0, -25), SupplierId: 2, UserId: 2, GrandTotal: 3200000, Status: "completed", Notes: "Tire replacement order", CreatedAt: now},
		{Date: now.AddDate(0, 0, -20), SupplierId: 3, UserId: 3, GrandTotal: 4200000, Status: "completed", Notes: "Engine oil bulk order", CreatedAt: now},
		{Date: now.AddDate(0, 0, -15), SupplierId: 4, UserId: 2, GrandTotal: 2400000, Status: "approved", Notes: "Battery stock order", CreatedAt: now},
		{Date: now.AddDate(0, 0, -10), SupplierId: 1, UserId: 3, GrandTotal: 1800000, Status: "pending", Notes: "Spare parts order", CreatedAt: now},
		{Date: now.AddDate(0, 0, -5), SupplierId: 5, UserId: 2, GrandTotal: 950000, Status: "pending", Notes: "Service equipment", CreatedAt: now},
		{Date: now.AddDate(0, 0, -2), SupplierId: 6, UserId: 3, GrandTotal: 1350000, Status: "approved", Notes: "Motor parts order", CreatedAt: now},
		{Date: now, SupplierId: 7, UserId: 2, GrandTotal: 1600000, Status: "pending", Notes: "Tire stock order", CreatedAt: now},
	}
}
