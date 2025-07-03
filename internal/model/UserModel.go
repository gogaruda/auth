package model

type UserModel struct {
	ID           string
	Username     string
	Email        string
	Password     string
	TokenVersion string
	IsVerified   bool
	Roles        []RoleModel
}
