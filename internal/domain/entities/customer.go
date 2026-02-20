package entities

import (
	"mecanica_xpto/internal/domain/valueobject"
)

type Customer struct {
	ID          uint                `json:"id"`
	UserID      uint                `json:"user_id"`
	User        *User               `json:"user-example,omitempty"`
	CpfCnpj     valueobject.CpfCnpj `json:"document"`
	PhoneNumber string              `json:"phone_number"`
	FullName    string              `json:"full_name"`
	Email       string              `json:"email"`
	Vehicles    []Vehicle           `json:"vehicles,omitempty"`
}
