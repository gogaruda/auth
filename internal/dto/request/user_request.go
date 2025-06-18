package request

type UpdateUserRequest struct {
	Username string   `json:"username" binding:"required"`
	Email    string   `json:"email" binding:"required,email"`
	RoleIDs  []string `json:"role_ids" binding:"required"`
}
