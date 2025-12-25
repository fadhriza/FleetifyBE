package models

import (
	"fleetify/internal/migration"
	"time"
)

type Roles struct {

	RoleOID string `db:"role_oid,notnull" json:"role_oid"`
	RoleName string `db:"role_name,notnull" json:"role_name"`
	RoleDescription string `db:"role_description" json:"role_description"`


	// Timestamps
	CreatedTimestamp time.Time `db:"created_timestamp" json:"created_timestamp"`
	UpdatedTimestamp time.Time `db:"updated_timestamp" json:"updated_timestamp"`
}

// TableName returns the table name
func (Roles) TableName() string {
	return "roles"
}

// GetID returns the primary key field name
func (Roles) GetID() string {
	return "roles_id"
}


func init() {
	migration.RegisterSeeder("Roles", func() interface{} {
		return SeedRoles()
	})
}

func SeedRoles() []Roles {
	return []Roles{
		{
			RoleOID:          "ADMIN",
			RoleName:         "Admin",
			RoleDescription:  "Administrator role with full permissions.",
			CreatedTimestamp: time.Now(),
			UpdatedTimestamp: time.Now(),
		},
		{
			RoleOID:          "MANAGER",
			RoleName:         "Manager",
			RoleDescription:  "Manager role with extended permissions.",
			CreatedTimestamp: time.Now(),
			UpdatedTimestamp: time.Now(),
		},
		{
			RoleOID:          "SUPPLIERS",
			RoleName:         "Suppliers",
			RoleDescription:  "Supplier users.",
			CreatedTimestamp: time.Now(),
			UpdatedTimestamp: time.Now(),
		},
		{
			RoleOID:          "MITRA",
			RoleName:         "Mitra",
			RoleDescription:  "Mitra users.",
			CreatedTimestamp: time.Now(),
			UpdatedTimestamp: time.Now(),
		},
	}
}
