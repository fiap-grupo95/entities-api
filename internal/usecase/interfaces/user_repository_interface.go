package interfaces

import (
	"mecanica_xpto/internal/domain/entities"
)

type IUserRepository interface {
	GetByID(id uint) (*entities.User, error)
	GetByEmail(email string) (*entities.User, error)
	Create(User *entities.User) error
	Update(User *entities.User) error
	Delete(id uint) error
	List() ([]entities.User, error)
}
