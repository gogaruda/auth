package response

type UserResponse struct {
	ID             string
	Username       *string
	Email          string
	GoogleID       *string
	IsVerified     bool
	CreatedByAdmin bool
	Roles          []RoleResponse
}
