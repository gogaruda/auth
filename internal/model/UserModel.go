package model

type UserModel struct {
	ID       string
	Username string
	Email    string
	Password string
	Roles    []RoleModel
}
