package models

import (
	"time"
)

type Users struct {
	UsersId string `db:"users_id" json:"users_id"`
	Username string `db:"username,unique,notnull" json:"username"`
	Password string `db:"password,notnull" json:"password"`
	Role string `db:"role,notnull" json:"role"`
	FullName string `db:"full_name,notnull" json:"full_name"`
	Email string `db:"email" json:"email"`
	Phone string `db:"phone" json:"phone"`
	IsActive bool `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	// Timestamps
	CreatedTimestamp time.Time `db:"created_timestamp" json:"created_timestamp"`
	UpdatedTimestamp time.Time `db:"updated_timestamp" json:"updated_timestamp"`
}

func (Users) TableName() string {
	return "users"
}

func (Users) GetID() string {
	return "users_id"
}
