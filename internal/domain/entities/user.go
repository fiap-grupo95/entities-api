package entities

import (
	"mecanica_xpto/internal/domain/valueobject"
	"time"
)

type User struct {
	ID        uint
	Email     string
	Password  string
	UserType  valueobject.UserType
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Customer  *Customer
}
