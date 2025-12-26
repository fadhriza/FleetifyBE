package models

import (
	"fleetify/internal/migration"
	"time"
)

type Users struct {
	UsersId   string    `db:"users_id" json:"users_id"`
	Id        int64     `db:"id" json:"id"`
	Username  string    `db:"username,unique,notnull" json:"username"`
	Password  string    `db:"password,notnull" json:"password"`
	Role      string    `db:"role,notnull" json:"role"`
	FullName  string    `db:"full_name,notnull" json:"full_name"`
	Email     string    `db:"email" json:"email"`
	Phone     string    `db:"phone" json:"phone"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	CreatedTimestamp time.Time `db:"created_timestamp" json:"created_timestamp"`
	UpdatedTimestamp time.Time `db:"updated_timestamp" json:"updated_timestamp"`
}

func (Users) TableName() string {
	return "users"
}

func (Users) GetID() string {
	return "users_id"
}

func init() {
	migration.RegisterSeeder("Users", func() interface{} {
		return SeedUsers()
	})
}

func SeedUsers() []Users {
	return []Users{
		{
			Username:         "admin",
			Password:         "admin123",
			Role:             "ADMIN",
			FullName:         "Administrator",
			Email:            "admin@fleetify.com",
			Phone:            "081234567890",
			IsActive:         true,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			CreatedTimestamp: time.Now(),
			UpdatedTimestamp: time.Now(),
		},
		{
			Username:         "manager1",
			Password:         "manager123",
			Role:             "MANAGER",
			FullName:         "Manager One",
			Email:            "manager1@fleetify.com",
			Phone:            "081234567891",
			IsActive:         true,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			CreatedTimestamp: time.Now(),
			UpdatedTimestamp: time.Now(),
		},
		{
			Username:         "manager2",
			Password:         "manager123",
			Role:             "MANAGER",
			FullName:         "Manager Two",
			Email:            "manager2@fleetify.com",
			Phone:            "081234567892",
			IsActive:         true,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			CreatedTimestamp: time.Now(),
			UpdatedTimestamp: time.Now(),
		},
		{
			Username:         "purchaser1",
			Password:         "purchaser123",
			Role:             "MANAGER",
			FullName:         "Purchaser One",
			Email:            "purchaser1@fleetify.com",
			Phone:            "081234567893",
			IsActive:         true,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			CreatedTimestamp: time.Now(),
			UpdatedTimestamp: time.Now(),
		},
		{
			Username:         "purchaser2",
			Password:         "purchaser123",
			Role:             "MANAGER",
			FullName:         "Purchaser Two",
			Email:            "purchaser2@fleetify.com",
			Phone:            "081234567894",
			IsActive:         true,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			CreatedTimestamp: time.Now(),
			UpdatedTimestamp: time.Now(),
		},
	}
}
