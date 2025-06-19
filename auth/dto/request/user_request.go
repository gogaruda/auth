package request

type CreateUserRequest struct {
	Username string   `json:"username" binding:"required"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=6"`
	RoleIDs  []string `json:"role_ids" binding:"required"`
}

type UpdateUserRequest struct {
	Username string   `json:"username" binding:"required"`
	Email    string   `json:"email" binding:"required,email"`
	RoleIDs  []string `json:"role_ids" binding:"required"`
}
