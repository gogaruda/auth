package model

import (
	"time"
)

type UserModel struct {
	ID           string
	Username     string
	Email        string
	Password     string
	TokenVersion string
	Roles        []RoleModel
	UpdatedAt    time.Time
	CreatedAt    time.Time
}
