package request

type UserCreateRequest struct {
	Username string   `json:"username" binding:"required,excludesall= "`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=6"`
	Roles    []string `json:"roles" binding:"required"`
}

func (u *UserCreateRequest) Sanitize() map[string]any {
	return map[string]any{
		"username": u.Username,
		"email":    u.Email,
	}
}
