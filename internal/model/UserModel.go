package model

type UserModel struct {
	ID             string
	Username       *string
	Email          string
	Password       *string
	TokenVersion   string
	GoogleID       *string
	IsVerified     bool
	CreatedByAdmin bool
	Roles          []RoleModel
}
