package response

type UserSwaggerResponse struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    []UserResponse `json:"data"`
}

type AuthSwaggerResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}
