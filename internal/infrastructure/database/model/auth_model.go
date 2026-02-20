package dto

type AuthDTO struct {
	Email    string `json:"email" binding:"required,email" example:"admin@xpto.com"`
	Password string `json:"password" binding:"required" example:"Q1w2e3r%"`
}
