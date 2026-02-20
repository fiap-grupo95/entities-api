package request

// UserCreateRequest represents the payload to create a new user.
type UserCreateRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	UserType string `json:"user_type" binding:"required"`
}

// UserUpdateRequest represents the payload to update a user.
type UserUpdateRequest struct {
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=6"`
	UserType string `json:"user_type" binding:"omitempty"`
}
