package model

import "time"

type RoleModel struct {
	ID        string
	Name      string
	Users     []UserModel
	UpdatedAt time.Time
	CreatedAt time.Time
}
