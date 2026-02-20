package response

// UserResponse represents the response payload for user operations.
type UserResponse struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	UserType string `json:"user_type"`
}
