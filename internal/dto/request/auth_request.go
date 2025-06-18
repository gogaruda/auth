package request

type AuthLoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

func (l *AuthLoginRequest) Sanitize() map[string]any {
	return map[string]any{
		"identifier": l.Identifier,
	}
}

type AuthRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"required,min=6"`
}

func (r *AuthRegisterRequest) Sanitize() map[string]any {
	return map[string]any{
		"identifier": r.Username,
		"email":      r.Email,
	}
}
