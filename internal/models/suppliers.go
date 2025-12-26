package models

import (
	"fleetify/internal/migration"
	"time"
)

type Suppliers struct {
	SuppliersId  string    `db:"suppliers_id" json:"suppliers_id"`
	Name         string    `db:"name,notnull" json:"name"`
	Email        string    `db:"email" json:"email"`
	Address      string    `db:"address" json:"address"`
	Phone        string    `db:"phone" json:"phone"`
	SupplierType string    `db:"supplier_type" json:"supplier_type"`
	IsActive     bool      `db:"is_active" json:"is_active"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

func (Suppliers) TableName() string {
	return "suppliers"
}

func (Suppliers) GetID() string {
	return "suppliers_id"
}

func init() {
	migration.RegisterSeeder("Suppliers", func() interface{} {
		return SeedSuppliers()
	})
}

func SeedSuppliers() []Suppliers {
	now := time.Now()
	return []Suppliers{
		{Name: "PT Auto Parts Indonesia", Email: "contact@autoparts.id", Address: "Jl. Sudirman No. 123, Jakarta", Phone: "021-12345678", SupplierType: "parts", IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "CV Tire Center", Email: "info@tirecenter.co.id", Address: "Jl. Gatot Subroto No. 456, Jakarta", Phone: "021-23456789", SupplierType: "tire", IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "PT Oil Distributor", Email: "sales@oildist.com", Address: "Jl. Thamrin No. 789, Jakarta", Phone: "021-34567890", SupplierType: "parts", IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Battery Pro Indonesia", Email: "order@batterypro.id", Address: "Jl. HR Rasuna Said No. 321, Jakarta", Phone: "021-45678901", SupplierType: "parts", IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "CV Service Equipment", Email: "info@serviceequip.co.id", Address: "Jl. Kebon Jeruk No. 654, Jakarta", Phone: "021-56789012", SupplierType: "service", IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "PT Motor Parts", Email: "sales@motorparts.id", Address: "Jl. Cikini Raya No. 987, Jakarta", Phone: "021-67890123", SupplierType: "parts", IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Tire Master Jakarta", Email: "contact@tiremaster.id", Address: "Jl. Kemang Raya No. 147, Jakarta", Phone: "021-78901234", SupplierType: "tire", IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "PT Lubricant Supply", Email: "order@lubricant.id", Address: "Jl. Senopati No. 258, Jakarta", Phone: "021-89012345", SupplierType: "parts", IsActive: true, CreatedAt: now, UpdatedAt: now},
	}
}
