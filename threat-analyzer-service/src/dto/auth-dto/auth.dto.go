package authdto

// RegisterRequestDTO defines the request payload for user registration
type RegisterRequestDTO struct {
	Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Password string `json:"password" validate:"required,min=8,max=100"`
	Email    string `json:"email" validate:"required,min=8,max=100"`
}

// RegisterResponseDTO defines the response payload for user registration
type RegisterResponseDTO struct {
	Message string `json:"message"`
}

// LoginRequestDTO defines the query parameters for user login
type LoginRequestDTO struct {
	Username string `form:"username" validate:"required,min=3,max=50,alphanum"`
	Password string `form:"password" validate:"required,min=8,max=100"`
}

// LoginResponseDTO defines the response payload for user login
type LoginResponseDTO struct {
	Token string `json:"token"`
}
